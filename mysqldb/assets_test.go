package mysqldb

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	AddAssetTest    = 0
	DeleteAssetTest = 1
)

func createAssetTestData(testID int) (*test.OrderedTests, DBConnectorMock, error) {
	dbConnector := DBConnectorMock{}
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}
	references := make(models.DataMap)

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, dbConnector, err
	}

	asset := models.Asset{
		ID:      assetID,
		DataMap: references,
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, dbConnector, err
	}

	binary, err := json.Marshal(asset.DataMap)
	if err != nil {
		return nil, dbConnector, err
	}

	data := test.Data{
		Data:     asset,
		Expected: nil,
	}

	switch testID {
	case AddAssetTest:
		testCase := "valid_user_asset"
		mock.ExpectBegin()
		mock.ExpectExec(AddAssetQuery).WithArgs(UserAssets, asset.ID, binary).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		data.Expected = errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectExec(AddAssetQuery).WithArgs(UserAssets, asset.ID, binary).WillReturnError(data.Expected.(error))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case DeleteAssetTest:
		testCase := "valid_user_asset"
		data.Expected = nil
		mock.ExpectBegin()
		mock.ExpectExec(DeleteAssetQuery).WithArgs(UserAssets, asset.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		data.Expected = fmt.Errorf(ErrAssetMissing, UserAssets)
		mock.ExpectBegin()
		mock.ExpectExec(DeleteAssetQuery).WithArgs(UserAssets, asset.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	}

	dbConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}
	Functions = MYSQLFunctions{}

	return &dataSet, dbConnector, nil
}

func TestAddAsset(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createAssetTestData(AddAssetTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := dbConnector.DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction: %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			var expectedError error
			if testCase.Expected != nil {
				expectedError = testCase.Expected.(error)
			}
			asset := testCase.Data.(models.Asset)

			err = Functions.AddAsset(UserAssets, &asset, tx)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, testCase.Expected)
				return
			}
		})
	}
}

func TestDeleteAsset(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createAssetTestData(DeleteAssetTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := dbConnector.DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction: %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			var expectedError error
			if testCase.Expected != nil {
				expectedError = testCase.Expected.(error)
			}
			asset := testCase.Data.(models.Asset)

			err = Functions.DeleteAsset(UserAssets, &asset.ID, tx)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, testCase.Expected)
				return
			}
		})
	}
}
