package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	AddProductTest                    = 0
	AddProductUsersTest               = 1
	DeleteProductUsersByProductIDTest = 2
	GetProductByIDTest                = 3
	GetProductByNameTest              = 4
	GetUserProductIDsTest             = 5
	GetProductsByUserIDTest           = 6
	DeleteProductTest                 = 7
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

	product := models.Product{
		ID:        productID,
		Name:      "Test",
		Public:    true,
		DetailsID: detailsID,
		AssetsID:  assetID,
	}

	return &product, nil
}

func addProductsToMock(products []models.Product) (*sqlmock.Rows, error) {
	rowsProducts := sqlmock.NewRows([]string{"id", "name", "public", "details", "product_assets_id"})
	for _, product := range products {
		binaryProductID, err := json.Marshal(product.ID)
		if err != nil {
			return nil, err
		}

		binaryDetailsID, err := json.Marshal(product.DetailsID)
		if err != nil {
			return nil, err
		}

		binaryAssetID, err := json.Marshal(product.AssetsID)
		if err != nil {
			return nil, err
		}

		rowsProducts.AddRow(binaryProductID, product.Name, product.Public, binaryDetailsID, binaryAssetID)
	}
	return rowsProducts, nil
}

func createTestProductList(quantity int) ([]models.Product, error) {
	// Create test data
	products := []models.Product{}
	for ; quantity > 0; quantity-- {
		product, err := createTestProductData()
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}
	return products, nil
}

func createTestUserProductsData(quantity int) (models.UserProducts, error) {
	userProducts := make(models.UserProducts)
	for ; quantity > 0; quantity-- {
		productID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		userProducts[productID] = 1
	}
	return userProducts, nil
}

func createTestProductUsersData() (models.ProductUsers, error) {
	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	owners := make(models.ProductUsers)
	owners[userID] = 1
	return owners, nil
}

