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
	addProductTest = iota
	addProductUsersTest
	deleteProductUsersByProductIDTest
	getProductByIDTest
	getProductByNameTest
	GetUserProductIDsTest
	deleteProductTest
	updateUsersProductsTest
)

func createTestProductData() (*models.Product, error) {
	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	detailsID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		ID:        productID,
		Name:      "Test",
		Public:    true,
		DetailsID: detailsID,
		AssetsID:  assetID,
	}

	return product, nil
}

func createTestUserProductsData(quantity int) (*models.UserProductIDs, error) {
	userProducts := &models.UserProductIDs{
		ProductMap:     make(map[uuid.UUID]int),
		ProductIDArray: make([]uuid.UUID, 0),
	}

	for ; quantity > 0; quantity-- {
		productID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		userProducts.ProductMap[productID] = 1
		userProducts.ProductIDArray = append(userProducts.ProductIDArray, productID)
	}
	return userProducts, nil
}

func createProductsTestData(testID int) (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, err
	}

	product, err := createTestProductData()
	if err != nil {
		return nil, err
	}

	binaryProductID, err := json.Marshal(product.ID)
	if err != nil {
		return nil, err
	}

	binaryAssetID, err := json.Marshal(product.AssetsID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	userProducts, err := createTestUserProductsData(2)
	if err != nil {
		return nil, err
	}

	productUsers, err := createTestProductUsersData()
	if err != nil {
		return nil, err
	}

	switch testID {
	case addProductTest:

		testCase := "valid_product"
		mock.ExpectBegin()
		mock.ExpectExec(AddProductQuery).WithArgs(product.ID, product.Name, product.Public, product.DetailsID, product.AssetsID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     product,
			Expected: nil,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		expected := errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectExec(AddProductQuery).WithArgs(product.ID, product.Name, product.Public, product.DetailsID, product.AssetsID).WillReturnError(expected)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     product,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "duplicate_name"
		expected = fmt.Errorf(ErrSQLDuplicateProductNameEntryString, product.Name)
		mock.ExpectBegin()
		mock.ExpectExec(AddProductQuery).WithArgs(product.ID, product.Name, product.Public, product.DetailsID, product.AssetsID).WillReturnError(expected)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     product,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case addProductUsersTest:
		testCase := "valid_products"
		data := make(map[string]interface{})
		data["product_id"] = product.ID
		data["product_users"] = productUsers
		mock.ExpectBegin()
		for _, userID := range productUsers.UserIDArray {
			privilege := productUsers.UserMap[userID]
			mock.ExpectExec(AddProductUsersQuery).WithArgs(userID, product.ID, privilege).WillReturnResult(sqlmock.NewResult(1, 1))
		}
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: nil,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		data = make(map[string]interface{})
		expected := errors.New("This is a failure test")
		data["product_id"] = product.ID
		data["product_users"] = productUsers
		mock.ExpectBegin()
		for _, userID := range productUsers.UserIDArray {
			privilege := productUsers.UserMap[userID]
			mock.ExpectExec(AddProductUsersQuery).WithArgs(userID, product.ID, privilege).WillReturnError(expected)
		}
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_to_add"
		data = make(map[string]interface{})
		expected = ErrNoProductUserAdded
		data["product_id"] = product.ID
		data["product_users"] = productUsers
		mock.ExpectBegin()
		for _, userID := range productUsers.UserIDArray {
			privilege := productUsers.UserMap[userID]
			mock.ExpectExec(AddProductUsersQuery).WithArgs(userID, product.ID, privilege).WillReturnResult(sqlmock.NewResult(1, 0))
		}
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case deleteProductUsersByProductIDTest:
		testCase := "valid_id"
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductUsersByProductIDQuery).WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     product.ID,
			Expected: nil,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_id"
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductUsersByProductIDQuery).WithArgs(product.ID).WillReturnError(ErrNoUserWithProduct)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     product.ID,
			Expected: ErrNoUserWithProduct,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case getProductByIDTest:

		testCase := "valid_id"
		expected := make(map[string]interface{})
		expected["data"] = product
		expected["error"] = nil
		rows := sqlmock.NewRows([]string{"id", "name", "public", "details", "product_assets_id"}).
			AddRow(binaryProductID, product.Name, product.Public, product.DetailsID, product.AssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByIDQuery).WithArgs(product.ID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     product.ID,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_id"
		expected = make(map[string]interface{})
		expected["data"] = nil
		expected["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByIDQuery).WithArgs(product.ID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     product.ID,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case getProductByNameTest:
		testCase := "valid_name"
		expected := make(map[string]interface{})
		expected["data"] = product
		expected["error"] = nil
		rows := sqlmock.NewRows([]string{"id", "name", "public", "details", "product_assets_id"}).
			AddRow(binaryProductID, product.Name, product.Public, product.DetailsID, binaryAssetID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByNameQuery).WithArgs(product.Name).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     product.Name,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_name"
		expected = make(map[string]interface{})
		expected["data"] = nil
		expected["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByNameQuery).WithArgs(product.Name).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     product.Name,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case GetUserProductIDsTest:
		testCase := "valid_id"
		expected := make(map[string]interface{})
		expected["data"] = userProducts
		expected["error"] = nil
		rows := sqlmock.NewRows([]string{"products_id", "privilege"})
		for _, productID := range userProducts.ProductIDArray {
			rows.AddRow(productID, userProducts.ProductMap[productID])
		}
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     userID,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_products"
		expected = make(map[string]interface{})
		expected["data"] = nil
		expected["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     userID,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case deleteProductTest:
		testCase := "valid_id"
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductQuery).WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     product.ID,
			Expected: nil,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "no_product"
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductQuery).WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     product.ID,
			Expected: ErrNoProductDeleted,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case updateUsersProductsTest:
		testCase := "valid_id"
		data := make(map[string]interface{})
		data["user_id"] = userID
		data["product_id"] = product.ID
		data["privilege"] = 1
		mock.ExpectBegin()
		mock.ExpectExec(UpdateUsersProductsQuery).WithArgs(data["privilege"].(int), userID, product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: nil,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "no_users_products"
		data = make(map[string]interface{})
		data["user_id"] = userID
		data["product_id"] = product.ID
		data["privilege"] = 1
		mock.ExpectBegin()
		mock.ExpectExec(UpdateUsersProductsQuery).WithArgs(data["privilege"].(int), userID, product.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: ErrNoUsersProductUpdate,
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

func TestAddProduct(t *testing.T) {
	// Create test data
	dataSet, err := createProductsTestData(addProductTest)
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
			var expectedError error
			if testCase.Expected != nil {
				expectedError = testCase.Expected.(error)
			}
			product := testCase.Data.(*models.Product)

			err = DBFunctions.AddProduct(product, tx)
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
				return
			}
		})
	}
}

func TestAddProductUsers(t *testing.T) {
	// Create test data
	dataSet, err := createProductsTestData(addProductUsersTest)
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
			var expectedError error
			if testCase.Expected != nil {
				expectedError = testCase.Expected.(error)
			}
			testData := testCase.Data.(map[string]interface{})
			productID := testData["product_id"].(uuid.UUID)

			productUsers := testData["product_users"].(*models.ProductUserIDs)
			err = DBFunctions.AddProductUsers(&productID, productUsers, tx)
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
				return
			}
		})
	}
}

