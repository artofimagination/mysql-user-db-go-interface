package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	GetUserByEmailTest = 0
	AddUserTest        = 1
)

func createUsersTestData(test int) (*test.OrderedTests, DBConnectorMock, error) {
	dbConnector := DBConnectorMock{}
	dataSet := test.OrderedTests{
		orderedList: make(test.OrderedTestList, 0),
		testDataSet: make(test.DataSet, 0),
	}
	data := test.TestData{}

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, dbConnector, err
	}

	settingsID, err := uuid.NewUUID()
	if err != nil {
		return nil, dbConnector, err
	}

	assetsID, err := uuid.NewUUID()
	if err != nil {
		return nil, dbConnector, err
	}

	user := models.User{
		ID:         userID,
		Name:       "testName",
		Email:      "test@test.com",
		Password:   []byte{},
		SettingsID: settingsID,
		AssetsID:   assetsID,
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, dbConnector, err
	}

	binaryUserID, err := json.Marshal(user.ID)
	if err != nil {
		return nil, dbConnector, err
	}

	binarySettingsID, err := json.Marshal(user.SettingsID)
	if err != nil {
		return nil, dbConnector, err
	}

	binaryAssetsID, err := json.Marshal(user.AssetsID)
	if err != nil {
		return nil, dbConnector, err
	}

	switch test {

	case GetUserByEmailTest:
		testCase := "valid_email"
		data = test.TestData{
			data:     user.Email,
			expected: make(map[string]interface{}),
		}
		data.expected = make(map[string]interface{})
		data.expected.(map[string]interface{})["data"] = &user
		data.expected.(map[string]interface{})["error"] = nil
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "user_settings_id", "user_assets_id"}).
			AddRow(binaryUserID, user.Name, user.Email, user.Password, binarySettingsID, binaryAssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "failed_query"
		data = test.TestData{
			data:     user.Email,
			expected: make(map[string]interface{}),
		}
		data.expected.(map[string]interface{})["data"] = nil
		data.expected.(map[string]interface{})["error"] = errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnError(data.expected.(map[string]interface{})["error"].(error))
		mock.ExpectRollback()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "invalid_email"
		data = test.TestData{
			data:     user.Email,
			expected: make(map[string]interface{}),
		}
		data.expected.(map[string]interface{})["data"] = nil
		data.expected.(map[string]interface{})["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)
	case AddUserTest:
		testCase := "valid_user"
		password := []byte{}
		data = test.TestData{
			data:     user,
			expected: nil,
		}
		data.data = user
		data.expected = nil
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "duplicate_name"
		data = test.TestData{
			data:     user,
			expected: fmt.Errorf(ErrSQLDuplicateUserNameEntryString, user.Name),
		}
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnError(data.expected.(error))
		mock.ExpectRollback()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "duplicate_email"
		data = test.TestData{
			data:     user,
			expected: fmt.Errorf(ErrSQLDuplicateEmailEntryString, user.Email),
		}
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnError(data.expected.(error))
		mock.ExpectRollback()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)
	default:
		return nil, dbConnector, fmt.Errorf("Unknown test %d", test)
	}

	dbConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}
	Functions = MYSQLFunctions{}

	return &dataSet, dbConnector, nil
}

func TestGetUserByEmail(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createUsersTestData(GetUserByEmailTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run test
	for _, testCaseString := range dataSet.orderedList {
		tx, err := dbConnector.DB.Begin()
		if err != nil {
			t.Errorf("Failed to setup DB transaction %s", err)
			return
		}
		testCase := dataSet.testDataSet[testCaseString]
		email := testCase.data.(string)
		var expectedData *models.User
		if testCase.expected.(map[string]interface{})["data"] != nil {
			expectedData = testCase.expected.(map[string]interface{})["data"].(*models.User)
		}
		var expectedError error
		if testCase.expected.(map[string]interface{})["error"] != nil {
			expectedError = testCase.expected.(map[string]interface{})["error"].(error)
		}

		output, err := Functions.GetUserByEmail(email, tx)
		if !cmp.Equal(output, expectedData) {
			t.Errorf("\n%s test failed.\n  Returned:\n   %+v\n  Expected:\n   %+v", testCaseString, output, expectedData)
			return
		}

		if !test.ErrEqual(err, expectedError) {
			t.Errorf("\n%s test failed.\n  Returned:\n   %+v\n  Expected:\n   %+v", testCaseString, err, expectedError)
			return
		}
	}
}

func TestAddUser(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createUsersTestData(AddUserTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.orderedList {
		tx, err := dbConnector.DB.Begin()
		if err != nil {
			t.Errorf("Failed to setup DB transaction %s", err)
			return
		}
		testCase := dataSet.testDataSet[testCaseString]
		user := testCase.data.(models.User)
		var expectedError error
		if testCase.expected != nil {
			expectedError = testCase.expected.(error)
		}

		err = Functions.AddUser(&user, tx)
		if !test.ErrEqual(err, expectedError) {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, expectedError)
			return
		}
	}
}
