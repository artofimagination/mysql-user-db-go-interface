package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

const (
	UserAssets     = "user_assets"
	ProductAssets  = "product_assets"
	ProductDetails = "product_details"
	UserSettings   = "user_settings"
)

var ErrAssetMissing = "This %s is missing"

var AddAssetQuery = "INSERT INTO ? (id, data) VALUES (UUID_TO_BIN(?), ?)"

func (*MYSQLFunctions) AddAsset(assetType string, asset *models.Asset, tx *sql.Tx) error {
	// Prepare data
	binary, err := json.Marshal(asset.DataMap)
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

var UpdateAssetQuery = "UPDATE ? set data = ? where id = UUID_TO_BIN(?)"

func UpdateAsset(assetType string, asset *models.Asset) error {
	refRaw, err := ConvertToJSONRaw(&asset.DataMap)
	if err != nil {
		return err
	}

	// Execute transaction
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	result, err := tx.Exec(UpdateAssetQuery, assetType, refRaw, asset.ID)
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
		return fmt.Errorf(ErrAssetMissing, assetType)
	}

	return tx.Commit()
}

var GetAssetQuery = "SELECT BIN_TO_UUID(id), data FROM ? WHERE id = UUID_TO_BIN(?)"

func GetAsset(assetType string, assetID *uuid.UUID) (*models.Asset, error) {
	asset := models.Asset{}

	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	query := tx.QueryRow(GetAssetQuery, assetType, assetID)

	dataMap := []byte{}
	err = query.Scan(&asset.ID, &dataMap)
	switch {
	case err == sql.ErrNoRows:
		if errRb := tx.Commit(); errRb != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	if err := json.Unmarshal(dataMap, &asset.DataMap); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return &asset, tx.Commit()
}

var DeleteAssetQuery = "DELETE FROM ? WHERE id=UUID_TO_BIN(?)"

func (*MYSQLFunctions) DeleteAsset(assetType string, assetID *uuid.UUID, tx *sql.Tx) error {
	result, err := tx.Exec(DeleteAssetQuery, assetType, *assetID)
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
		return fmt.Errorf(ErrAssetMissing, assetType)
	}
	return nil
}
