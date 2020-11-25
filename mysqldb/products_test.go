package mysqldb

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

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

func createTestProductUsersData() (models.ProductUsers, error) {
	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	owners := make(models.ProductUsers)
	owners[userID] = 1
	return owners, nil
}

func createTestUserProductsData(quantity int) (*models.UserProducts, error) {
	userProducts := make(models.UserProducts)
	for ; quantity > 0; quantity-- {
		productID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		userProducts[productID] = 1
	}
	return &userProducts, nil
}

func TestAddProduct_ValidProduct(t *testing.T) {
	// Create test data
	product, err := createTestProductData()
	if err != nil {
		t.Errorf("Failed to generate product data %s", err)
		return
	}

	jsonRaw, err := ConvertToJSONRaw(product.Details)
	if err != nil {
		t.Errorf("Failed to generate JSONRaw from details %s", err)
		return
	}

	// Create mock conditions
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate DB mock %s", err)
		return
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(AddProductQuery).WithArgs(product.ID, product.Name, product.Public, jsonRaw, product.AssetsID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}

	// Run test
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to setup DB transaction %s", err)
		return
	}

	if err := FunctionInterface.AddProduct(product, tx); err != nil {
		t.Errorf("Failed to add product %s", err)
		return
	}
}

func TestAddProductUsers_ValidProductTests(t *testing.T) {
	// Create test data
	product, err := createTestProductData()
	if err != nil {
		t.Errorf("Failed to generate product data %s", err)
		return
	}

	productUsers, err := createTestProductUsersData()
	if err != nil {
		t.Errorf("Failed to generate product users data %s", err)
		return
	}

	// Create mock conditions
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to generate DB mock %s", err)
		return
	}
	defer db.Close()

	mock.ExpectBegin()
	for userID, privilege := range productUsers {
		mock.ExpectExec("INSERT INTO users_products").WithArgs(product.ID, userID, privilege).WillReturnResult(sqlmock.NewResult(1, 1))
	}
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}

	FunctionInterface = MYSQLFunctionInterface{}

	// Run test
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to setup DB transaction %s", err)
		return
	}

	if err := FunctionInterface.AddProductUsers(&product.ID, productUsers, tx); err != nil {
		t.Errorf("Failed to add product users %s", err)
		return
	}
}

func TestDeleteProductUsersByProductID_ValidID(t *testing.T) {
	// Create test data
	productID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to generate product UUID %s", err)
		return
	}

	// Create mock conditions
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to generate DB mock %s", err)
		return
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM users_products").WithArgs(productID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}

	// Run test
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to create transaction %s", err)
		return
	}

	if err := deleteProductUsersByProductID(&productID, tx); err != nil {
		t.Errorf("Failed to delete product users %s", err)
		return
	}
}

func TestGetProductByID_ValidID(t *testing.T) {
	// Create test data
	product, err := createTestProductData()
	if err != nil {
		t.Errorf("Failed to generate product data %s", err)
		return
	}

	jsonRaw, err := ConvertToJSONRaw(product.Details)
	if err != nil {
		t.Errorf("Failed to generate JSONRaw from details %s", err)
		return
	}

	binaryID, err := json.Marshal(product.ID)
	if err != nil {
		t.Errorf("Failed to generate binary from UUID %s", err)
		return
	}

	// Create mock conditions
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate DB mock %s", err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "public", "details", "product_assets_id"}).
		AddRow(binaryID, product.Name, product.Public, jsonRaw, product.AssetsID)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT BIN_TO_UUID(id), name, public, details, BIN_TO_UUID(product_assets_id) FROM products WHERE id = UUID_TO_BIN(?)").WithArgs(product.ID).WillReturnRows(rows)
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}

	// Run test
	data, err := GetProductByID(product.ID)
	if err != nil {
		t.Errorf("Failed to get product %s", err)
		return
	}

	if !cmp.Equal(*data, *product) {
		t.Errorf("Test returned:\n %+v\nExpected:\n %+v", *data, *product)
		return
	}
}

func TestGetProductByName_ValidName(t *testing.T) {
	// Create test data
	product, err := createTestProductData()
	if err != nil {
		t.Errorf("Failed to generate product data %s", err)
		return
	}

	jsonRaw, err := ConvertToJSONRaw(product.Details)
	if err != nil {
		t.Errorf("Failed to generate JSONRaw from details %s", err)
		return
	}

	binaryID, err := json.Marshal(product.ID)
	if err != nil {
		t.Errorf("Failed to generate binary from UUID %s", err)
		return
	}

	// Create mock conditions
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate DB mock %s", err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "public", "details", "product_assets_id"}).
		AddRow(binaryID, product.Name, product.Public, jsonRaw, product.AssetsID)
	mock.ExpectBegin()
	mock.ExpectQuery(GetProductByNameQuery).WithArgs(product.Name).WillReturnRows(rows)
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}

	FunctionInterface = MYSQLFunctionInterface{}

	// Run test
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to create transaction %s", err)
		return
	}

	data, err := FunctionInterface.GetProductByName(product.Name, tx)
	if err != nil {
		t.Errorf("Failed to get product %s", err)
		return
	}

	if !cmp.Equal(*data, *product) {
		t.Errorf("Test returned:\n %+v\nExpected:\n %+v", *data, *product)
		return
	}
}