func createProductsTestData(testID int) (*test.OrderedTests, DBConnectorMock, error) {
	dbConnector := DBConnectorMock{}
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, dbConnector, err
	}

	data := test.Data{}

	product, err := createTestProductData()
	if err != nil {
		return nil, dbConnector, err
	}

	binaryProductID, err := json.Marshal(product.ID)
	if err != nil {
		return nil, dbConnector, err
	}

	binaryAssetID, err := json.Marshal(product.AssetsID)
	if err != nil {
		return nil, dbConnector, err
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, dbConnector, err
	}

	userProducts, err := createTestUserProductsData(2)
	if err != nil {
		return nil, dbConnector, err
	}

	switch testID {
	case AddProductTest:

		testCase := "valid_product"
		data = test.Data{
			Data:     product,
			Expected: nil,
		}
		mock.ExpectBegin()
		mock.ExpectExec(AddProductQuery).WithArgs(product.ID, product.Name, product.Public, product.DetailsID, product.AssetsID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		data = test.Data{
			Data:     product,
			Expected: errors.New("This is a failure test"),
		}
		mock.ExpectBegin()
		mock.ExpectExec(AddProductQuery).WithArgs(product.ID, product.Name, product.Public, product.DetailsID, product.AssetsID).WillReturnError(data.Expected.(error))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "duplicate_name"
		data = test.Data{
			Data:     product,
			Expected: fmt.Errorf(ErrSQLDuplicateProductNameEntryString, product.Name),
		}
		mock.ExpectBegin()
		mock.ExpectExec(AddProductQuery).WithArgs(product.ID, product.Name, product.Public, product.DetailsID, product.AssetsID).WillReturnError(data.Expected.(error))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case AddProductUsersTest:

		testCase := "valid_products"
		productUsers, err := createTestProductUsersData()
		if err != nil {
			return nil, dbConnector, err
		}

		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: nil,
		}
		data.Data.(map[string]interface{})["product_id"] = product.ID
		data.Data.(map[string]interface{})["product_users"] = productUsers

		mock.ExpectBegin()
		for userID, privilege := range productUsers {
			mock.ExpectExec(AddProductUsersQuery).WithArgs(product.ID, userID, privilege).WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: errors.New("This is a failure test"),
		}
		data.Data.(map[string]interface{})["product_id"] = product.ID
		data.Data.(map[string]interface{})["product_users"] = productUsers
		mock.ExpectBegin()
		for userID, privilege := range productUsers {
			mock.ExpectExec(AddProductUsersQuery).WithArgs(product.ID, userID, privilege).WillReturnError(data.Expected.(error))
		}
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case DeleteProductUsersByProductIDTest:

		testCase := "valid_id"
		data = test.Data{
			Data:     product.ID,
			Expected: nil,
		}
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductUsersByProductIDQuery).WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_id"
		data = test.Data{
			Data:     product.ID,
			Expected: ErrNoUserWithProduct,
		}
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductUsersByProductIDQuery).WithArgs(product.ID).WillReturnError(ErrNoUserWithProduct)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case GetProductByIDTest:

		testCase := "valid_id"
		data = test.Data{
			Data:     product.ID,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = product
		data.Expected.(map[string]interface{})["error"] = nil
		rows := sqlmock.NewRows([]string{"id", "name", "public", "details", "product_assets_id"}).
			AddRow(binaryProductID, product.Name, product.Public, product.DetailsID, product.AssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByIDQuery).WithArgs(product.ID).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_id"
		data = test.Data{
			Data:     product.ID,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = sql.ErrNoRows
		data.Data = product.ID
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByIDQuery).WithArgs(product.ID).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case GetProductByNameTest:

		testCase := "valid_name"
		data = test.Data{
			Data:     product.Name,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = product
		data.Expected.(map[string]interface{})["error"] = nil
		rows := sqlmock.NewRows([]string{"id", "name", "public", "details", "product_assets_id"}).
			AddRow(binaryProductID, product.Name, product.Public, product.DetailsID, binaryAssetID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByNameQuery).WithArgs(product.Name).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_name"
		data = test.Data{
			Data:     product.Name,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByNameQuery).WithArgs(product.Name).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case GetUserProductIDsTest:

		testCase := "valid_id"
		data = test.Data{
			Data:     userID,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = userProducts
		data.Expected.(map[string]interface{})["error"] = nil
		rows := sqlmock.NewRows([]string{"products_id", "privilege"})
		for productID, privilege := range userProducts {
			rows.AddRow(productID, privilege)
		}

		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_products"
		data = test.Data{
			Data:     userID,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = sql.ErrNoRows

		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case GetProductsByUserIDTest:

		products, err := createTestProductList(2)
		if err != nil {
			return nil, dbConnector, err
		}

		testCase := "valid_id"
		data = test.Data{
			Data:     userID,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = products
		data.Expected.(map[string]interface{})["error"] = nil
		rowsUserProducts := sqlmock.NewRows([]string{"products_id", "privilege"})
		// Ranging through map is random, need to collect product ID-s in fixed order in order to
		// have the correct order of sql mock expectations.
		orderedProductIDs := make([]uuid.UUID, 0)
		for productID, privilege := range userProducts {

			rowsUserProducts.AddRow(productID, privilege)
			orderedProductIDs = append(orderedProductIDs, productID)
		}

		products[0].ID = orderedProductIDs[0]
		products[1].ID = orderedProductIDs[1]

		rowsProducts, err := addProductsToMock(products)
		if err != nil {
			return nil, dbConnector, err
		}

		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnRows(rowsUserProducts)
		for _, productID := range orderedProductIDs {
			mock.ExpectQuery(GetProductByIDQuery).WithArgs(productID).WillReturnRows(rowsProducts)
		}
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "no_products"
		data = test.Data{
			Data:     userID,
			Expected: make(map[string]interface{}),
		}
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case DeleteProductTest:

		testCase := "valid_id"
		data = test.Data{
			Data:     product.ID,
			Expected: nil,
		}
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM users_products where products_id = UUID_TO_BIN(?)").WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("DELETE FROM products where id = UUID_TO_BIN(?)").WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	default:
		return nil, dbConnector, fmt.Errorf("Unknown test %d", testID)
	}

	dbConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}
	Functions = MYSQLFunctions{}

	return &dataSet, dbConnector, nil
}

func createPrivilegesTestData() (*test.OrderedTests, DBConnectorMock, error) {
	dbConnector := DBConnectorMock{}
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, dbConnector, err
	}

	privileges := make(models.Privileges, 2)
	privileges[0].ID = 0
	privileges[0].Name = "test0"
	privileges[0].Description = "description0"
	privileges[1].ID = 1
	privileges[1].Name = "test1"
	privileges[1].Description = "description1"

	testCase := "valid_id"
	data := test.Data{
		Data:     nil,
		Expected: make(map[string]interface{}),
	}
	data.Expected.(map[string]interface{})["data"] = privileges
	data.Expected.(map[string]interface{})["error"] = nil
	rows := sqlmock.NewRows([]string{"id", "name", "description"})
	for _, privilege := range privileges {
		rows.AddRow(privilege.ID, privilege.Name, privilege.Description)
	}

	mock.ExpectBegin()
	mock.ExpectQuery(GetPrivilegesQuery).WillReturnRows(rows)
	mock.ExpectCommit()
	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "invalid_id"
	data = test.Data{
		Data:     nil,
		Expected: make(map[string]interface{}),
	}
	data.Expected.(map[string]interface{})["data"] = nil
	data.Expected.(map[string]interface{})["error"] = sql.ErrNoRows

	mock.ExpectBegin()
	mock.ExpectQuery(GetPrivilegesQuery).WillReturnError(sql.ErrNoRows)
	mock.ExpectCommit()
	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	dbConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}
	Functions = MYSQLFunctions{}

	return &dataSet, dbConnector, nil
}

