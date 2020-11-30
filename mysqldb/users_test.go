package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	GetUserByEmailTest = 0
	AddUserTest        = 1
)

func createUsersTestData(testID int) (*test.OrderedTests, DBConnectorMock, error) {
	dbConnector := DBConnectorMock{}
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}
	data := test.Data{}

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

	switch testID {

	case GetUserByEmailTest:
		testCase := "valid_email"
		data = test.Data{
			Data:     user.Email,
			Expected: make(map[string]interface{}),
		}
		data.Expected = make(map[string]interface{})
		data.Expected.(map[string]interface{})["data"] = &user
		data.Expected.(map[string]interface{})["error"] = nil
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "user_settings_id", "user_assets_id"}).
			AddRow(binaryUserID, user.Name, user.Email, user.Password, binarySettingsID, binaryAssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		data = test.Data{
			Data:     user.Email,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnError(data.Expected.(map[string]interface{})["error"].(error))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_email"
		data = test.Data{
			Data:     user.Email,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case AddUserTest:
		testCase := "valid_user"
		password := []byte{}
		data = test.Data{
			Data:     user,
			Expected: nil,
		}
		data.Data = user
		data.Expected = nil
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "duplicate_name"
		data = test.Data{
			Data:     user,
			Expected: fmt.Errorf(ErrSQLDuplicateUserNameEntryString, user.Name),
		}
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnError(data.Expected.(error))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "duplicate_email"
		data = test.Data{
			Data:     user,
			Expected: fmt.Errorf(ErrSQLDuplicateEmailEntryString, user.Email),
		}
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnError(data.Expected.(error))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	default:
		return nil, dbConnector, fmt.Errorf("Unknown test %d", testID)
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
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := dbConnector.DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			email := testCase.Data.(string)
			var expectedData *models.User
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.User)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}

			output, err := Functions.GetUserByEmail(email, tx)
			if !cmp.Equal(output, expectedData) {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData)
				return
			}

			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}
		})
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
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := dbConnector.DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			user := testCase.Data.(models.User)
			var expectedError error
			if testCase.Expected != nil {
				expectedError = testCase.Expected.(error)
			}

			err = Functions.AddUser(&user, tx)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}
		})
	}
}
