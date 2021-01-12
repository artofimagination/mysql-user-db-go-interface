package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
)

const (
	AddAssetTest = iota
	DeleteAssetTest
	UpdateAssetTest
	GetAssetTest
)

func createAssetTestData(testID int) (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}
	references := make(models.DataMap)

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	asset := models.Asset{
		ID:      assetID,
		DataMap: references,
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, err
	}

	binaryDataMap, err := json.Marshal(asset.DataMap)
	if err != nil {
		return nil, err
	}

	binaryID, err := json.Marshal(asset.ID)
	if err != nil {
		return nil, err
	}

	data := test.Data{
		Data:     asset,
		Expected: nil,
	}

	switch testID {
	case AddAssetTest:
		testCase := "valid_user_asset"
		mock.ExpectBegin()
		query := fmt.Sprintf(AddAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(asset.ID, binaryDataMap).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		data.Expected = errors.New("This is a failure test")
		mock.ExpectBegin()
		query = fmt.Sprintf(AddAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(asset.ID, binaryDataMap).WillReturnError(data.Expected.(error))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case DeleteAssetTest:
		testCase := "valid_user_asset"
		data.Expected = nil
		mock.ExpectBegin()
		query := fmt.Sprintf(DeleteAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(asset.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		data.Expected = fmt.Errorf(ErrAssetMissing, UserAssets)
		mock.ExpectBegin()
		query = fmt.Sprintf(DeleteAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(asset.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case UpdateAssetTest:
		testCase := "valid_asset"

		data.Expected = nil
		mock.ExpectBegin()
		query := fmt.Sprintf(UpdateAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(binaryDataMap, &asset.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_asset"
		data.Expected = fmt.Errorf(ErrAssetMissing, UserAssets)
		mock.ExpectBegin()
		query = fmt.Sprintf(UpdateAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(binaryDataMap, &asset.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetAssetTest:
		testCase := "valid_asset"
		data := test.Data{
			Data:     asset,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = &asset
		data.Expected.(map[string]interface{})["error"] = nil
		rows := sqlmock.NewRows([]string{"id", "data"}).AddRow(binaryID, binaryDataMap)
		mock.ExpectBegin()
		query := fmt.Sprintf(GetAssetQuery, UserAssets)
		mock.ExpectQuery(query).WithArgs(&asset.ID).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_asset"
		data = test.Data{
			Data:     asset,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		query = fmt.Sprintf(GetAssetQuery, UserAssets)
		mock.ExpectQuery(query).WithArgs(&asset.ID).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	}

	DBFunctions = &MYSQLFunctions{
		DBConnector: &DBConnectorMock{
			DB:   db,
			Mock: mock,
		},
	}

	return dataSet, nil
}

func TestAddAsset(t *testing.T) {
	// Create test data
	dataSet, err := createAssetTestData(AddAssetTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
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

			err = DBFunctions.AddAsset(UserAssets, &asset, tx)
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
				return
			}
		})
	}
}

func TestDeleteAsset(t *testing.T) {
	// Create test data
	dataSet, err := createAssetTestData(DeleteAssetTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
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

			err = DBFunctions.DeleteAsset(UserAssets, &asset.ID, tx)
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
				return
			}
		})
	}
}

func TestUpdateAsset(t *testing.T) {
	// Create test data
	dataSet, err := createAssetTestData(UpdateAssetTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			var expectedError error
			if testCase.Expected != nil {
				expectedError = testCase.Expected.(error)
			}
			asset := testCase.Data.(models.Asset)

			err = DBFunctions.UpdateAsset(UserAssets, &asset)
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
				return
			}
		})
	}
}

func TestGetAsset(t *testing.T) {
	// Create test data
	dataSet, err := createAssetTestData(GetAssetTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}
			var expectedData *models.Asset
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.Asset)
			}
			asset := testCase.Data.(models.Asset)

			output, err := DBFunctions.GetAsset(UserAssets, &asset.ID)
			if diff := pretty.Diff(output, expectedData); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData, diff)
				return
			}

			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
				return
			}
		})
	}
}
