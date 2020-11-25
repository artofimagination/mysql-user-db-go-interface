package mysqldb

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

func createTestAsset() (*models.Asset, error) {
	references := make(models.References)

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	asset := models.Asset{
		ID:         assetID,
		References: references,
	}
	return &asset, nil
}

func TestAddAsset_ValidUserAsset(t *testing.T) {
	// Create test data
	asset, err := createTestAsset()
	if err != nil {
		t.Errorf("Failed to generate asset test data %s", err)
		return
	}

	binary, err := json.Marshal(asset.References)
	if err != nil {
		t.Errorf("Failed to marshal asset references map %s", err)
		return
	}

	// Prepare mock
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate user test data %s", err)
		return
	}

	mock.ExpectBegin()
	mock.ExpectExec(AddAssetQuery).WithArgs(UserAssets, asset.ID, binary).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}
	defer db.Close()

	FunctionInterface = MYSQLFunctionInterface{}

	// Run test
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to setup DB transaction %s", err)
		return
	}

	err = FunctionInterface.AddAsset(UserAssets, asset, tx)
	if err != nil {
		t.Errorf("Failed to add asset %s", err)
		return
	}
}
