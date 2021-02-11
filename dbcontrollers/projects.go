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
var ErrNoProjectForProduct = errors.New("No projects for this product")
var ErrNoProjectDetailsUpdate = errors.New("Details for the selected project not found or no change happened")
var ErrNoProjectAssetsUpdate = errors.New("Assets for the selected project not found or no change happened")
var ErrEmptyProjectIDList = errors.New("Request does not contain any project identifiers")

func (c *MYSQLController) CreateProject(name string, visibility string, owner *uuid.UUID, productID *uuid.UUID, generateAssetPath func(assetID *uuid.UUID) (string, error)) (*models.ProjectData, error) {
	references := make(models.DataMap)
	asset, err := c.ModelFunctions.NewAsset(references, generateAssetPath)
	if err != nil {
		return nil, err
	}

	details := make(models.DataMap)
	details["name"] = name
	details["visibility"] = visibility
	projectDetails, err := c.ModelFunctions.NewAsset(details, generateAssetPath)
	if err != nil {
		return nil, err
	}

	project, err := c.ModelFunctions.NewProject(productID, &projectDetails.ID, &asset.ID)
	if err != nil {
		return nil, err
	}

	// Start a DB transaction and do all inserts within the same transaction to improve consistency.
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	users := models.ProjectUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	privilege, err := c.DBFunctions.GetPrivilege("Owner")
	if err != nil {
		return nil, err
	}

	if err := c.DBFunctions.AddAsset(mysqldb.ProjectDetails, projectDetails, tx); err != nil {
		return nil, err
	}

	if err := c.DBFunctions.AddAsset(mysqldb.ProjectAssets, asset, tx); err != nil {
		return nil, err
	}

	if err := c.DBFunctions.AddProject(project, tx); err != nil {
		return nil, err
	}

	users.UserMap[*owner] = privilege.ID
	if err := c.DBFunctions.AddProjectUsers(&project.ID, &users, tx); err != nil {
		return nil, err
	}

	projectData := models.ProjectData{
		ID:        project.ID,
		ProductID: project.ProductID,
		Details:   projectDetails,
		Assets:    asset,
	}

	return &projectData, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) DeleteProject(projectID *uuid.UUID) error {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	if err := c.deleteProject(projectID, tx); err != nil {
		return err
	}

	return c.DBConnector.Commit(tx)
}

func (c *MYSQLController) deleteProject(projectID *uuid.UUID, tx *sql.Tx) error {
	project, err := c.DBFunctions.GetProjectByID(projectID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrProjectNotFound
		}
	}

	if err := c.DBFunctions.DeleteProjectUsersByProjectID(projectID, tx); err != nil {
		if err == mysqldb.ErrNoProductDeleted {
			return ErrProjectNotFound
		}
	}

	if err := c.DBFunctions.DeleteProject(projectID, tx); err != nil {
		if err == mysqldb.ErrNoProductDeleted {
			return ErrProjectNotFound
		}
		return err
	}

	if err := c.DBFunctions.DeleteAsset(mysqldb.ProjectAssets, &project.AssetsID, tx); err != nil {
		return err
	}

	if err := c.DBFunctions.DeleteAsset(mysqldb.ProjectDetails, &project.DetailsID, tx); err != nil {
		return err
	}

	return nil
}

func (c *MYSQLController) GetProject(projectID *uuid.UUID) (*models.ProjectData, error) {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	project, err := c.DBFunctions.GetProjectByID(projectID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := c.DBConnector.Rollback(tx); err != nil {
				return nil, err
			}
			return nil, ErrProjectNotFound
		}
	}

	details, err := c.DBFunctions.GetAsset(mysqldb.ProjectDetails, &project.DetailsID)
	if err != nil {
		return nil, err
	}

	assets, err := c.DBFunctions.GetAsset(mysqldb.ProjectAssets, &project.AssetsID)
	if err != nil {
		return nil, err
	}

	projectData := models.ProjectData{
		ID:        project.ID,
		ProductID: project.ProductID,
		Details:   details,
		Assets:    assets,
	}

	return &projectData, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) UpdateProjectDetails(projectData *models.ProjectData) error {
	if err := c.DBFunctions.UpdateAsset(mysqldb.ProjectDetails, projectData.Details); err != nil {
		if fmt.Errorf(mysqldb.ErrAssetMissing, mysqldb.ProjectDetails).Error() == err.Error() {
			return ErrNoProjectDetailsUpdate
		}
		return err
	}
	return nil
}

func (c *MYSQLController) UpdateProjectAssets(projectData *models.ProjectData) error {
	if err := c.DBFunctions.UpdateAsset(mysqldb.ProjectAssets, projectData.Assets); err != nil {
		if fmt.Errorf(mysqldb.ErrAssetMissing, mysqldb.ProjectAssets).Error() == err.Error() {
			return ErrNoProjectAssetsUpdate
		}
		return err
	}
	return nil
}

func (c *MYSQLController) UpdateProjectUser(projectID *uuid.UUID, userID *uuid.UUID, privilege int) error {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	return c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetProjectsByUserID(userID *uuid.UUID) ([]models.UserProject, error) {
	projects := make([]models.UserProject, 0)
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	ownershipMap, err := c.DBFunctions.GetUserProjectIDs(userID, tx)
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

	return projects, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) buildProjectData(projects []models.Project, tx *sql.Tx) ([]models.ProjectData, error) {
	projectDataList := make([]models.ProjectData, len(projects))
	assetIDs := make([]uuid.UUID, 0)
	detailsIDs := make([]uuid.UUID, 0)
	for _, project := range projects {
		assetIDs = append(assetIDs, project.AssetsID)
		detailsIDs = append(detailsIDs, project.DetailsID)
	}

	details, err := c.DBFunctions.GetAssets(mysqldb.ProjectDetails, detailsIDs, tx)
	if err != nil {
		return nil, err
	}

	assets, err := c.DBFunctions.GetAssets(mysqldb.ProjectAssets, assetIDs, tx)
	if err != nil {
		return nil, err
	}

	for index, project := range projects {
		projectData := models.ProjectData{
			ID:        project.ID,
			ProductID: project.ProductID,
			Details:   &details[index],
			Assets:    &assets[index],
		}
		projectDataList[index] = projectData
	}
	return projectDataList, nil
}

func (c *MYSQLController) GetProjectsByProductID(productID *uuid.UUID) ([]models.ProjectData, error) {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	projects, err := c.DBFunctions.GetProductProjects(productID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoProjectForProduct
		}
		return nil, err
	}

	projectDataList, err := c.buildProjectData(projects, tx)
	if err != nil {
		return nil, err
	}

	return projectDataList, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetProjects(projectIDs []uuid.UUID) ([]models.ProjectData, error) {
	if len(projectIDs) == 0 {
		return nil, ErrEmptyProjectIDList
	}

	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	projects, err := c.DBFunctions.GetProjectsByIDs(projectIDs, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := c.DBConnector.Rollback(tx); err != nil {
				return nil, err
			}
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	projectDataList, err := c.buildProjectData(projects, tx)
	if err != nil {
		return nil, err
	}

	return projectDataList, c.DBConnector.Commit(tx)
}
