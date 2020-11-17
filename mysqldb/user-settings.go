package mysqldb

import (
	"encoding/json"

	"github.com/artofimagination/mysql-user-db-go-interface/models"

	"github.com/google/uuid"
)

func AddSettings() (*uuid.UUID, error) {
	queryString := "INSERT INTO user_settings (id, settings) VALUES (UUID_TO_BIN(?), ?)"
	db, err := DBInterface.ConnectSystem()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	newID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	binary, err := json.Marshal(models.Settings{})
	if err != nil {
		return nil, err
	}

	query, err := db.Query(queryString, newID, binary)
	if err != nil {
		return nil, err
	}

	defer query.Close()
	return &newID, nil
}

func GetSettings(settingsID *uuid.UUID) (*models.UserSetting, error) {
	settings := models.UserSetting{}
	queryString := "SELECT settings FROM user_settings WHERE id = UUID_TO_BIN(?)"
	db, err := DBInterface.ConnectSystem()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := db.QueryRow(queryString, *settingsID)

	if err != nil {
		return nil, err
	}

	settingsJSON := json.RawMessage{}
	if err := query.Scan(&settingsJSON); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(settingsJSON, &settings.Settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func DeleteSettings(settingsID *uuid.UUID) error {
	query := "DELETE FROM user_settings WHERE id=UUID_TO_BIN(?)"
	db, err := DBInterface.ConnectSystem()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query, *settingsID)
	if err != nil {
		return err
	}
	return nil
}
