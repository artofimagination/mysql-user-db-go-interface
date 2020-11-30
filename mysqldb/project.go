package mysqldb

import (
	"github.com/artofimagination/mysql-user-db-go-interface/models"
)

// UpdateProject updates the selected project.
func UpdateProject(project models.Project) error {
	query := "UPDATE projects set user_id = ?, features_id = ?, name = ?, config = ? where id = ?"
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	_, err = tx.Exec(query, project.UserID, project.FeatureID, project.Name, project.Config, project.ID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}
	return tx.Commit()
}

// AddProject adds a new project to the database.
func AddProject(project models.Project) error {
	query := "INSERT INTO projects (id, user_id, features_id, name, config) VALUES (UUID_TO_BIN(UUID()), UUID_TO_BIN(?), ?, ?, ?)"
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, project.UserID, project.FeatureID, project.Name, project.Config)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}
	return tx.Commit()
}

// GetProjectByName returns the project by name.
func GetProjectByName(name string) (*models.Project, error) {
	var project models.Project
	queryString := "select BIN_TO_UUID(id), name, config from projects where name = ?"
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	query, err := tx.Query(queryString, name)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	defer query.Close()

	query.Next()
	if err := query.Scan(&project.ID, &project.Name, &project.Config); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return &project, tx.Commit()
}
