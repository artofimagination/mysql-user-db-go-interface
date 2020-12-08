package dbcontrollers

import (
	"fmt"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
)

func createTestUsersData() (*models.ProductUserIDs, models.Privileges) {
	privileges := make(models.Privileges, 2)
	privileges[0].ID = 0
	privileges[0].Name = "Owner"
	privileges[0].Description = "description0"
	privileges[1].ID = 1
	privileges[1].Name = "User"
	privileges[1].Description = "description1"
	mysqldb.DBConnector = DBConnectorMock{}

	users := models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}

	return &users, privileges
}

func createProductTestData() (*test.OrderedTests, error) {
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	dbController = &MYSQLController{}

	_, privileges := createTestUsersData()

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	models.Interface = ModelInterfaceMock{
		assetID:   assetID,
		productID: productID,
	}

	product := models.Product{
		Name:      "testProduct",
		Public:    true,
		ID:        productID,
		DetailsID: assetID,
		AssetsID:  assetID,
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	dataMap := make(models.DataMap)
	assets := models.Asset{
		ID:      assetID,
		DataMap: dataMap,
	}

	details := models.Asset{
		ID:      assetID,
		DataMap: dataMap,
	}

	productData := models.ProductData{
		ID:      product.ID,
		Name:    product.Name,
		Public:  product.Public,
		Details: details,
		Assets:  assets,
	}

	testCase := "no_existing_product"
	data := test.Data{
		Data:     make(map[string]interface{}),
		Expected: make(map[string]interface{}),
	}

	data.Data.(map[string]interface{})["input"] = productData
	data.Data.(map[string]interface{})["db_mock"] = nil
	data.Data.(map[string]interface{})["user_id"] = userID
	data.Data.(map[string]interface{})["privileges"] = privileges
	data.Expected.(map[string]interface{})["data"] = &productData
	data.Expected.(map[string]interface{})["error"] = nil
	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "existing_product"
	data = test.Data{
		Data:     make(map[string]interface{}),
		Expected: make(map[string]interface{}),
	}

	data.Data.(map[string]interface{})["input"] = productData
	data.Data.(map[string]interface{})["db_mock"] = &product
	data.Data.(map[string]interface{})["user_id"] = userID
	data.Data.(map[string]interface{})["privileges"] = privileges
	data.Expected.(map[string]interface{})["data"] = nil
	data.Expected.(map[string]interface{})["error"] = fmt.Errorf(ErrProductExistsString, product.Name)
	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	mysqldb.Functions = DBFunctionInterfaceMock{}
	mysqldb.DBConnector = DBConnectorMock{}
	return &dataSet, nil
}

func createValidationTestData() (*test.OrderedTests, error) {
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	productUsers, privileges := createTestUsersData()

	testCase := "valid_data"
	data := test.Data{
		Data:     make(map[string]interface{}),
		Expected: nil,
	}

	productUsers.UserMap[userID] = 0
	data.Data.(map[string]interface{})["input"] = productUsers
	data.Data.(map[string]interface{})["privileges"] = privileges

	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "empty_user_list"
	data = test.Data{
		Data:     make(map[string]interface{}),
		Expected: ErrEmptyUsersList,
	}

	productUsers = &models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	data.Data.(map[string]interface{})["input"] = productUsers
	data.Data.(map[string]interface{})["privileges"] = privileges

	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "nil_user_list"
	data = test.Data{
		Data:     make(map[string]interface{}),
		Expected: ErrEmptyUsersList,
	}

	data.Data.(map[string]interface{})["input"] = nil
	data.Data.(map[string]interface{})["privileges"] = privileges

	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "no_owner"
	data = test.Data{
		Data:     make(map[string]interface{}),
		Expected: ErrInvalidOwnerCount,
	}

	productUsers = &models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	productUsers.UserMap[userID] = 1
	data.Data.(map[string]interface{})["input"] = productUsers
	data.Data.(map[string]interface{})["privileges"] = privileges

	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "multiple_owners"
	data = test.Data{
		Data:     make(map[string]interface{}),
		Expected: ErrInvalidOwnerCount,
	}

	productUsers = &models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	productUsers.UserMap[userID] = 0
	userID2, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	productUsers.UserMap[userID2] = 0
	data.Data.(map[string]interface{})["input"] = productUsers
	data.Data.(map[string]interface{})["privileges"] = privileges

	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "invalid_privilege"
	data = test.Data{
		Data:     make(map[string]interface{}),
		Expected: fmt.Errorf(ErrUnknownPrivilegeString, 2, userID.String()),
	}

	productUsers = &models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	productUsers.UserMap[userID] = 2
	data.Data.(map[string]interface{})["input"] = productUsers
	data.Data.(map[string]interface{})["privileges"] = privileges

	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	return &dataSet, nil
}

func TestCreateProduct(t *testing.T) {
	// Create test data
	dataSet, err := createProductTestData()
	if err != nil {
		t.Errorf("Failed to create test data: %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			var expectedData *models.ProductData
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.ProductData)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}
			input := testCase.Data.(map[string]interface{})["input"].(models.ProductData)
			userID := testCase.Data.(map[string]interface{})["user_id"].(uuid.UUID)
			privileges := testCase.Data.(map[string]interface{})["privileges"].(models.Privileges)
			var DBMock *models.Product
			if testCase.Data.(map[string]interface{})["db_mock"] != nil {
				DBMock = testCase.Data.(map[string]interface{})["db_mock"].(*models.Product)
			}

			mockCopy := DBMock
			mysqldb.Functions = DBFunctionInterfaceMock{
				product:      mockCopy,
				privileges:   privileges,
				productAdded: test.NewBool(false),
			}

			output, err := dbController.CreateProduct(
				input.Name,
				input.Public,
				&userID,
				func(*uuid.UUID) (string, error) {
					return "testPath", nil
				})

			if output != nil {
				if output.Name != expectedData.Name || output.Public != expectedData.Public {
					t.Errorf(test.TestResultString, testCaseString, output, expectedData)
					return
				}
			} else {
				if output != expectedData {
					t.Errorf(test.TestResultString, testCaseString, output, expectedData)
					return
				}
			}

			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}
		})
	}
}

func TestValidateUsers(t *testing.T) {
	// Create test data
	dataSet, err := createValidationTestData()
	if err != nil {
		t.Errorf("Failed to create test data: %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			var expectedError error
			if testCase.Expected != nil {
				expectedError = testCase.Expected.(error)
			}
			var input *models.ProductUserIDs
			if testCase.Data.(map[string]interface{})["input"] != nil {
				input = testCase.Data.(map[string]interface{})["input"].(*models.ProductUserIDs)
			}
			privileges := testCase.Data.(map[string]interface{})["privileges"].(models.Privileges)

			mysqldb.Functions = DBFunctionInterfaceMock{
				privileges: privileges,
			}

			err := validateUsers(input)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}
		})
	}
}
