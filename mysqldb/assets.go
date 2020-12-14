package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

const (
	UserAssets     = "user_assets"
	ProductAssets  = "product_assets"
	ProductDetails = "product_details"
	UserSettings   = "user_settings"
	ProjectDetails = "project_details"
	ProjectAssets  = "project_assets"
)

var ErrAssetMissing = "This %s is missing"

var AddAssetQuery = "INSERT INTO %s (id, data) VALUES (UUID_TO_BIN(?), ?)"

func (*MYSQLFunctions) AddAsset(assetType string, asset *models.Asset, tx *sql.Tx) error {
	// Prepare data
	binary, err := json.Marshal(asset.DataMap)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	// Execute transaction
	query := fmt.Sprintf(AddAssetQuery, assetType)
	_, err = tx.Exec(query, asset.ID, binary)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	return nil
}

var UpdateAssetQuery = "UPDATE %s set data = ? where id = UUID_TO_BIN(?)"

func UpdateAsset(assetType string, asset *models.Asset) error {
	binary, err := json.Marshal(asset.DataMap)
	if err != nil {
		return err
	}

	// Execute transaction
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(UpdateAssetQuery, assetType)
	result, err := tx.Exec(query, binary, asset.ID)
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

var GetAssetQuery = "SELECT BIN_TO_UUID(id), data FROM %s WHERE id = UUID_TO_BIN(?)"

func GetAsset(assetType string, assetID *uuid.UUID) (*models.Asset, error) {
	asset := &models.Asset{}

	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(GetAssetQuery, assetType)
	result := tx.QueryRow(query, assetID)

	dataMap := []byte{}
	err = result.Scan(&asset.ID, &dataMap)
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

	return asset, tx.Commit()
}

var DeleteAssetQuery = "DELETE FROM %s WHERE id=UUID_TO_BIN(?)"

func (*MYSQLFunctions) DeleteAsset(assetType string, assetID *uuid.UUID, tx *sql.Tx) error {
	query := fmt.Sprintf(DeleteAssetQuery, assetType)
	result, err := tx.Exec(query, *assetID)
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

var GetAssetsQuery = "SELECT BIN_TO_UUID(id), data FROM %s WHERE id IN (UUID_TO_BIN(?)"

func (MYSQLFunctions) GetAssets(assetType string, IDs []uuid.UUID, tx *sql.Tx) ([]models.Asset, error) {
	query := GetAssetsQuery + strings.Repeat(",UUID_TO_BIN(?)", len(IDs)-1) + ")"
	interfaceList := make([]interface{}, len(IDs))
	for i := range IDs {
		interfaceList[i] = IDs[i]
	}
	query = fmt.Sprintf(query, assetType)
	rows, err := tx.Query(query, interfaceList...)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	defer rows.Close()

	assets := make([]models.Asset, 0)
	for rows.Next() {
		dataMap := []byte{}
		asset := models.Asset{}
		err := rows.Scan(&asset.ID, &dataMap)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		if err := json.Unmarshal(dataMap, &asset.DataMap); err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		assets = append(assets, asset)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	if len(assets) == 0 {
		return nil, sql.ErrNoRows
	}

	return assets, nil
}
