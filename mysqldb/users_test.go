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

type UserExpectedData struct {
	user         *models.User
	productUsers *models.ProductUserIDs
	err          error
}

type UserInputData struct {
	productID *uuid.UUID
	user      *models.User
	queryType int
	keyValue  interface{}
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

	user := &models.User{
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
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "user_settings_id", "user_assets_id"}).
			AddRow(binaryUserID, user.Name, user.Email, user.Password, binarySettingsID, binaryAssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				queryType: ByEmail,
				keyValue:  user.Email,
			},
			Expected: UserExpectedData{
				user: user,
				err:  nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query_email"
		err := errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnError(err)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				queryType: ByEmail,
				keyValue:  user.Email,
			},
			Expected: UserExpectedData{
				user: nil,
				err:  err,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_email"
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				queryType: ByEmail,
				keyValue:  user.Email,
			},
			Expected: UserExpectedData{
				user: nil,
				err:  sql.ErrNoRows,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "valid_ID"
		rows = sqlmock.NewRows([]string{"id", "name", "email", "password", "user_settings_id", "user_assets_id"}).
			AddRow(binaryUserID, user.Name, user.Email, user.Password, binarySettingsID, binaryAssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByIDQuery).WithArgs(user.ID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				queryType: ByID,
				keyValue:  user.ID,
			},
			Expected: UserExpectedData{
				user: user,
				err:  nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query_ID"
		err = errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByIDQuery).WithArgs(user.ID).WillReturnError(err)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				queryType: ByID,
				keyValue:  user.ID,
			},
			Expected: UserExpectedData{
				user: nil,
				err:  err,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_ID"
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserByIDQuery).WithArgs(user.ID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				queryType: ByID,
				keyValue:  user.ID,
			},
			Expected: UserExpectedData{
				user: nil,
				err:  sql.ErrNoRows,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case AddUserTest:
		testCase := "valid_user"
		password := ""
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				user: user,
			},
			Expected: UserExpectedData{
				err: nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "duplicate_name"
		expected := fmt.Errorf(ErrSQLDuplicateUserNameEntryString, user.Name)
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnError(expected)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				user: user,
			},
			Expected: UserExpectedData{
				err: expected,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "duplicate_email"
		expected = fmt.Errorf(ErrSQLDuplicateEmailEntryString, user.Email)
		mock.ExpectBegin()
		mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnError(expected)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				user: user,
			},
			Expected: UserExpectedData{
				err: expected,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case DeleteUserTest:
		testCase := "valid_user"
		mock.ExpectBegin()
		mock.ExpectExec(DeleteUserQuery).WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				user: user,
			},
			Expected: UserExpectedData{
				err: nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_user"
		mock.ExpectBegin()
		mock.ExpectExec(DeleteUserQuery).WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				user: user,
			},
			Expected: UserExpectedData{
				err: ErrNoUserDeleted,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case GetProductUserIDsTest:
		testCase := "valid_id"
		rows := sqlmock.NewRows([]string{"products_id", "privilege"})
		for _, userID := range productUsers.UserIDArray {
			rows.AddRow(userID, productUsers.UserMap[userID])
		}
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductUserIDsQuery).WithArgs(productID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				productID: &productID,
			},
			Expected: UserExpectedData{
				productUsers: productUsers,
				err:          nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_users"
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductUserIDsQuery).WithArgs(productID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data: UserInputData{
				productID: &productID,
			},
			Expected: UserExpectedData{
				productUsers: nil,
				err:          sql.ErrNoRows,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case DeleteProductUserTest:
		testCase := "valid_ids"
		expected := UserExpectedData{
			err: nil,
		}
		input := UserInputData{
			user:      user,
			productID: &productID,
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

	DBFunctions = &MYSQLFunctions{
		DBConnector: &DBConnectorMock{
			DB:   db,
			Mock: mock,
		},
	}

	return dataSet, nil
}

func TestGetUser(t *testing.T) {
	// Create test data
	dataSet, err := createUsersTestData(GetUserTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run test
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(UserExpectedData)
			inputData := testCase.Data.(UserInputData)

			output, err := DBFunctions.GetUser(inputData.queryType, inputData.keyValue, tx)
			if diff := pretty.Diff(output, expectedData.user); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData.user, diff)
				return
			}

			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
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
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(UserExpectedData)
			inputData := testCase.Data.(UserInputData)

			err = DBFunctions.AddUser(inputData.user, tx)
			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
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
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(UserExpectedData)
			inputData := testCase.Data.(UserInputData)

			err = DBFunctions.DeleteUser(&inputData.user.ID, tx)
			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
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

	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(UserExpectedData)
			inputData := testCase.Data.(UserInputData)

			output, err := DBFunctions.GetProductUserIDs(inputData.productID, tx)
			if diff := pretty.Diff(output, expectedData.productUsers); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData.productUsers, diff)
				return
			}

			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
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

	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			expectedData := dataSet.TestDataSet[testCaseString].Expected.(UserExpectedData)
			inputData := dataSet.TestDataSet[testCaseString].Data.(UserInputData)

			err = DBFunctions.DeleteProductUser(inputData.productID, &inputData.user.ID, tx)
			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err)
				return
			}
		})
	}
}
