package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	UserAssets    = "user"
	ProductAssets = "product"
)

var AddAssetQuery = "INSERT INTO ?_assets (id, refs) VALUES (UUID_TO_BIN(?), ?)"

func (MYSQLFunctions) AddAsset(assetType string, asset *models.Asset, tx *sql.Tx) error {
	// Prepare data
	binary, err := json.Marshal(asset.References)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	// Execute transaction
	_, err = tx.Exec(AddAssetQuery, assetType, asset.ID, binary)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	return tx.Commit()
}

func UpdateAsset(assetType string, asset *models.Asset) error {
	// Prepare data
	queryString := "UPDATE ?_assets set refs = ? where id = UUID_TO_BIN(?)"

	refRaw, err := ConvertToJSONRaw(&asset.References)
	if err != nil {
		return err
	}

	// Execute transaction
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	query, err := tx.Query(queryString, assetType, refRaw, asset.ID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	defer query.Close()
	return tx.Commit()
}

func GetAsset(assetType string, assetID *uuid.UUID) (*models.Asset, error) {
	asset := models.Asset{}
	queryString := "SELECT BIN_TO_UUID(id), refs FROM ?_assets WHERE id = UUID_TO_BIN(?)"

	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	query := tx.QueryRow(queryString, assetType, *assetID)
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

func DeleteAsset(assetType string, assetID *uuid.UUID) error {
	query := "DELETE FROM ?_assets WHERE id=UUID_TO_BIN(?)"
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, assetType, *assetID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}
	return tx.Commit()
}
