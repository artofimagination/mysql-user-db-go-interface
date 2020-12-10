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
	GetUserTest = iota
	AddUserTest
	DeleteUserTest
	GetProductUserIDsTest
	DeleteProductUserTest
)

type UserExpectedTestData struct {
	err error
}

type UserInputTestData struct {
	productID uuid.UUID
	user      models.User
}

func createUsersTestData(testID int) (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	settingsID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	assetsID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	binaryUserID, err := json.Marshal(user.ID)
	if err != nil {
		return nil, err
	}

	binarySettingsID, err := json.Marshal(user.SettingsID)
	if err != nil {
		return nil, err
	}

	binaryAssetsID, err := json.Marshal(user.AssetsID)
	if err != nil {
		return nil, err
	}

	productUsers, err := createTestProductUsersData()
	if err != nil {
		return nil, err
	}

	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	switch testID {

	case GetUserTest:
		testCase := "valid_email"
		data := make(map[string]interface{})
		expected := make(map[string]interface{})
		data["query_type"] = ByEmail
		data["key_value"] = user.Email
		expected["data"] = &user
		expected["error"] = nil
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "user_settings_id", "user_assets_id"}).
			AddRow(binaryUserID, user.Name, user.Email, user.Password, binarySettingsID, binaryAssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query_email"
		data = make(map[string]interface{})
		expected = make(map[string]interface{})
		data["query_type"] = ByEmail
		data["key_value"] = user.Email
		expected["data"] = nil
		expected["error"] = errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnError(expected["error"].(error))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_email"
		data = make(map[string]interface{})
		expected = make(map[string]interface{})
		data["query_type"] = ByEmail
		data["key_value"] = user.Email
		expected["data"] = nil
		expected["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "valid_ID"
		data = make(map[string]interface{})
		data["query_type"] = ByID
		data["key_value"] = user.ID
		expected = make(map[string]interface{})
		expected["data"] = &user
		expected["error"] = nil
		rows = sqlmock.NewRows([]string{"id", "name", "email", "password", "user_settings_id", "user_assets_id"}).
			AddRow(binaryUserID, user.Name, user.Email, user.Password, binarySettingsID, binaryAssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByIDQuery).WithArgs(user.ID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query_ID"
		data = make(map[string]interface{})
		expected = make(map[string]interface{})
		data["query_type"] = ByID
		data["key_value"] = user.ID
		expected["data"] = nil
		expected["error"] = errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByIDQuery).WithArgs(user.ID).WillReturnError(expected["error"].(error))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_ID"
		data = make(map[string]interface{})
		expected = make(map[string]interface{})
		data["query_type"] = ByID
		data["key_value"] = user.ID
		expected["data"] = nil
		expected["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByIDQuery).WithArgs(user.ID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case AddUserTest:
		testCase := "valid_user"
		password := ""
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     user,
			Expected: nil,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "duplicate_name"
		expected := fmt.Errorf(ErrSQLDuplicateUserNameEntryString, user.Name)
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnError(expected)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     user,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "duplicate_email"
		expected = fmt.Errorf(ErrSQLDuplicateEmailEntryString, user.Email)
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnError(expected)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     user,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case DeleteUserTest:
		testCase := "valid_user"
		mock.ExpectBegin()
		mock.ExpectExec(DeleteUserQuery).WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     user.ID,
			Expected: nil,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_user"

		mock.ExpectBegin()
		mock.ExpectExec(DeleteUserQuery).WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     user.ID,
			Expected: ErrNoUserDeleted,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetProductUserIDsTest:

		testCase := "valid_id"
		expected := make(map[string]interface{})
		expected["data"] = productUsers
		expected["error"] = nil
		rows := sqlmock.NewRows([]string{"products_id", "privilege"})
		for _, userID := range productUsers.UserIDArray {
			rows.AddRow(userID, productUsers.UserMap[userID])
		}

		mock.ExpectBegin()
		mock.ExpectQuery(GetProductUserIDsQuery).WithArgs(userID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     userID,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_users"
		expected = make(map[string]interface{})
		expected["data"] = nil
		expected["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductUserIDsQuery).WithArgs(userID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     userID,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case DeleteProductUserTest:
		testCase := "valid_ids"
		expected := UserExpectedTestData{
			err: nil,
		}
		input := UserInputTestData{
			user:      user,
			productID: productID,
		}

		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductUserQuery).WithArgs(productID, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		dataSet.TestDataSet[testCase] = test.Data{
			Data:     input,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	default:
		return nil, fmt.Errorf("Unknown test %d", testID)
	}

	DBConnector = &DBConnectorMock{
		DB:   db,
		Mock: mock,
	}
	Functions = &MYSQLFunctions{}

	return dataSet, nil
}

func TestGetUser(t *testing.T) {
	// Create test data
	dataSet, err := createUsersTestData(GetUserTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBConnector.(*DBConnectorMock).DB.Close()

	// Run test
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			keyValue := testCase.Data.(map[string]interface{})["key_value"]
			query := testCase.Data.(map[string]interface{})["query_type"].(int)
			var expectedData *models.User
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.User)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}

			output, err := Functions.GetUser(query, keyValue, tx)
			if diff := pretty.Diff(output, expectedData); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData)
				return
			}

			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}
		})
	}
}

func TestAddUser(t *testing.T) {
	// Create test data
	dataSet, err := createUsersTestData(AddUserTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBConnector.(*DBConnectorMock).DB.Begin()
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
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	// Create test data
	dataSet, err := createUsersTestData(DeleteUserTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBConnector.(*DBConnectorMock).DB.Begin()
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
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}
		})
	}
}

func TestGetProductUserIDs(t *testing.T) {
	// Create test data
	dataSet, err := createUsersTestData(GetProductUserIDsTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			productID := testCase.Data.(uuid.UUID)
			var expectedData *models.ProductUserIDs
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.ProductUserIDs)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}

			output, err := Functions.GetProductUserIDs(&productID, tx)
			if diff := pretty.Diff(output, expectedData); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData)
				return
			}

			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}
		})
	}
}

func TestDeleteProductUser(t *testing.T) {
	// Create test data
	dataSet, err := createUsersTestData(DeleteProductUserTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			expected := dataSet.TestDataSet[testCaseString].Expected.(UserExpectedTestData)
			input := dataSet.TestDataSet[testCaseString].Data.(UserInputTestData)

			err = Functions.DeleteProductUser(&input.productID, &input.user.ID, tx)
			if diff := pretty.Diff(err, expected.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expected.err)
				return
			}
		})
	}
}
