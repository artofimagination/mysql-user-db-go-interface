package dbcontrollers

import (
	"fmt"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

func createTestProductsUsersData() (*models.ProductUserIDs, models.Privileges) {
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
	mysqldb.DBConnector = &DBConnectorMock{}

	users := &models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}

	return users, privileges
}

func createProductTestData() (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	dbController = &MYSQLController{}

	_, privileges := createTestProductsUsersData()

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
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
	assets := &models.Asset{
		ID:      assetID,
		DataMap: dataMap,
	}

	productData := models.ProductData{
		ID:      product.ID,
		Name:    product.Name,
		Public:  product.Public,
		Details: assets,
		Assets:  assets,
	}

	models.Interface = &ModelInterfaceMock{
		assetID:   assetID,
		productID: productID,
		asset:     assets,
	}

	testCase := "no_existing_product"
	data := make(map[string]interface{})
	data["input"] = productData
	data["db_mock"] = nil
	data["user_id"] = userID
	data["privileges"] = privileges
	expected := make(map[string]interface{})
	expected["data"] = &productData
	expected["error"] = nil
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "existing_product"
	data = make(map[string]interface{})
	data["input"] = productData
	data["db_mock"] = &product
	data["user_id"] = userID
	data["privileges"] = privileges
	expected = make(map[string]interface{})
	expected["data"] = nil
	expected["error"] = fmt.Errorf(ErrProductExistsString, product.Name)
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	mysqldb.Functions = &DBFunctionInterfaceMock{}
	mysqldb.DBConnector = &DBConnectorMock{}

	return dataSet, nil
}

func createValidationTestData() (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	productUsers, privileges := createTestProductsUsersData()

	testCase := "valid_data"
	productUsers.UserMap[userID] = 0
	data := make(map[string]interface{})
	data["input"] = productUsers
	data["privileges"] = privileges
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: nil,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "empty_user_list"
	productUsers = &models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	data = make(map[string]interface{})
	data["input"] = productUsers
	data["privileges"] = privileges
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: ErrEmptyUsersList,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "nil_user_list"
	data = make(map[string]interface{})
	data["input"] = nil
	data["privileges"] = privileges

	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: ErrEmptyUsersList,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "no_owner"
	productUsers = &models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	productUsers.UserMap[userID] = 1
	data = make(map[string]interface{})
	data["input"] = productUsers
	data["privileges"] = privileges

	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: ErrInvalidOwnerCount,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "multiple_owners"
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
	data = make(map[string]interface{})
	data["input"] = productUsers
	data["privileges"] = privileges
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: ErrInvalidOwnerCount,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "invalid_privilege"
	productUsers = &models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	productUsers.UserMap[userID] = 2
	data = make(map[string]interface{})
	data["input"] = productUsers
	data["privileges"] = privileges
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: fmt.Errorf(ErrUnknownPrivilegeString, 2, userID.String()),
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	return dataSet, nil
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
			mysqldb.Functions = &DBFunctionInterfaceMock{
				product:      mockCopy,
				privileges:   privileges,
				productAdded: false,
			}

			output, err := dbController.CreateProduct(
				input.Name,
				input.Public,
				&userID,
				func(*uuid.UUID) (string, error) {
					return "testPath", nil
				})

			if diff := pretty.Diff(expectedData, output); len(diff) != 0 {
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

			mysqldb.Functions = &DBFunctionInterfaceMock{
				privileges: privileges,
			}

			err := validateOwnership(input)
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
				return
			}
		})
	}
}