func TestAddProduct(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createProductsTestData(AddProductTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := dbConnector.DB.Begin()
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

			err = Functions.AddProduct(product, tx)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, testCase.Expected)
				return
			}
		})
	}
}

func TestAddProductUsers(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createProductsTestData(AddProductUsersTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := dbConnector.DB.Begin()
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

			err = Functions.AddProductUsers(&productID, testData["product_users"].(models.ProductUsers), tx)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, testCase.Expected)
				return
			}
		})
	}
}

func TestDeleteProductUsersByProductID(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createProductsTestData(DeleteProductUsersByProductIDTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := dbConnector.DB.Begin()
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

			err = deleteProductUsersByProductID(&productID, tx)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, testCase.Expected)
				return
			}
		})
	}
}

func TestGetProductByID(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createProductsTestData(GetProductByIDTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	DBConnector = dbConnector
	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
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

			output, err := Functions.GetProductByID(productID)
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

func TestGetProductByName(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createProductsTestData(GetProductByNameTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := dbConnector.DB.Begin()
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

			output, err := Functions.GetProductByName(productName, tx)
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

func TestGetUserProductIDs(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createProductsTestData(GetUserProductIDsTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := dbConnector.DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			userID := testCase.Data.(uuid.UUID)
			var expectedData models.UserProducts
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(models.UserProducts)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}

			output, err := Functions.GetUserProductIDs(userID, tx)
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

func TestGetProductsByUserID(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createProductsTestData(GetProductsByUserIDTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	DBConnector = dbConnector
	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			userID := testCase.Data.(uuid.UUID)
			var expectedData []models.Product
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].([]models.Product)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}

			output, err := Functions.GetProductsByUserID(userID)
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

func TestDeleteProduct(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createProductsTestData(DeleteProductTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	DBConnector = dbConnector
	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		testCase := dataSet.TestDataSet[testCaseString]
		var expectedError error
		if testCase.Expected != nil {
			expectedError = testCase.Expected.(error)
		}
		data := testCase.Data.(uuid.UUID)

		err = Functions.DeleteProduct(&data)
		if !test.ErrEqual(err, expectedError) {
			t.Errorf(test.TestResultString, testCaseString, err, expectedError)
			return
		}
	}
}

func TestGetPrivileges(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createPrivilegesTestData()
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	DBConnector = dbConnector
	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		testCase := dataSet.TestDataSet[testCaseString]
		var expectedData models.Privileges
		if testCase.Expected.(map[string]interface{})["data"] != nil {
			expectedData = testCase.Expected.(map[string]interface{})["data"].(models.Privileges)
		}
		var expectedError error
		if testCase.Expected.(map[string]interface{})["error"] != nil {
			expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
		}

		output, err := Functions.GetPrivileges()
		if !cmp.Equal(output, expectedData) {
			t.Errorf(test.TestResultString, testCaseString, output, expectedData)
			return
		}

		if !test.ErrEqual(err, expectedError) {
			t.Errorf(test.TestResultString, testCaseString, err, expectedError)
			return
		}
	}
}
