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

func createSettingsTestData() (*test.OrderedTests, DBConnectorMock, error) {
	dbConnector := DBConnectorMock{}
	dataSet := test.OrderedTests{
		orderedList: make(test.OrderedTestList, 0),
		testDataSet: make(test.DataSet, 0),
	}
	settings := make(models.Settings)

	settingID, err := uuid.NewUUID()
	if err != nil {
		return nil, dbConnector, err
	}

	setting := models.UserSettings{
		ID:       settingID,
		Settings: settings,
	}

	binary, err := json.Marshal(setting.Settings)
	if err != nil {
		return nil, dbConnector, err
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, dbConnector, err
	}

	data := test.TestData{
		data:     setting,
		expected: nil,
	}

	testCase := "valid_user_settings"
	mock.ExpectBegin()
	mock.ExpectExec(AddUserSettingsQuery).WithArgs(setting.ID, binary).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	dataSet.testDataSet[testCase] = data
	dataSet.orderedList = append(dataSet.orderedList, testCase)

	testCase = "failed_query"
	data.expected = errors.New("This is a failure test")
	mock.ExpectBegin()
	mock.ExpectExec(AddUserSettingsQuery).WithArgs(setting.ID, binary).WillReturnError(data.expected.(error))
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

func TestAddSettings_ValidUserSetting(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createSettingsTestData()
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
		userSettings := testCase.data.(models.UserSettings)

		err = Functions.AddSettings(&userSettings, tx)
		if !test.ErrEqual(err, testCase.expected) {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, testCase.expected)
			return
		}
	}
}
