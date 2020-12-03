package controllers

import (
	"fmt"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func createTestUsersData() (models.ProductUsers, models.Privileges) {
	privileges := make(models.Privileges, 2)
	privileges[0].ID = 0
	privileges[0].Name = "Owner"
	privileges[0].Description = "description0"
	privileges[1].ID = 1
	privileges[1].Name = "User"
	privileges[1].Description = "description1"
	mysqldb.DBConnector = &DBConnectorMock{}

	users := make(models.ProductUsers)

	return users, privileges
}

func createProductTestData() (*test.OrderedTests, error) {
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	productUsers, privileges := createTestUsersData()

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

	testCase := "no_existing_product"
	data := make(map[string]interface{})
	data["input"] = product
	data["db_mock"] = nil
	productUsers[userID] = 0
	data["product_users"] = productUsers
	data["privileges"] = privileges
	expected := make(map[string]interface{})
	expected["data"] = &product
	expected["error"] = nil
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "existing_product"
	data = make(map[string]interface{})
	data["input"] = product
	data["db_mock"] = &product
	productUsers[userID] = 0
	data["product_users"] = productUsers
	data["privileges"] = privileges
	expected = make(map[string]interface{})
	expected["data"] = nil
	expected["error"] = fmt.Errorf(ErrProductExistsString, product.Name)
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "incorrect_product_users"
	data = make(map[string]interface{})
	data["input"] = product
	data["db_mock"] = &product
	productUsers = make(models.ProductUsers)
	data["product_users"] = productUsers
	data["privileges"] = privileges
	expected = make(map[string]interface{})
	expected["data"] = nil
	expected["error"] = ErrEmptyUsersList
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	mysqldb.Functions = DBFunctionInterfaceMock{}
	mysqldb.DBConnector = &DBConnectorMock{}
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
	productUsers[userID] = 0
	data := make(map[string]interface{})
	data["input"] = productUsers
	data["privileges"] = privileges

	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: nil,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "empty_user_list"
	data = make(map[string]interface{})
	productUsers = make(models.ProductUsers)
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
	productUsers = make(models.ProductUsers)
	productUsers[userID] = 1
	data = make(map[string]interface{})
	data["input"] = productUsers
	data["privileges"] = privileges

	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: ErrInvalidOwnerCount,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "multiple_owners"
	productUsers = make(models.ProductUsers)
	productUsers[userID] = 0
	userID2, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	productUsers[userID2] = 0
	data = make(map[string]interface{})
	data["input"] = productUsers
	data["privileges"] = privileges

	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: ErrInvalidOwnerCount,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "invalid_privilege"
	productUsers = make(models.ProductUsers)
	productUsers[userID] = 2
	data = make(map[string]interface{})
	data["input"] = productUsers
	data["privileges"] = privileges

	dataSet.TestDataSet[testCase] = test.Data{
		Data:     data,
		Expected: fmt.Errorf(ErrUnknownPrivilegeString, 2, userID.String()),
	}
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
			var expectedData *models.Product
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.Product)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}
			input := testCase.Data.(map[string]interface{})["input"].(models.Product)
			productUsers := testCase.Data.(map[string]interface{})["product_users"].(models.ProductUsers)
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

			output, err := CreateProduct(
				input.Name,
				input.Public,
				productUsers,
				func(*uuid.UUID) string {
					return "testPath"
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
			var input models.ProductUsers
			if testCase.Data.(map[string]interface{})["input"] != nil {
				input = testCase.Data.(map[string]interface{})["input"].(models.ProductUsers)
			}
			privileges := testCase.Data.(map[string]interface{})["privileges"].(models.Privileges)

			mysqldb.Functions = DBFunctionInterfaceMock{
				privileges: privileges,
			}

			err := validateOwnership(input)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}
		})
	}
}
