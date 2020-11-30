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
	GetUserTest    = 0
	AddUserTest    = 1
	DeleteUserTest = 2
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

	case GetUserTest:
		testCase := "valid_email"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["key_name"] = ByEmail
		data.Data.(map[string]interface{})["key_value"] = user.Email
		data.Expected.(map[string]interface{})["data"] = &user
		data.Expected.(map[string]interface{})["error"] = nil
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "user_settings_id", "user_assets_id"}).
			AddRow(binaryUserID, user.Name, user.Email, user.Password, binarySettingsID, binaryAssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserQuery).WithArgs(ByEmail, user.Email).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query_email"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["key_name"] = ByEmail
		data.Data.(map[string]interface{})["key_value"] = user.Email
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserQuery).WithArgs(ByEmail, user.Email).WillReturnError(data.Expected.(map[string]interface{})["error"].(error))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_email"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["key_name"] = ByEmail
		data.Data.(map[string]interface{})["key_value"] = user.Email
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserQuery).WithArgs(ByEmail, user.Email).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "valid_ID"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["key_name"] = ByID
		data.Data.(map[string]interface{})["key_value"] = user.ID
		data.Expected = make(map[string]interface{})
		data.Expected.(map[string]interface{})["data"] = &user
		data.Expected.(map[string]interface{})["error"] = nil
		rows = sqlmock.NewRows([]string{"id", "name", "email", "password", "user_settings_id", "user_assets_id"}).
			AddRow(binaryUserID, user.Name, user.Email, user.Password, binarySettingsID, binaryAssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserQuery).WithArgs(ByID, user.ID).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query_ID"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["key_name"] = ByID
		data.Data.(map[string]interface{})["key_value"] = user.ID
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserQuery).WithArgs(ByID, user.ID).WillReturnError(data.Expected.(map[string]interface{})["error"].(error))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_ID"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["key_name"] = ByID
		data.Data.(map[string]interface{})["key_value"] = user.ID
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserQuery).WithArgs(ByID, user.ID).WillReturnError(sql.ErrNoRows)
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

	case DeleteUserTest:
		testCase := "valid_user"
		data = test.Data{
			Data:     user.ID,
			Expected: nil,
		}
		mock.ExpectBegin()
		mock.ExpectExec(DeleteUserQuery).WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_user"
		data = test.Data{
			Data:     user.ID,
			Expected: ErrNoUserDeleted,
		}
		mock.ExpectBegin()
		mock.ExpectExec(DeleteUserQuery).WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 0))
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

func TestGetUser(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createUsersTestData(GetUserTest)
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
			keyValue := testCase.Data.(map[string]interface{})["key_value"]
			keyName := testCase.Data.(map[string]interface{})["key_name"].(string)
			var expectedData *models.User
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.User)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}

			output, err := Functions.GetUser(keyName, keyValue, tx)
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

func TestDeleteUser(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createUsersTestData(DeleteUserTest)
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
			userID := testCase.Data.(uuid.UUID)
			var expectedError error
			if testCase.Expected != nil {
				expectedError = testCase.Expected.(error)
			}

			err = Functions.DeleteUser(&userID, tx)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}
		})
	}
}
