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
	UserAssets     = "user_assets"
	ProductAssets  = "product_assets"
	ProductDetails = "product_details"
	UserSettings   = "user_settings"
)

var ErrAssetMissing = "This %s is missing"

var AddAssetQuery = "INSERT INTO ? (id, data) VALUES (UUID_TO_BIN(?), ?)"

func (MYSQLFunctions) AddAsset(assetType string, asset *models.Asset, tx *sql.Tx) error {
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

func UpdateAsset(assetType string, asset *models.Asset) error {
	// Prepare data
	queryString := "UPDATE ? set data = ? where id = UUID_TO_BIN(?)"

	refRaw, err := ConvertToJSONRaw(&asset.DataMap)
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
	queryString := "SELECT BIN_TO_UUID(id), data FROM ? WHERE id = UUID_TO_BIN(?)"

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

	if err := json.Unmarshal(refs, &asset.DataMap); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return &asset, tx.Commit()
}

var DeleteAssetQuery = "DELETE FROM ? WHERE id=UUID_TO_BIN(?)"

func (MYSQLFunctions) DeleteAsset(assetType string, assetID *uuid.UUID, tx *sql.Tx) error {
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
