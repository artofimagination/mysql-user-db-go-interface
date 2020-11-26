package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"models"

	"github.com/pkg/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
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
	details := make(models.Details)
	details[models.SupportClients] = true
	details[models.ProjectUI] = true
	details[models.Requires3D] = true
	details[models.ClientUI] = true

	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	product := models.Product{
		ID:       productID,
		Name:     "Test",
		Public:   true,
		Details:  details,
		AssetsID: assetID,
	}

	return &product, nil
}

func addProductsToMock(products *[]models.Product) (*sqlmock.Rows, error) {
	rowsProducts := sqlmock.NewRows([]string{"id", "name", "public", "details", "product_assets_id"})
	for _, product := range *products {
		binaryID, err := json.Marshal(product.ID)
		if err != nil {
			return nil, err
		}
		jsonRaw, err := ConvertToJSONRaw(product.Details)
		if err != nil {
			return nil, err
		}
		rowsProducts.AddRow(binaryID, product.Name, product.Public, jsonRaw, product.AssetsID)
	}
	return rowsProducts, nil
}

func createTestProductList(quantity int) (*[]models.Product, error) {
	// Create test data
	products := []models.Product{}
	for ; quantity > 0; quantity-- {
		product, err := createTestProductData()
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}
	return &products, nil
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

func createProductsTestData(test int) (*testhelpers.OrderedTests, DBConnectorMock, error) {
	dbConnector := DBConnectorMock{}
	dataSet := testhelpers.OrderedTests{
		orderedList: make(OrderedTestList, 0),
		testDataSet: make(TestDataSet, 0),
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, dbConnector, err
	}

	data := testhelpers.TestData{}

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

	testCase := "valid_product"
	jsonRaw, err := ConvertToJSONRaw(product.Details)
	if err != nil {
		return nil, dbConnector, err
	}

	userProducts, err := createTestUserProductsData(2)
	if err != nil {
		return nil, dbConnector, err
	}

	switch test {
	case AddProductTest:

		data.expected = nil
		mock.ExpectBegin()
		mock.ExpectExec(AddProductQuery).WithArgs(product.ID, product.Name, product.Public, jsonRaw, product.AssetsID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "failed_query"
		data.expected = errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectExec(GetUserByEmailQuery).WithArgs(product.ID, product.Name, product.Public, jsonRaw, product.AssetsID).WillReturnError(data.expected.(error))
		mock.ExpectRollback()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "duplicate_name"
		data.expected = fmt.Errorf(ErrSQLDuplicateProductNameEntryString, product.Name)
		mock.ExpectBegin()
		mock.ExpectExec(GetUserByEmailQuery).WithArgs(product.ID, product.Name, product.Public, jsonRaw, product.AssetsID).WillReturnError(data.expected.(error))
		mock.ExpectRollback()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

	case AddProductUsersTest:

		testCase := "valid_products"
		data.expected = nil

		productUsers, err := createTestProductUsersData()
		if err != nil {
			return nil, dbConnector, err
		}

		data.data = make(map[string]interface{})
		data.data.(map[string]interface{})["product_id"] = product.ID
		data.data.(map[string]interface{})["product_users"] = productUsers

		mock.ExpectBegin()
		for userID, privilege := range productUsers {
			mock.ExpectExec(AddProductUsersQuery).WithArgs(product.ID, userID, privilege).WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "failed_query"
		data.expected = errors.New("This is a failure test")
		mock.ExpectBegin()
		for userID, privilege := range productUsers {
			mock.ExpectExec(AddProductUsersQuery).WithArgs(product.ID, userID, privilege).WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectRollback()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

	case DeleteProductUsersByProductIDTest:

		testCase := "valid_id"
		data.expected = nil
		data.data = product.ID
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductUsersByProductIDQuery).WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "missing_id"
		data.expected = ErrNoUserWithProduct
		data.data = product.ID
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductUsersByProductIDQuery).WithArgs(product.ID).WillReturnError(ErrNoUserWithProduct)
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

	case GetProductByIDTest:

		testCase := "valid_id"
		data.expected = make(map[string]interface{})
		data.expected.(map[string]interface{})["data"] = product
		data.expected.(map[string]interface{})["error"] = nil
		data.data = product.ID
		rows := sqlmock.NewRows([]string{"id", "name", "public", "details", "product_assets_id"}).
			AddRow(binaryProductID, product.Name, product.Public, jsonRaw, product.AssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByIDQuery).WithArgs(product.ID).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "missing_id"
		data.expected = make(map[string]interface{})
		data.expected.(map[string]interface{})["data"] = nil
		data.expected.(map[string]interface{})["error"] = sql.ErrNoRows
		data.data = product.ID
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByIDQuery).WithArgs(product.ID).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

	case GetProductByNameTest:

		testCase := "valid_name"
		data.expected = make(map[string]interface{})
		data.expected.(map[string]interface{})["data"] = product
		data.expected.(map[string]interface{})["error"] = nil
		data.data = product.Name
		rows := sqlmock.NewRows([]string{"id", "name", "public", "details", "product_assets_id"}).
			AddRow(binaryProductID, product.Name, product.Public, jsonRaw, binaryAssetID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByNameQuery).WithArgs(product.ID).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "missing_name"
		data.expected = make(map[string]interface{})
		data.expected.(map[string]interface{})["data"] = nil
		data.expected.(map[string]interface{})["error"] = sql.ErrNoRows
		data.data = product.Name
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByNameQuery).WithArgs(product.ID).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

	case GetUserProductIDsTest:

		testCase := "valid_id"
		data.expected = make(map[string]interface{})
		data.expected.(map[string]interface{})["data"] = userProducts
		data.expected.(map[string]interface{})["error"] = nil
		data.data = userID
		rows := sqlmock.NewRows([]string{"products_id", "privilege"})
		for productID, privilege := range userProducts {
			rows.AddRow(productID, privilege)
		}

		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "missing_products"
		data.expected = make(map[string]interface{})
		data.expected.(map[string]interface{})["data"] = nil
		data.expected.(map[string]interface{})["error"] = sql.ErrNoRows
		data.data = userID

		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

	case GetProductsByUserIDTest:

		products, err := createTestProductList(2)
		if err != nil {
			return nil, dbConnector, err
		}

		testCase := "valid_id"
		data.data = userID
		data.expected.(map[string]interface{})["data"] = products
		data.expected.(map[string]interface{})["error"] = nil
		rowsUserProducts := sqlmock.NewRows([]string{"products_id", "privilege"})
		// Ranging through map is random, need to collect product ID-s in fixed order in order to
		// have the correct order of sql mock expectations.
		orderedProductIDs := make([]uuid.UUID, 0)
		for productID, privilege := range userProducts {

			rowsUserProducts.AddRow(productID, privilege)
			orderedProductIDs = append(orderedProductIDs, productID)
		}

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
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

		testCase = "no_products"
		data.data = userID
		data.expected.(map[string]interface{})["data"] = nil
		data.expected.(map[string]interface{})["error"] = sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

	case DeleteProductTest:

		testCase := "valid_id"
		data.data = product.ID
		data.expected = nil
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM users_products where products_id = UUID_TO_BIN(?)").WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("DELETE FROM products where id = UUID_TO_BIN(?)").WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		dataSet.testDataSet[testCase] = data
		dataSet.orderedList = append(dataSet.orderedList, testCase)

	default:
		return nil, dbConnector, fmt.Errorf("Unknown test %d", test)
	}

	dbConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}
	Functions = MYSQLFunctions{}

	return &dataSet, dbConnector, nil
}