func TestGetUserProductIDs_ValidID(t *testing.T) {
	// Create test data
	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to generate user UUID %s", err)
	}

	userProducts, err := createTestUserProductsData(2)
	if err != nil {
		t.Errorf("Failed to generate user products data %s", err)
		return
	}

	// Create mock conditions
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate DB mock %s", err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"products_id", "privilege"})
	for productID, privilege := range *userProducts {
		rows.AddRow(productID, privilege)
	}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT BIN_TO_UUID(products_id), privilege FROM users_products where users_id = UUID_TO_BIN(?)").WithArgs(userID).WillReturnRows(rows)
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}

	FunctionInterface = MYSQLFunctionInterface{}

	// Run test
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to create transaction %s", err)
		return
	}

	data, err := getUserProductIDs(userID, tx)
	if err != nil {
		t.Errorf("Failed to get product %s", err)
		return
	}

	if !cmp.Equal(*data, *userProducts) {
		t.Errorf("Test returned:\n %+v\nExpected:\n %+v", *data, *userProducts)
		return
	}
}

func TestGetProductsByUserID_ValidID(t *testing.T) {
	// Create test data
	products, err := createTestProductList(2)
	if err != nil {
		t.Errorf("Failed to generate product list %s", err)
		return
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to generate user UUID %s", err)
		return
	}

	userProducts := make(models.UserProducts)
	for _, product := range *products {
		userProducts[product.ID] = 1
	}

	// Create mock conditions
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate DB mock %s", err)
		return
	}

	rowsUserProducts := sqlmock.NewRows([]string{"products_id", "privilege"})
	// Ranging through map is random, need to collect product ID-s ina fixed order in order to
	// have the correct order of sql mock expectations.
	orderedProductIDs := make([]uuid.UUID, 0)
	for productID, privilege := range userProducts {

		rowsUserProducts.AddRow(productID, privilege)
		orderedProductIDs = append(orderedProductIDs, productID)
	}

	rowsProducts, err := addProductsToMock(products)
	if err != nil {
		t.Errorf("Failed to generate product mock rows %s", err)
		return
	}

	mock.ExpectBegin()
	mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnRows(rowsUserProducts)
	for _, productID := range orderedProductIDs {
		mock.ExpectQuery(GetProductByIDQuery).WithArgs(productID).WillReturnRows(rowsProducts)
	}
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}

	// Run test
	data, err := GetProductsByUserID(userID)
	if err != nil {
		t.Errorf("Failed to get product %s", err)
		return
	}

	if !cmp.Equal(*data, *products) {
		t.Errorf("Test returned:\n %+v\nExpected:\n %+v", *data, *products)
		return
	}
}

func TestDeleteProduct(t *testing.T) {
	// Create test data
	products, err := createTestProductList(1)
	if err != nil {
		t.Errorf("Failed to generate product list %s", err)
		return
	}

	// Create mock conditions
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate DB mock %s", err)
		return
	}

	mock.ExpectBegin()
	for _, product := range *products {
		mock.ExpectExec("DELETE FROM users_products where products_id = UUID_TO_BIN(?)").WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("DELETE FROM products where id = UUID_TO_BIN(?)").WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	}
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}

	// Run test
	for _, product := range *products {
		if err := DeleteProduct(&product.ID); err != nil {
			t.Errorf("Failed to delete product %s", err)
		}
	}
}

func TestGetPrivileges_ValidRows(t *testing.T) {
	// Create test data
	expected := make(models.Privileges, 2)
	expected[0].ID = 0
	expected[0].Name = "test0"
	expected[0].Description = "description0"
	expected[1].ID = 1
	expected[1].Name = "test1"
	expected[1].Description = "description1"

	// Create mock conditions
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate DB mock %s", err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description"})
	for _, privilege := range expected {
		rows.AddRow(privilege.ID, privilege.Name, privilege.Description)
	}

	mock.ExpectBegin()
	mock.ExpectQuery(GetPrivilegesQuery).WillReturnRows(rows)
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}

	FunctionInterface = MYSQLFunctionInterface{}

	// Run test
	data, err := FunctionInterface.GetPrivileges()
	if err != nil {
		t.Errorf("Failed to get product %s", err)
		return
	}

	if !cmp.Equal(data, expected) {
		t.Errorf("Test returned:\n %+v\nExpected:\n %+v", data, expected)
		return
	}
}

func TestGetPrivileges_NoRows(t *testing.T) {
	// Create mock conditions
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate DB mock %s", err)
		return
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(GetPrivilegesQuery).WillReturnError(sql.ErrNoRows)
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}

	FunctionInterface = MYSQLFunctionInterface{}

	// Run test
	_, err = FunctionInterface.GetPrivileges()
	if err == nil || (err != nil && err != sql.ErrNoRows) {
		t.Errorf("Test returned:\n %+v\nExpected:\n %+v", err, sql.ErrNoRows)
		return
	}
}
