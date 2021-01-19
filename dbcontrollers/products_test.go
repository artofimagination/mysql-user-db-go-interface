package dbcontrollers

import (
	"fmt"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

type ProductExpectedData struct {
	productData *models.ProductData
	err         error
}

type ProductMockData struct {
	productData *models.ProductData
	product     *models.Product
	privileges  models.Privileges
}

type ProductInputData struct {
	productData *models.ProductData
	userID      uuid.UUID
}

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

	_, privileges := createTestProductsUsersData()

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		Name:      "testProduct",
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

	productData := &models.ProductData{
		ID:      product.ID,
		Name:    product.Name,
		Details: assets,
		Assets:  assets,
	}

	dbController = &MYSQLController{
		DBFunctions: &DBFunctionMock{},
		DBConnector: &DBConnectorMock{},
		ModelFunctions: &ModelMock{
			assetID:   assetID,
			productID: productID,
			asset:     assets,
		},
	}

	testCase := "no_existing_product"
	expected := ProductExpectedData{
		productData: productData,
		err:         nil,
	}
	input := ProductInputData{
		productData: productData,
		userID:      userID,
	}
	mock := ProductMockData{
		productData: nil,
		privileges:  privileges,
	}
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     input,
		Mock:     mock,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "existing_product"
	expected = ProductExpectedData{
		productData: nil,
		err:         fmt.Errorf(ErrProductExistsString, product.Name),
	}
	input = ProductInputData{
		productData: productData,
		userID:      userID,
	}
	mock = ProductMockData{
		product:    product,
		privileges: privileges,
	}
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     input,
		Mock:     mock,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	return dataSet, nil
}

type ValidationExpectedData struct {
	err error
}

type ValidationMockData struct {
	privileges models.Privileges
}

type ValidationInputData struct {
	productUsers *models.ProductUserIDs
}

func createValidationTestData() (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	dbController = &MYSQLController{
		DBFunctions: &DBFunctionMock{},
		DBConnector: &DBConnectorMock{},
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	productUsers, privileges := createTestProductsUsersData()

	testCase := "valid_data"
	productUsers.UserMap[userID] = 0
	expected := ValidationExpectedData{
		err: nil,
	}
	input := ValidationInputData{
		productUsers: productUsers,
	}
	mock := ValidationMockData{
		privileges: privileges,
	}
	productUsers.UserMap[userID] = 0
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     input,
		Mock:     mock,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "empty_user_list"
	productUsers.UserMap[userID] = 0
	productUsers = &models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	expected = ValidationExpectedData{
		err: ErrEmptyUsersList,
	}
	input = ValidationInputData{
		productUsers: productUsers,
	}
	mock = ValidationMockData{
		privileges: privileges,
	}
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     input,
		Mock:     mock,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "nil_user_list"
	expected = ValidationExpectedData{
		err: ErrEmptyUsersList,
	}
	input = ValidationInputData{
		productUsers: nil,
	}
	mock = ValidationMockData{
		privileges: privileges,
	}
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     input,
		Mock:     mock,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "no_owner"
	productUsers = &models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	productUsers.UserMap[userID] = 1
	expected = ValidationExpectedData{
		err: ErrInvalidOwnerCount,
	}
	input = ValidationInputData{
		productUsers: productUsers,
	}
	mock = ValidationMockData{
		privileges: privileges,
	}
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     input,
		Mock:     mock,
		Expected: expected,
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
	expected = ValidationExpectedData{
		err: ErrInvalidOwnerCount,
	}
	input = ValidationInputData{
		productUsers: productUsers,
	}
	mock = ValidationMockData{
		privileges: privileges,
	}
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     input,
		Mock:     mock,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "invalid_privilege"
	productUsers = &models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	productUsers.UserMap[userID] = 2
	expected = ValidationExpectedData{
		err: fmt.Errorf(ErrUnknownPrivilegeString, 2, userID.String()),
	}
	input = ValidationInputData{
		productUsers: productUsers,
	}
	mock = ValidationMockData{
		privileges: privileges,
	}
	dataSet.TestDataSet[testCase] = test.Data{
		Data:     input,
		Mock:     mock,
		Expected: expected,
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
			expectedData := testCase.Expected.(ProductExpectedData)
			inputData := testCase.Data.(ProductInputData)
			mockData := testCase.Mock.(ProductMockData)

			dbController.DBFunctions = &DBFunctionMock{
				product:      mockData.product,
				privileges:   mockData.privileges,
				projectAdded: false,
			}

			output, err := dbController.CreateProduct(
				inputData.productData.Name,
				&inputData.userID,
				func(*uuid.UUID) (string, error) {
					return "testPath", nil
				})

			if diff := pretty.Diff(expectedData.productData, output); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData.productData, diff)
				return
			}

			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
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
			expectedData := testCase.Expected.(ValidationExpectedData)
			inputData := testCase.Data.(ValidationInputData)
			mockData := testCase.Mock.(ValidationMockData)

			dbController.DBFunctions = &DBFunctionMock{
				privileges: mockData.privileges,
			}

			err := dbController.validateOwnership(inputData.productUsers)
			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
				return
			}
		})
	}
}
