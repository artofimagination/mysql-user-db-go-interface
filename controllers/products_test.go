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
	mysqldb.DBConnector = DBConnectorMock{}

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
	data := test.Data{
		Data:     make(map[string]interface{}),
		Expected: make(map[string]interface{}),
	}

	data.Data.(map[string]interface{})["input"] = product
	data.Data.(map[string]interface{})["db_mock"] = nil
	productUsers[userID] = 0
	data.Data.(map[string]interface{})["product_users"] = productUsers
	data.Data.(map[string]interface{})["privileges"] = privileges
	data.Expected.(map[string]interface{})["data"] = &product
	data.Expected.(map[string]interface{})["error"] = nil
	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "existing_product"
	data = test.Data{
		Data:     make(map[string]interface{}),
		Expected: make(map[string]interface{}),
	}

	data.Data.(map[string]interface{})["input"] = product
	data.Data.(map[string]interface{})["db_mock"] = &product
	productUsers[userID] = 0
	data.Data.(map[string]interface{})["product_users"] = productUsers
	data.Data.(map[string]interface{})["privileges"] = privileges
	data.Expected.(map[string]interface{})["data"] = nil
	data.Expected.(map[string]interface{})["error"] = fmt.Errorf(ErrProductExistsString, product.Name)
	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "incorrect_product_users"
	data = test.Data{
		Data:     make(map[string]interface{}),
		Expected: make(map[string]interface{}),
	}

	data.Data.(map[string]interface{})["input"] = product
	data.Data.(map[string]interface{})["db_mock"] = &product
	productUsers = make(models.ProductUsers)
	data.Data.(map[string]interface{})["product_users"] = productUsers
	data.Data.(map[string]interface{})["privileges"] = privileges
	data.Expected.(map[string]interface{})["data"] = nil
	data.Expected.(map[string]interface{})["error"] = ErrEmptyUsersList
	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	mysqldb.Functions = DBFunctionInterfaceMock{}
	mysqldb.DBConnector = DBConnectorMock{}
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

			mysqldb.Functions = DBFunctionInterfaceMock{
				product:    DBMock,
				privileges: privileges,
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

func TestValidateUsers_ValidData(t *testing.T) {
	// Create test data
	users, privileges := createTestUsersData()
	mysqldb.Functions = DBFunctionInterfaceMock{
		privileges: privileges,
	}
	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 0

	// Execute test
	if err := validateUsers(users); err != nil {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, nil)
		return
	}
}

func TestValidateUsers_EmptyUsersList(t *testing.T) {
	// Create test data
	users, privileges := createTestUsersData()
	mysqldb.Functions = DBFunctionInterfaceMock{
		privileges: privileges,
	}

	// Execute test
	if err := validateUsers(users); err == nil || (err != nil && err != ErrEmptyUsersList) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrEmptyUsersList)
		return
	}
}

func TestValidateUsers_NilUserList(t *testing.T) {
	// Create test data
	var users models.ProductUsers

	// Execute test
	if err := validateUsers(users); err == nil || (err != nil && err != ErrEmptyUsersList) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrEmptyUsersList)
		return
	}
}

func TestValidateUsers_NoOwner(t *testing.T) {
	// Create test data
	users, privileges := createTestUsersData()
	mysqldb.Functions = DBFunctionInterfaceMock{
		privileges: privileges,
	}
	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 1

	userID, err = uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 1

	// Execute test
	if err := validateUsers(users); err == nil || (err != nil && err != ErrInvalidOwnerCount) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrInvalidOwnerCount)
		return
	}
}

func TestValidateUsers_MultipleOwners(t *testing.T) {
	// Create test data
	users, privileges := createTestUsersData()
	mysqldb.Functions = DBFunctionInterfaceMock{
		privileges: privileges,
	}
	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 0

	userID, err = uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 0

	// Execute test
	if err := validateUsers(users); err == nil || (err != nil && err != ErrInvalidOwnerCount) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrInvalidOwnerCount)
		return
	}
}

func TestValidateUsers_InvalidPrivilege(t *testing.T) {
	// Create test data
	users, privileges := createTestUsersData()
	mysqldb.Functions = DBFunctionInterfaceMock{
		privileges: privileges,
	}
	userIDInvalid, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userIDInvalid] = 2

	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 0

	// Execute test
	Expected := fmt.Errorf(ErrUnknownPrivilegeString, 2, userIDInvalid.String())
	if err := validateUsers(users); err == nil || (err != nil && !cmp.Equal(err.Error(), Expected.Error())) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, Expected)
		return
	}
}
