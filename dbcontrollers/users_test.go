package dbcontrollers

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
	CreateUserTest = 0
	DeleteUserTest = 1
)

func createUserTestData(testID int) (*test.OrderedTests, error) {
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	dbController = &MYSQLController{}

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

	models.Interface = ModelInterfaceMock{
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

	usersProducts := models.UserProductIDs{
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
		data := test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["input"] = user
		data.Data.(map[string]interface{})["db_mock"] = nil
		data.Expected.(map[string]interface{})["data"] = &user
		data.Expected.(map[string]interface{})["error"] = nil
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "existing_user"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["input"] = user
		data.Data.(map[string]interface{})["db_mock"] = &user
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = ErrDuplicateEmailEntry
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case DeleteUserTest:

		testCase := "valid_data_has_nominee"
		data := test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		usersProducts.ProductMap[productID] = 0
		nominees := make(map[uuid.UUID]uuid.UUID)
		nominees[productID] = nomineeID
		data.Data.(map[string]interface{})["user_id"] = user.ID
		data.Data.(map[string]interface{})["nominees"] = nominees
		data.Data.(map[string]interface{})["db_mock_product"] = &product
		data.Data.(map[string]interface{})["db_mock_user"] = &user
		data.Data.(map[string]interface{})["db_mock_users_products"] = usersProducts
		data.Data.(map[string]interface{})["db_mock_privileges"] = privileges
		data.Expected.(map[string]interface{})["error"] = nil
		data.Expected.(map[string]interface{})["user_deleted"] = true
		data.Expected.(map[string]interface{})["product_deleted"] = false
		data.Expected.(map[string]interface{})["users_products_updated"] = true
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "valid_data_has_no_nominee"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		usersProducts.ProductMap[productID] = 0
		data.Data.(map[string]interface{})["user_id"] = user.ID
		data.Data.(map[string]interface{})["nominees"] = nil
		data.Data.(map[string]interface{})["db_mock_product"] = &product
		data.Data.(map[string]interface{})["db_mock_user"] = &user
		data.Data.(map[string]interface{})["db_mock_users_products"] = usersProducts
		data.Data.(map[string]interface{})["db_mock_privileges"] = privileges
		data.Expected.(map[string]interface{})["error"] = nil
		data.Expected.(map[string]interface{})["user_deleted"] = true
		data.Expected.(map[string]interface{})["product_deleted"] = true
		data.Expected.(map[string]interface{})["users_products_updated"] = false
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_user"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["user_id"] = user.ID
		data.Data.(map[string]interface{})["nominees"] = nil
		data.Data.(map[string]interface{})["db_mock_product"] = &product
		data.Data.(map[string]interface{})["db_mock_user"] = &user
		data.Data.(map[string]interface{})["db_mock_users_products"] = usersProducts
		data.Data.(map[string]interface{})["db_mock_error"] = sql.ErrNoRows
		data.Data.(map[string]interface{})["db_mock_privileges"] = privileges
		data.Expected.(map[string]interface{})["error"] = ErrUserNotFound
		data.Expected.(map[string]interface{})["user_deleted"] = false
		data.Expected.(map[string]interface{})["product_deleted"] = false
		data.Expected.(map[string]interface{})["users_products_updated"] = false
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "has_no_products"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		usersProducts = models.UserProductIDs{
			ProductMap: make(map[uuid.UUID]int),
		}
		data.Data.(map[string]interface{})["user_id"] = user.ID
		data.Data.(map[string]interface{})["nominees"] = nil
		data.Data.(map[string]interface{})["db_mock_product"] = &product
		data.Data.(map[string]interface{})["db_mock_user"] = &user
		data.Data.(map[string]interface{})["db_mock_users_products"] = usersProducts
		data.Data.(map[string]interface{})["db_mock_error"] = nil
		data.Data.(map[string]interface{})["db_mock_privileges"] = privileges
		data.Expected.(map[string]interface{})["error"] = nil
		data.Expected.(map[string]interface{})["user_deleted"] = true
		data.Expected.(map[string]interface{})["product_deleted"] = false
		data.Expected.(map[string]interface{})["users_products_updated"] = false
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	}

	mysqldb.Functions = DBFunctionInterfaceMock{}
	mysqldb.DBConnector = DBConnectorMock{}
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
			mysqldb.Functions = DBFunctionInterfaceMock{
				user:         mockCopy,
				userAdded:    test.NewBool(false),
				productAdded: test.NewBool(false),
			}

			output, err := dbController.CreateUser(
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
			dbMockUsersProducts := testCase.Data.(map[string]interface{})["db_mock_users_products"].(models.UserProductIDs)

			var dbMockError error
			if testCase.Data.(map[string]interface{})["db_mock_error"] != nil {
				dbMockError = testCase.Data.(map[string]interface{})["db_mock_error"].(error)
			}

			mysqldb.Functions = DBFunctionInterfaceMock{
				user:                 dbMockUser,
				product:              dbMockProduct,
				userProducts:         &dbMockUsersProducts,
				err:                  dbMockError,
				userDeleted:          test.NewBool(false),
				productDeleted:       test.NewBool(false),
				usersProductsUpdated: test.NewBool(false),
				privileges:           testCase.Data.(map[string]interface{})["db_mock_privileges"].(models.Privileges),
			}

			err := dbController.DeleteUser(&userID, nominatedOwners)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}

			if !cmp.Equal(*mysqldb.Functions.(DBFunctionInterfaceMock).userDeleted, expectedUserDeleted) {
				t.Errorf(test.TestResultString, testCaseString, *mysqldb.Functions.(DBFunctionInterfaceMock).userDeleted, expectedUserDeleted)
				return
			}

			if !cmp.Equal(*mysqldb.Functions.(DBFunctionInterfaceMock).productDeleted, expectedProductDeleted) {
				t.Errorf(test.TestResultString, testCaseString, *mysqldb.Functions.(DBFunctionInterfaceMock).productDeleted, expectedProductDeleted)
				return
			}

			if !cmp.Equal(*mysqldb.Functions.(DBFunctionInterfaceMock).usersProductsUpdated, expectedUsersProducts) {
				t.Errorf(test.TestResultString, testCaseString, *mysqldb.Functions.(DBFunctionInterfaceMock).usersProductsUpdated, expectedUsersProducts)
				return
			}
		})
	}
}
