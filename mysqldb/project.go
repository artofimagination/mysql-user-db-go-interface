package mysqldb

import (
	"database/sql"
	"strings"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrNoUserWithProject = errors.New("No user is associated to this project")
var ErrNoProjectUserAdded = errors.New("No project user relation has been added")
var ErrNoProjectDeleted = errors.New("No project was deleted")
var ErrNoUsersProjectUpdate = errors.New("No users project was updated")

var AddProjectUsersQuery = "INSERT INTO users_projects (users_id, projects_id, privileges_id) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?)"

func (*MYSQLFunctions) AddProjectUsers(projectID *uuid.UUID, projectUsers *models.ProjectUserIDs, tx *sql.Tx) error {
	for userID, privilege := range projectUsers.UserMap {
		result, err := tx.Exec(AddProjectUsersQuery, userID, projectID, privilege)
		if err != nil {
			return RollbackWithErrorStack(tx, err)
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return RollbackWithErrorStack(tx, err)
		}

		if affected == 0 {
			if errRb := tx.Rollback(); errRb != nil {
				return err
			}
			return ErrNoProjectUserAdded
		}
	}
	return nil
}

var DeleteProjectUsersByProjectIDQuery = "DELETE FROM users_projects where projects_id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) DeleteProjectUsersByProjectID(projectID *uuid.UUID, tx *sql.Tx) error {
	result, err := tx.Exec(DeleteProjectUsersByProjectIDQuery, projectID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		if errRb := tx.Rollback(); errRb != nil {
			return err
		}
		return ErrNoUserWithProject
	}

	return nil
}

var UpdateUsersProjectsQuery = "UPDATE users_projects set privileges_id = ? where users_id = UUID_TO_BIN(?) AND projects_id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) UpdateUsersProjects(userID *uuid.UUID, projectID *uuid.UUID, privilege int, tx *sql.Tx) error {
	result, err := tx.Exec(UpdateUsersProjectsQuery, privilege, userID, projectID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		if errRb := tx.Rollback(); errRb != nil {
			return err
		}
		return ErrNoUsersProjectUpdate
	}

	return nil
}

var AddProjectQuery = "INSERT INTO projects (id, products_id, project_details_id, project_assets_id) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?))"

func (*MYSQLFunctions) AddProject(project *models.Project, tx *sql.Tx) error {
	_, err := tx.Exec(AddProjectQuery, project.ID, project.ProductID, project.DetailsID, project.AssetsID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}
	return nil
}

var GetProjectByIDQuery = "SELECT BIN_TO_UUID(id), BIN_TO_UUID(products_id), BIN_TO_UUID(project_details_id), BIN_TO_UUID(project_assets_id) FROM projects WHERE id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) GetProjectByID(ID *uuid.UUID, tx *sql.Tx) (*models.Project, error) {
	project := &models.Project{}
	query := tx.QueryRow(GetProjectByIDQuery, ID)
	err := query.Scan(&project.ID, &project.ProductID, &project.DetailsID, &project.AssetsID)
	switch {
	case err == sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	return project, nil
}

var GetProjectsByIDsQuery = "SELECT BIN_TO_UUID(id), BIN_TO_UUID(products_id), BIN_TO_UUID(project_details_id), BIN_TO_UUID(project_assets_id) FROM projects WHERE id IN (UUID_TO_BIN(?)"

func (*MYSQLFunctions) GetProjectsByIDs(IDs []uuid.UUID, tx *sql.Tx) ([]models.Project, error) {
	query := GetProjectsByIDsQuery + strings.Repeat(",UUID_TO_BIN(?)", len(IDs)-1) + ")"
	interfaceList := make([]interface{}, len(IDs))
	for i := range IDs {
		interfaceList[i] = IDs[i]
	}
	rows, err := tx.Query(query, interfaceList...)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	defer rows.Close()

	projects := make([]models.Project, 0)
	for rows.Next() {
		project := models.Project{}
		err := rows.Scan(&project.ID, &project.ProductID, &project.DetailsID, &project.AssetsID)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		projects = append(projects, project)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	if len(projects) == 0 {
		return nil, sql.ErrNoRows
	}

	return projects, nil
}

var GetUserProjectIDsQuery = "SELECT BIN_TO_UUID(projects_id), privileges_id FROM users_projects where users_id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) GetUserProjectIDs(userID *uuid.UUID, tx *sql.Tx) (*models.UserProjectIDs, error) {
	rows, err := tx.Query(GetUserProjectIDsQuery, userID)
	switch {
	case err == sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	defer rows.Close()
	userProjects := &models.UserProjectIDs{
		ProjectMap:     make(map[uuid.UUID]int),
		ProjectIDArray: make([]uuid.UUID, 0),
	}
	for rows.Next() {
		projectID := uuid.UUID{}
		privilege := -1
		err := rows.Scan(&projectID, &privilege)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		userProjects.ProjectMap[projectID] = privilege
		userProjects.ProjectIDArray = append(userProjects.ProjectIDArray, projectID)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}
	return userProjects, nil
}

var GetProductProjectsQuery = "SELECT BIN_TO_UUID(id), BIN_TO_UUID(products_id), BIN_TO_UUID(project_details_id), BIN_TO_UUID(project_assets_id) FROM projects where products_id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) GetProductProjects(productID *uuid.UUID, tx *sql.Tx) ([]models.Project, error) {
	projects := make([]models.Project, 0)
	rows, err := tx.Query(GetProductProjectsQuery, &productID)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	defer rows.Close()
	for rows.Next() {
		project := models.Project{}
		err := rows.Scan(&project.ID, &project.ProductID, &project.DetailsID, &project.AssetsID)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		projects = append(projects, project)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	if len(projects) == 0 {
		return nil, sql.ErrNoRows
	}

	return projects, nil
}

var DeleteProjectQuery = "DELETE FROM projects where id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) DeleteProject(projectID *uuid.UUID, tx *sql.Tx) error {
	result, err := tx.Exec(DeleteProjectQuery, projectID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		if errRb := tx.Rollback(); errRb != nil {
			return err
		}
		return ErrNoProjectDeleted
	}

	return nil
}

var DeleteProjectsByProductIDQuery = "DELETE FROM projects where products_id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) DeleteProjectsByProductID(productID *uuid.UUID, tx *sql.Tx) error {
	result, err := tx.Exec(DeleteProjectsByProductIDQuery, productID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		return ErrNoProjectDeleted
	}

	return nil
}
