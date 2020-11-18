package mysqldb

import (
	"database/sql"
	"encoding/json"

	"github.com/artofimagination/mysql-user-db-go-interface/models"

	"github.com/google/uuid"
)

func AddSettings(tx *sql.Tx) (*uuid.UUID, error) {
	queryString := "INSERT INTO user_settings (id, settings) VALUES (UUID_TO_BIN(?), ?)"

	newID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	binary, err := json.Marshal(models.Settings{})
	if err != nil {
		return nil, err
	}

	query, err := tx.Query(queryString, newID, binary)
	if err != nil {
		return nil, err
	}

	defer query.Close()
	return &newID, nil
}

func GetSettings(settingsID *uuid.UUID) (*models.UserSetting, error) {
	settings := models.UserSetting{}
	queryString := "SELECT settings FROM user_settings WHERE id = UUID_TO_BIN(?)"
	tx, err := DBInterface.ConnectSystem()
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
	tx, err := DBInterface.ConnectSystem()
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