func TestUpdateUsersProducts(t *testing.T) {
	// Create test data
	dataSet, err := createProductsTestData(updateUsersProductsTest)
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
			var expectedError error
			if testCase.Expected != nil {
				expectedError = testCase.Expected.(error)
			}
			userID := testCase.Data.(map[string]interface{})["user_id"].(uuid.UUID)
			productID := testCase.Data.(map[string]interface{})["product_id"].(uuid.UUID)
			privilege := testCase.Data.(map[string]interface{})["privilege"].(int)

			err = DBFunctions.UpdateUsersProducts(&userID, &productID, privilege, tx)
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
				return
			}
		})
	}
}

func TestDeleteProductUsersByProductID(t *testing.T) {
	// Create test data
	dataSet, err := createProductsTestData(deleteProductUsersByProductIDTest)
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
			var expectedError error
			if testCase.Expected != nil {
				expectedError = testCase.Expected.(error)
			}
			productID := testCase.Data.(uuid.UUID)

			err = DBFunctions.DeleteProductUsersByProductID(&productID, tx)
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
				return
			}
		})
	}
}

func TestGetProductByID(t *testing.T) {
	// Create test data
	dataSet, err := createProductsTestData(getProductByIDTest)
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
			productID := testCase.Data.(uuid.UUID)
			var expectedData *models.Product
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.Product)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}

			output, err := DBFunctions.GetProductByID(&productID, tx)
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

func TestGetProductByName(t *testing.T) {
	// Create test data
	dataSet, err := createProductsTestData(getProductByNameTest)
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
			productName := testCase.Data.(string)
			var expectedData *models.Product
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.Product)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}

			output, err := DBFunctions.GetProductByName(productName, tx)
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

func TestGetUserProductIDs(t *testing.T) {
	// Create test data
	dataSet, err := createProductsTestData(GetUserProductIDsTest)
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
			userID := testCase.Data.(uuid.UUID)
			var expectedData *models.UserProductIDs
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.UserProductIDs)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}

			output, err := DBFunctions.GetUserProductIDs(&userID, tx)
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

func TestDeleteProduct(t *testing.T) {
	// Create test data
	dataSet, err := createProductsTestData(deleteProductTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
		if err != nil {
			t.Errorf("Failed to setup DB transaction %s", err)
			return
		}

		testCaseString := testCaseString
		testCase := dataSet.TestDataSet[testCaseString]
		var expectedError error
		if testCase.Expected != nil {
			expectedError = testCase.Expected.(error)
		}
		data := testCase.Data.(uuid.UUID)

		err = DBFunctions.DeleteProduct(&data, tx)
		if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
			t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
			return
		}
	}
}
