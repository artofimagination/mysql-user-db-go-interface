package mysqldb

import (
	"database/sql"
	"encoding/json"

	"github.com/artofimagination/mysql-user-db-go-interface/models"

	"github.com/google/uuid"
)

var AddUserSettingsQuery = "INSERT INTO user_settings (id, settings) VALUES (UUID_TO_BIN(?), ?)"

func (MYSQLFunctions) AddSettings(settings *models.UserSettings, tx *sql.Tx) error {

	binary, err := json.Marshal(settings.Settings)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	_, err = tx.Exec(AddUserSettingsQuery, settings.ID, binary)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	return tx.Commit()
}

func GetSettings(settingsID *uuid.UUID) (*models.UserSettings, error) {
	settings := models.UserSettings{}
	queryString := "SELECT settings FROM user_settings WHERE id = UUID_TO_BIN(?)"
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	query := tx.QueryRow(queryString, *settingsID)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	settingsJSON := json.RawMessage{}
	if err := query.Scan(&settingsJSON); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	if err := json.Unmarshal(settingsJSON, &settings.Settings); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return &settings, tx.Commit()
}

func DeleteSettings(settingsID *uuid.UUID) error {
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	if err := deleteSettings(settingsID, tx); err != nil {
		return RollbackWithErrorStack(tx, err)
	}
	return tx.Commit()
}

func deleteSettings(settingsID *uuid.UUID, tx *sql.Tx) error {
	query := "DELETE FROM user_settings WHERE id=UUID_TO_BIN(?)"

	_, err := tx.Exec(query, *settingsID)
	if err != nil {
		return err
	}
	return nil
}
