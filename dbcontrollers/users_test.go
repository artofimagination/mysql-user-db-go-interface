package dbcontrollers

import (
	"database/sql"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

const (
	CreateUserTest = iota
	DeleteUserTest
)

func createUserTestData(testID int) (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
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

	dataMap := make(models.DataMap)
	assets := &models.Asset{
		ID:      assetID,
		DataMap: dataMap,
	}

	models.Interface = &ModelInterfaceMock{
		assetID:    assetID,
		settingsID: settingsID,
		userID:     userID,
		asset:      assets,
	}

	userData := models.UserData{
		ID:       user.ID,
		Name:     user.Name,
		Settings: assets,
		Assets:   assets,
	}

	privileges := make(models.Privileges, 2)
	privilege := &models.Privilege{
		ID:          0,
		Name:        "Owner",
		Description: "description0",
	}
	privileges[0] = privilege
	privilege = &models.Privilege{
		ID:          1,
		Name:        "User",
		Description: "description1",
	}
	privileges[1] = privilege

	switch testID {
	case CreateUserTest:
		testCase := "no_existing_user"
		data := make(map[string]interface{})
		data["input"] = userData
		data["password"] = user.Password
		data["db_mock"] = nil
		expected := make(map[string]interface{})
		expected["data"] = &userData
		expected["error"] = nil
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "existing_user"
		data = make(map[string]interface{})
		data["input"] = userData
		data["password"] = user.Password
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
		usersProducts = models.UserProductIDs{
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
	return dataSet, nil
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
			var expectedData *models.UserData
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.UserData)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}
			input := testCase.Data.(map[string]interface{})["input"].(models.UserData)
			var DBMock *models.User
			if testCase.Data.(map[string]interface{})["db_mock"] != nil {
				DBMock = testCase.Data.(map[string]interface{})["db_mock"].(*models.User)
			}

			password := testCase.Data.(map[string]interface{})["password"].([]byte)

			mockCopy := DBMock
			mysqldb.Functions = &DBFunctionInterfaceMock{
				user:         mockCopy,
				userAdded:    false,
				productAdded: false,
			}

			output, err := dbController.CreateUser(
				input.Name,
				input.Email,
				password,
				func(*uuid.UUID) string {
					return "testPath"
				}, func([]byte) ([]byte, error) {
					return []byte{}, nil
				})

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

			err := dbController.DeleteUser(&userID, nominatedOwners)
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
				return
			}

			if diff := pretty.Diff(mysqldb.Functions.(*DBFunctionInterfaceMock).userDeleted, expectedUserDeleted); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, mysqldb.Functions.(*DBFunctionInterfaceMock).userDeleted, expectedUserDeleted)
				return
			}

			if diff := pretty.Diff(mysqldb.Functions.(*DBFunctionInterfaceMock).productDeleted, expectedProductDeleted); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, mysqldb.Functions.(*DBFunctionInterfaceMock).productDeleted, expectedProductDeleted)
				return
			}

			if diff := pretty.Diff(mysqldb.Functions.(*DBFunctionInterfaceMock).usersProductsUpdated, expectedUsersProducts); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, mysqldb.Functions.(*DBFunctionInterfaceMock).usersProductsUpdated, expectedUsersProducts)
				return
			}
		})
	}
}
