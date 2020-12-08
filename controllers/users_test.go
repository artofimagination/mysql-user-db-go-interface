package controllers

import (
	"database/sql"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

const (
	CreateUserTest = iota
	DeleteUserTest
)

func createUserTestData(testID int) (*test.OrderedTests, error) {
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	settingsID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	models.Interface = &ModelInterfaceMock{
		assetID:    assetID,
		settingsID: settingsID,
		userID:     userID,
	}

	user := models.User{
		Name:       "testName",
		Email:      "testEmail",
		Password:   []byte{},
		ID:         userID,
		SettingsID: assetID,
		AssetsID:   assetID,
	}

	nomineeID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	usersProducts := models.UserProducts{
		ProductMap: make(map[uuid.UUID]int),
	}

	product := models.Product{
		ID:     productID,
		Public: true,
	}

	privileges := make(models.Privileges, 2)
	privileges[0].ID = 0
	privileges[0].Name = "Owner"
	privileges[0].Description = "description0"
	privileges[1].ID = 1
	privileges[1].Name = "User"
	privileges[1].Description = "description1"

	switch testID {
	case CreateUserTest:
		testCase := "no_existing_user"
		data := make(map[string]interface{})
		data["input"] = user
		data["db_mock"] = nil
		expected := make(map[string]interface{})
		expected["data"] = &user
		expected["error"] = nil
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "existing_user"
		data = make(map[string]interface{})
		data["input"] = user
		data["db_mock"] = &user
		expected = make(map[string]interface{})
		expected["data"] = nil
		expected["error"] = ErrDuplicateEmailEntry
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case DeleteUserTest:

		testCase := "valid_data_has_nominee"
		usersProducts.ProductMap[productID] = 0
		nominees := make(map[uuid.UUID]uuid.UUID)
		nominees[productID] = nomineeID
		data := make(map[string]interface{})
		data["user_id"] = user.ID
		data["nominees"] = nominees
		data["db_mock_product"] = &product
		data["db_mock_user"] = &user
		data["db_mock_users_products"] = usersProducts
		data["db_mock_privileges"] = privileges
		expected := make(map[string]interface{})
		expected["error"] = nil
		expected["user_deleted"] = true
		expected["product_deleted"] = false
		expected["users_products_updated"] = true
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "valid_data_has_no_nominee"
		usersProducts.ProductMap[productID] = 0
		data = make(map[string]interface{})
		data["user_id"] = user.ID
		data["nominees"] = nil
		data["db_mock_product"] = &product
		data["db_mock_user"] = &user
		data["db_mock_users_products"] = usersProducts
		data["db_mock_privileges"] = privileges
		expected = make(map[string]interface{})
		expected["error"] = nil
		expected["user_deleted"] = true
		expected["product_deleted"] = true
		expected["users_products_updated"] = false
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_user"
		data = make(map[string]interface{})
		data["user_id"] = user.ID
		data["nominees"] = nil
		data["db_mock_product"] = &product
		data["db_mock_user"] = &user
		data["db_mock_users_products"] = usersProducts
		data["db_mock_error"] = sql.ErrNoRows
		data["db_mock_privileges"] = privileges
		expected = make(map[string]interface{})
		expected["error"] = ErrUserNotFound
		expected["user_deleted"] = false
		expected["product_deleted"] = false
		expected["users_products_updated"] = false
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "has_no_products"
		usersProducts = models.UserProducts{
			ProductMap: make(map[uuid.UUID]int),
		}
		data = make(map[string]interface{})
		data["user_id"] = user.ID
		data["nominees"] = nil
		data["db_mock_product"] = &product
		data["db_mock_user"] = &user
		data["db_mock_users_products"] = usersProducts
		data["db_mock_error"] = nil
		data["db_mock_privileges"] = privileges
		expected = make(map[string]interface{})
		expected["error"] = nil
		expected["user_deleted"] = true
		expected["product_deleted"] = false
		expected["users_products_updated"] = false
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	}

	mysqldb.Functions = &DBFunctionInterfaceMock{}
	mysqldb.DBConnector = &DBConnectorMock{}
	projectdb = ProjectDBDummy{}
	return &dataSet, nil
}

func TestCreateUser(t *testing.T) {
	// Create test data
	dataSet, err := createUserTestData(CreateUserTest)
	if err != nil {
		t.Errorf("Failed to create test data %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			var expectedData *models.User
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.User)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}
			input := testCase.Data.(map[string]interface{})["input"].(models.User)
			var DBMock *models.User
			if testCase.Data.(map[string]interface{})["db_mock"] != nil {
				DBMock = testCase.Data.(map[string]interface{})["db_mock"].(*models.User)
			}

			mockCopy := DBMock
			mysqldb.Functions = &DBFunctionInterfaceMock{
				user:         mockCopy,
				userAdded:    false,
				productAdded: false,
			}

			output, err := CreateUser(
				input.Name,
				input.Email,
				input.Password,
				func(*uuid.UUID) string {
					return "testPath"
				}, func([]byte) ([]byte, error) {
					return []byte{}, nil
				})

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

func TestDeleteUser(t *testing.T) {
	// Create test data
	dataSet, err := createUserTestData(DeleteUserTest)
	if err != nil {
		t.Errorf("Failed to create test data %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}

			expectedUserDeleted := testCase.Expected.(map[string]interface{})["user_deleted"].(bool)
			expectedProductDeleted := testCase.Expected.(map[string]interface{})["product_deleted"].(bool)
			expectedUsersProducts := testCase.Expected.(map[string]interface{})["users_products_updated"].(bool)

			userID := testCase.Data.(map[string]interface{})["user_id"].(uuid.UUID)
			var nominatedOwners map[uuid.UUID]uuid.UUID
			if testCase.Data.(map[string]interface{})["nominees"] != nil {
				nominatedOwners = testCase.Data.(map[string]interface{})["nominees"].(map[uuid.UUID]uuid.UUID)
			}

			dbMockUser := testCase.Data.(map[string]interface{})["db_mock_user"].(*models.User)
			dbMockProduct := testCase.Data.(map[string]interface{})["db_mock_product"].(*models.Product)
			dbMockUsersProducts := testCase.Data.(map[string]interface{})["db_mock_users_products"].(models.UserProducts)

			var dbMockError error
			if testCase.Data.(map[string]interface{})["db_mock_error"] != nil {
				dbMockError = testCase.Data.(map[string]interface{})["db_mock_error"].(error)
			}

			mysqldb.Functions = &DBFunctionInterfaceMock{
				user:                 dbMockUser,
				product:              dbMockProduct,
				userProducts:         &dbMockUsersProducts,
				err:                  dbMockError,
				userDeleted:          false,
				productDeleted:       false,
				usersProductsUpdated: false,
				privileges:           testCase.Data.(map[string]interface{})["db_mock_privileges"].(models.Privileges),
			}

			err := DeleteUser(&userID, nominatedOwners)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}

			if !cmp.Equal(mysqldb.Functions.(*DBFunctionInterfaceMock).userDeleted, expectedUserDeleted) {
				t.Errorf(test.TestResultString, testCaseString, mysqldb.Functions.(*DBFunctionInterfaceMock).userDeleted, expectedUserDeleted)
				return
			}

			if !cmp.Equal(mysqldb.Functions.(*DBFunctionInterfaceMock).productDeleted, expectedProductDeleted) {
				t.Errorf(test.TestResultString, testCaseString, mysqldb.Functions.(*DBFunctionInterfaceMock).productDeleted, expectedProductDeleted)
				return
			}

			if !cmp.Equal(mysqldb.Functions.(*DBFunctionInterfaceMock).usersProductsUpdated, expectedUsersProducts) {
				t.Errorf(test.TestResultString, testCaseString, mysqldb.Functions.(*DBFunctionInterfaceMock).usersProductsUpdated, expectedUsersProducts)
				return
			}
		})
	}
}
