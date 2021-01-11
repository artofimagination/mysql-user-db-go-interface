package dbcontrollers

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/google/uuid"
)

var ErrProjectExistsString = "Project with name %s already exists"
var ErrProjectNotFound = errors.New("The selected project not found")
var ErrMissingProjectDetail = errors.New("Details for the selected project not found")
var ErrMissingProjectAsset = errors.New("Assets for the selected project not found")
var ErrEmptyProjectIDList = errors.New("Request does not contain any project identifiers")

func (*MYSQLController) CreateProject(name string, visibility string, owner *uuid.UUID, productID *uuid.UUID, generateAssetPath func(assetID *uuid.UUID) (string, error)) (*models.ProjectData, error) {
	references := make(models.DataMap)
	asset, err := models.Interface.NewAsset(references, generateAssetPath)
	if err != nil {
		return nil, err
	}

	details := make(models.DataMap)
	details["name"] = name
	details["visibility"] = visibility
	projectDetails, err := models.Interface.NewAsset(details, generateAssetPath)
	if err != nil {
		return nil, err
	}

	project, err := models.Interface.NewProject(productID, &projectDetails.ID, &asset.ID)
	if err != nil {
		return nil, err
	}

	// Start a DB transaction and do all inserts within the same transaction to improve consistency.
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	users := models.ProjectUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	privilege, err := mysqldb.Functions.GetPrivilege("Owner")
	if err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddAsset(mysqldb.ProjectDetails, projectDetails, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddAsset(mysqldb.ProjectAssets, asset, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddProject(project, tx); err != nil {
		return nil, err
	}

	users.UserMap[*owner] = privilege.ID
	if err := mysqldb.Functions.AddProjectUsers(&project.ID, &users, tx); err != nil {
		return nil, err
	}

	projectData := models.ProjectData{
		ID:        project.ID,
		ProductID: project.ProductID,
		Details:   projectDetails,
		Assets:    asset,
	}

	return &projectData, mysqldb.DBConnector.Commit(tx)
}

func (*MYSQLController) DeleteProject(projectID *uuid.UUID) error {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	if err := deleteProject(projectID, tx); err != nil {
		return err
	}

	return mysqldb.DBConnector.Commit(tx)
}

func deleteProject(projectID *uuid.UUID, tx *sql.Tx) error {
	project, err := mysqldb.Functions.GetProjectByID(projectID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrProjectNotFound
		}
	}

	if err := mysqldb.Functions.DeleteProjectUsersByProjectID(projectID, tx); err != nil {
		if err == mysqldb.ErrNoProductDeleted {
			return ErrProjectNotFound
		}
	}

	if err := mysqldb.Functions.DeleteProject(projectID, tx); err != nil {
		if err == mysqldb.ErrNoProductDeleted {
			return ErrProjectNotFound
		}
		return err
	}

	if err := mysqldb.Functions.DeleteAsset(mysqldb.ProjectAssets, &project.AssetsID, tx); err != nil {
		return err
	}

	if err := mysqldb.Functions.DeleteAsset(mysqldb.ProjectDetails, &project.DetailsID, tx); err != nil {
		return err
	}

	return nil
}

func (*MYSQLController) GetProject(projectID *uuid.UUID) (*models.ProjectData, error) {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	project, err := mysqldb.Functions.GetProjectByID(projectID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := mysqldb.DBConnector.Rollback(tx); err != nil {
				return nil, err
			}
			return nil, ErrProjectNotFound
		}
	}

	details, err := mysqldb.GetAsset(mysqldb.ProjectDetails, &project.DetailsID)
	if err != nil {
		return nil, err
	}

	assets, err := mysqldb.GetAsset(mysqldb.ProjectAssets, &project.AssetsID)
	if err != nil {
		return nil, err
	}

	projectData := models.ProjectData{
		ID:        project.ID,
		ProductID: project.ProductID,
		Details:   details,
		Assets:    assets,
	}

	return &projectData, mysqldb.DBConnector.Commit(tx)
}

func (*MYSQLController) UpdateProjectDetails(projectData *models.ProjectData) error {

	if err := mysqldb.UpdateAsset(mysqldb.ProjectDetails, projectData.Details); err != nil {
		if fmt.Errorf(mysqldb.ErrAssetMissing, mysqldb.ProjectDetails).Error() == err.Error() {
			return ErrMissingProjectDetail
		}
		return err
	}
	return nil
}

func (*MYSQLController) UpdateProjectAssets(projectData *models.ProjectData) error {
	if err := mysqldb.UpdateAsset(mysqldb.ProjectAssets, projectData.Assets); err != nil {
		if fmt.Errorf(mysqldb.ErrAssetMissing, mysqldb.ProjectAssets).Error() == err.Error() {
			return ErrMissingProjectAsset
		}
		return err
	}
	return nil
}

func (*MYSQLController) UpdateProjectUser(projectID *uuid.UUID, userID *uuid.UUID, privilege int) error {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	return mysqldb.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetProjectsByUserID(userID *uuid.UUID) ([]models.UserProject, error) {
	projects := make([]models.UserProject, 0)
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	ownershipMap, err := mysqldb.Functions.GetUserProjectIDs(userID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoProjectsForUser
		}
		return nil, err
	}

	for projectID, privilege := range ownershipMap.ProjectMap {
		projectID := projectID
		project, err := c.GetProject(&projectID)
		if err != nil {
			return nil, err
		}

		userProject := models.UserProject{
			ProjectData: project,
			Privilege:   privilege,
		}

		projects = append(projects, userProject)
	}

	return projects, mysqldb.DBConnector.Commit(tx)
}

func (*MYSQLController) GetProjects(projectIDs []uuid.UUID) ([]models.ProjectData, error) {
	if len(projectIDs) == 0 {
		return nil, ErrEmptyProjectIDList
	}

	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	projects, err := mysqldb.Functions.GetProjectsByIDs(projectIDs, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := mysqldb.DBConnector.Rollback(tx); err != nil {
				return nil, err
			}
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	assetIDs := make([]uuid.UUID, 0)
	detailsIDs := make([]uuid.UUID, 0)
	for _, product := range projects {
		assetIDs = append(assetIDs, product.AssetsID)
		detailsIDs = append(detailsIDs, product.DetailsID)
	}

	details, err := mysqldb.Functions.GetAssets(mysqldb.ProjectDetails, detailsIDs, tx)
	if err != nil {
		return nil, err
	}

	assets, err := mysqldb.Functions.GetAssets(mysqldb.ProjectAssets, assetIDs, tx)
	if err != nil {
		return nil, err
	}

	projectDataList := make([]models.ProjectData, len(projects))
	for index, project := range projects {
		projectData := models.ProjectData{
			ID:        project.ID,
			ProductID: project.ProductID,
			Details:   &details[index],
			Assets:    &assets[index],
		}
		projectDataList[index] = projectData
	}

	return projectDataList, mysqldb.DBConnector.Commit(tx)
}
