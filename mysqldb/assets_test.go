package mysqldb

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func createAssetTestData() (*test.OrderedTests, DBConnectorMock, error) {
	dbConnector := DBConnectorMock{}
	dataSet := test.OrderedTests{
		orderedList: make(test.OrderedTestList, 0),
		testDataSet: make(test.DataSet, 0),
	}
	references := make(models.References)

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, dbConnector, err
	}

	asset := models.Asset{
		ID:         assetID,
		References: references,
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, dbConnector, err
	}

	binary, err := json.Marshal(asset.References)
	if err != nil {
		return nil, dbConnector, err
	}

	data := test.TestData{
		data:     asset,
		expected: nil,
	}

	testCase := "valid_user_asset"
	mock.ExpectBegin()
	mock.ExpectExec(AddAssetQuery).WithArgs(UserAssets, asset.ID, binary).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	dataSet.testDataSet[testCase] = data
	dataSet.orderedList = append(dataSet.orderedList, testCase)

	testCase = "failed_query"
	data.expected = errors.New("This is a failure test")
	mock.ExpectBegin()
	mock.ExpectExec(AddAssetQuery).WithArgs(UserAssets, asset.ID, binary).WillReturnError(data.expected.(error))
	mock.ExpectRollback()
	dataSet.testDataSet[testCase] = data
	dataSet.orderedList = append(dataSet.orderedList, testCase)

	dbConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}
	Functions = MYSQLFunctions{}

	return &dataSet, dbConnector, nil
}

func TestAddAsset(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createAssetTestData()
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.orderedList {
		tx, err := dbConnector.DB.Begin()
		if err != nil {
			t.Errorf("Failed to setup DB transaction: %s", err)
			return
		}
		testCase := dataSet.testDataSet[testCaseString]
		asset := testCase.data.(models.Asset)

		err = Functions.AddAsset(UserAssets, &asset, tx)
		if !test.ErrEqual(err, testCase.expected) {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, testCase.expected)
			return
		}
	}
}
