package mysqldb

import (
	"database/sql"
	"encoding/json"

	"github.com/artofimagination/mysql-user-db-go-interface/models"

	"github.com/google/uuid"
)

var AddProductDetailsQuery = "INSERT INTO product_details (id, details) VALUES (UUID_TO_BIN(?), ?)"

func (MYSQLFunctions) AddDetails(details *models.ProductDetails, tx *sql.Tx) error {
	binary, err := json.Marshal(details.Details)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	_, err = tx.Exec(AddProductDetailsQuery, details.ID, binary)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	return tx.Commit()
}

var GetProductDetailsQuery = "SELECT details FROM product_details WHERE id = UUID_TO_BIN(?)"

func GetDetails(settingsID *uuid.UUID) (*models.ProductDetails, error) {
	details := models.ProductDetails{}
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	query := tx.QueryRow(GetProductDetailsQuery, *settingsID)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	detailsJSON := json.RawMessage{}
	if err := query.Scan(&detailsJSON); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	if err := json.Unmarshal(detailsJSON, &details.Details); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return &details, tx.Commit()
}

var DeleteDetailsQuery = "DELETE FROM user_settings WHERE id=UUID_TO_BIN(?)"

func DeleteDetails(settingsID *uuid.UUID, tx *sql.Tx) error {
	_, err := tx.Exec(DeleteDetailsQuery, *settingsID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}
	return tx.Commit()
}