func createPrivilegesTestData() (*testhelpers.OrderedTests, DBConnectorMock, error) {
	dbConnector := DBConnectorMock{}
	dataSet := testhelpers.OrderedTests{
		orderedList: make(OrderedTestList, 0),
		testDataSet: make(TestDataSet, 0),
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, dbConnector, err
	}

	data := testhelpers.TestData{}

	privileges := make(models.Privileges, 2)
	privileges[0].ID = 0
	privileges[0].Name = "test0"
	privileges[0].Description = "description0"
	privileges[1].ID = 1
	privileges[1].Name = "test1"
	privileges[1].Description = "description1"

	testCase := "valid_id"
	data.expected.(map[string]interface{})["data"] = privileges
	data.expected.(map[string]interface{})["error"] = nil
	rows := sqlmock.NewRows([]string{"id", "name", "description"})
	for _, privilege := range privileges {
		rows.AddRow(privilege.ID, privilege.Name, privilege.Description)
	}

	mock.ExpectBegin()
	mock.ExpectQuery(GetPrivilegesQuery).WillReturnRows(rows)
	mock.ExpectCommit()
	dataSet.testDataSet[testCase] = data
	dataSet.orderedList = append(dataSet.orderedList, testCase)

	testCase = "invalid_id"
	data.expected.(map[string]interface{})["data"] = nil
	data.expected.(map[string]interface{})["error"] = sql.ErrNoRows

	mock.ExpectBegin()
	mock.ExpectQuery(GetPrivilegesQuery).WillReturnError(sql.ErrNoRows)
	mock.ExpectCommit()
	dataSet.testDataSet[testCase] = data
	dataSet.orderedList = append(dataSet.orderedList, testCase)

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
	for _, testCaseString := range dataSet.orderedList {
		tx, err := dbConnector.DB.Begin()
		if err != nil {
			t.Errorf("Failed to setup DB transaction %s", err)
			return
		}
		testCase := dataSet.testDataSet[testCaseString]
		product := testCase.data.(models.Product)

		err = Functions.AddProduct(&product, tx)
		if err != testCase.expected {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, testCase.expected)
			return
		}
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
	for _, testCaseString := range dataSet.orderedList {
		tx, err := dbConnector.DB.Begin()
		if err != nil {
			t.Errorf("Failed to setup DB transaction %s", err)
			return
		}
		testCase := dataSet.testDataSet[testCaseString]
		testData := testCase.data.(map[string]interface{})
		productID := testData["product_id"].(uuid.UUID)

		err = Functions.AddProductUsers(&productID, testData["product_users"].(models.ProductUsers), tx)
		if err != testCase.expected {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, testCase.expected)
			return
		}
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
	for _, testCaseString := range dataSet.orderedList {
		tx, err := dbConnector.DB.Begin()
		if err != nil {
			t.Errorf("Failed to setup DB transaction %s", err)
			return
		}
		testCase := dataSet.testDataSet[testCaseString]
		productID := testCase.data.(uuid.UUID)

		err = Functions.deleteProductUsersByProductID(&productID, tx)
		if err != testCase.expected {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, testCase.expected)
			return
		}
	}
}

