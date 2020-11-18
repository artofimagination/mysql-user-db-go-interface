package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func AddAsset(tx *sql.Tx) (*uuid.UUID, error) {
	// Prepare data
	queryString := "INSERT INTO user_assets (id, refs) VALUES (UUID_TO_BIN(?), ?)"

	newID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	binary, err := json.Marshal(models.References{})
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

func UpdateAsset(asset *models.Asset) error {
	// Prepare data
	queryString := "UPDATE user_assets set refs = ? where id = UUID_TO_BIN(?)"

	refRaw, err := ConvertToJSONRaw(&asset.References)
	if err != nil {
		return err
	}

	// Execute transaction
	tx, err := DBInterface.ConnectSystem()
	if err != nil {
		return err
	}

	query, err := tx.Query(queryString, refRaw, asset.ID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	defer query.Close()
	return tx.Commit()
}

func GetAsset(assetID *uuid.UUID) (*models.Asset, error) {
	asset := models.Asset{}
	queryString := "SELECT BIN_TO_UUID(id), refs FROM user_assets WHERE id = UUID_TO_BIN(?)"

	tx, err := DBInterface.ConnectSystem()
	if err != nil {
		return nil, err
	}

	query := tx.QueryRow(queryString, *assetID)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	refs := json.RawMessage{}
	if err := query.Scan(&asset.ID, &refs); err != nil {
		return nil, RollbackWithErrorStack(tx, errors.Wrap(errors.WithStack(err), fmt.Sprintf("Asset %s not found", assetID.String())))
	}

	if err := json.Unmarshal(refs, &asset.References); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return &asset, tx.Commit()
}

func DeleteAsset(assetID *uuid.UUID) error {
	query := "DELETE FROM user_assets WHERE id=UUID_TO_BIN(?)"
	tx, err := DBInterface.ConnectSystem()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, *assetID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}
	return tx.Commit()
}
