package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/tests"
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

type AssetExpectedData struct {
	asset *models.Asset
	err   error
}

type AssetInputData struct {
	asset *models.Asset
}

func createAssetTestData(testID int) (*tests.OrderedTests, error) {
	dataSet := &tests.OrderedTests{
		OrderedList: make(tests.OrderedTestList, 0),
		TestDataSet: make(tests.DataSet),
	}
	references := make(models.DataMap)

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	asset := &models.Asset{
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

	switch testID {
	case AddAssetTest:
		testCase := "valid_user_asset"
		mock.ExpectBegin()
		query := fmt.Sprintf(AddAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(asset.ID, binaryDataMap).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: AssetInputData{
				asset: asset,
			},
			Mock: nil,
			Expected: AssetExpectedData{
				err: nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		err := errors.New("This is a failure test")
		mock.ExpectBegin()
		query = fmt.Sprintf(AddAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(asset.ID, binaryDataMap).WillReturnError(err)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: AssetInputData{
				asset: asset,
			},
			Mock: nil,
			Expected: AssetExpectedData{
				err: err,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case DeleteAssetTest:
		testCase := "valid_user_asset"
		mock.ExpectBegin()
		query := fmt.Sprintf(DeleteAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(asset.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: AssetInputData{
				asset: asset,
			},
			Mock: nil,
			Expected: AssetExpectedData{
				err: nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		mock.ExpectBegin()
		query = fmt.Sprintf(DeleteAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(asset.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: AssetInputData{
				asset: asset,
			},
			Mock: nil,
			Expected: AssetExpectedData{
				err: fmt.Errorf(ErrAssetMissing, UserAssets),
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case UpdateAssetTest:
		testCase := "valid_asset"

		mock.ExpectBegin()
		query := fmt.Sprintf(UpdateAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(binaryDataMap, &asset.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: AssetInputData{
				asset: asset,
			},
			Mock: nil,
			Expected: AssetExpectedData{
				err: nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_asset"
		mock.ExpectBegin()
		query = fmt.Sprintf(UpdateAssetQuery, UserAssets)
		mock.ExpectExec(query).WithArgs(binaryDataMap, &asset.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: AssetInputData{
				asset: asset,
			},
			Mock: nil,
			Expected: AssetExpectedData{
				err: fmt.Errorf(ErrAssetMissing, UserAssets),
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetAssetTest:
		testCase := "valid_asset"

		rows := sqlmock.NewRows([]string{"id", "data"}).AddRow(binaryID, binaryDataMap)
		mock.ExpectBegin()
		query := fmt.Sprintf(GetAssetQuery, UserAssets)
		mock.ExpectQuery(query).WithArgs(&asset.ID).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: AssetInputData{
				asset: asset,
			},
			Mock: nil,
			Expected: AssetExpectedData{
				asset: asset,
				err:   nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_asset"
		mock.ExpectBegin()
		query = fmt.Sprintf(GetAssetQuery, UserAssets)
		mock.ExpectQuery(query).WithArgs(&asset.ID).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: AssetInputData{
				asset: asset,
			},
			Mock: nil,
			Expected: AssetExpectedData{
				asset: nil,
				err:   sql.ErrNoRows,
			},
		}
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
			expectedData := testCase.Expected.(AssetExpectedData)
			inputData := testCase.Data.(AssetInputData)

			err = DBFunctions.AddAsset(UserAssets, inputData.asset, tx)
			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(tests.TestResultString, testCaseString, err, expectedData.err, diff)
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
			expectedData := testCase.Expected.(AssetExpectedData)
			inputData := testCase.Data.(AssetInputData)

			err = DBFunctions.DeleteAsset(UserAssets, &inputData.asset.ID, tx)
			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(tests.TestResultString, testCaseString, err, expectedData.err, diff)
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
			expectedData := testCase.Expected.(AssetExpectedData)
			inputData := testCase.Data.(AssetInputData)

			err = DBFunctions.UpdateAsset(UserAssets, inputData.asset)
			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(tests.TestResultString, testCaseString, err, expectedData.err, diff)
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
			expectedData := testCase.Expected.(AssetExpectedData)
			inputData := testCase.Data.(AssetInputData)

			output, err := DBFunctions.GetAsset(UserAssets, &inputData.asset.ID)
			tests.CheckResult(output, expectedData.asset, err, expectedData.err, testCaseString, t)
		})
	}
}
