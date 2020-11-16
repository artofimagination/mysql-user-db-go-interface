package mysqldb

import (
	"github.com/artofimagination/mysql-user-db-go-interface/models"
)

// UpdateProject updates the selected project.
func UpdateProject(project models.Project) error {
	query := "UPDATE projects set user_id = ?, features_id = ?, name = ?, config = ? where id = ?"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Query(query, project.UserID, project.FeatureID, project.Name, project.Config, project.ID)
	if err != nil {
		return err
	}
	return nil
}

// AddProject adds a new project to the database.
func AddProject(project models.Project) error {
	query := "INSERT INTO projects (id, user_id, features_id, name, config) VALUES (UUID_TO_BIN(UUID()), UUID_TO_BIN(?), ?, ?, ?)"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Query(query, project.UserID, project.FeatureID, project.Name, project.Config)
	if err != nil {
		return err
	}
	return nil
}

// GetProjectByName returns the project by name.
func GetProjectByName(name string) (models.Project, error) {
	var project models.Project
	queryString := "select BIN_TO_UUID(id), name, config from projects where name = ?"
	db, err := ConnectSystem()
	if err != nil {
		return project, err
	}
	defer db.Close()

	query, err := db.Query(queryString, name)

	if err != nil {
		return project, err
	}

	defer query.Close()

	query.Next()
	if err := query.Scan(&project.ID, &project.Name, &project.Config); err != nil {
		return project, err
	}

	return project, nil
}