func TestGetProductByID(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createProductsTestData(GetProductByIDTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.orderedList {
		testCase := dataSet.testDataSet[testCaseString]
		productID := testCase.data.(uuid.UUID)
		expectedData := testCase.expected.(map[string]interface{})["data"].(models.Product)
		expectedError := testCase.expected.(map[string]interface{})["error"].(error)

		output, err := Functions.GetProductByID(productID)
		if !cmp.Equal(output, &expectedData) {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, output, &expectedData)
			return
		}

		if err != expectedError {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, expectedError)
			return
		}
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
	for _, testCaseString := range dataSet.orderedList {
		tx, err := dbConnector.DB.Begin()
		if err != nil {
			t.Errorf("Failed to setup DB transaction %s", err)
			return
		}
		testCase := dataSet.testDataSet[testCaseString]
		productName := testCase.data.(string)
		expectedData := testCase.expected.(map[string]interface{})["data"].(models.Product)
		expectedError := testCase.expected.(map[string]interface{})["error"].(error)

		output, err := Functions.GetProductByName(productName, tx)
		if !cmp.Equal(output, &expectedData) {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, output, &expectedData)
			return
		}

		if err != expectedError {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, expectedError)
			return
		}
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
	for _, testCaseString := range dataSet.orderedList {
		tx, err := dbConnector.DB.Begin()
		if err != nil {
			t.Errorf("Failed to setup DB transaction %s", err)
			return
		}
		testCase := dataSet.testDataSet[testCaseString]
		userID := testCase.data.(uuid.UUID)
		expectedData := testCase.expected.(map[string]interface{})["data"].(models.UserProducts)
		expectedError := testCase.expected.(map[string]interface{})["error"].(error)

		output, err := Functions.GetUserProductIDs(userID, tx)
		if !cmp.Equal(output, &expectedData) {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, output, &expectedData)
			return
		}

		if err != expectedError {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, expectedError)
			return
		}
	}
}

func TestGetProductsByUserID(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createProductsTestData(GetProductsByUserIDTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.orderedList {
		testCase := dataSet.testDataSet[testCaseString]
		userID := testCase.data.(uuid.UUID)
		expectedData := testCase.expected.(map[string]interface{})["data"].(models.UserProducts)
		expectedError := testCase.expected.(map[string]interface{})["error"].(error)

		output, err := Functions.GetProductsByUserID(userID)
		if !cmp.Equal(output, &expectedData) {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, output, &expectedData)
			return
		}

		if err != expectedError {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, expectedError)
			return
		}
	}
}

func TestDeleteProduct(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createProductsTestData(DeleteProductTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.orderedList {
		testCase := dataSet.testDataSet[testCaseString]
		data := testCase.data.(uuid.UUID)

		err = Functions.DeleteProduct(&data)
		if err != testCase.expected {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, testCase.expected.(error))
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

	defer dbConnector.DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.orderedList {
		testCase := dataSet.testDataSet[testCaseString]
		expectedData := testCase.expected.(map[string]interface{})["data"].(models.Privileges)
		expectedError := testCase.expected.(map[string]interface{})["error"].(error)

		output, err := Functions.GetPrivileges()
		if !cmp.Equal(output, &expectedData) {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, output, &expectedData)
			return
		}

		if err != expectedError {
			t.Errorf("\n%s test failed.\n  Returned:\n %+v\n  Expected:\n %+v", testCaseString, err, expectedError)
			return
		}
	}
}
