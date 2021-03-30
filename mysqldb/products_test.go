package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/tests"
	"github.com/google/uuid"
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

type ProductInputData struct {
	userID       *uuid.UUID
	product      *models.Product
	productUsers *models.ProductUserIDs
	privilege    int
}

type ProductExpectedData struct {
	userProducts *models.UserProductIDs
	product      *models.Product
	err          error
}

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

func createProductsTestData(testID int) (*tests.OrderedTests, error) {
	dataSet := &tests.OrderedTests{
		OrderedList: make(tests.OrderedTestList, 0),
		TestDataSet: make(tests.DataSet),
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
		mock.ExpectExec(AddProductQuery).WithArgs(product.ID, product.Name, product.DetailsID, product.AssetsID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product: product,
			},
			Expected: ProductExpectedData{
				err: nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		expected := errors.New("This is a failure test")
		mock.ExpectBegin()
		mock.ExpectExec(AddProductQuery).WithArgs(product.ID, product.Name, product.DetailsID, product.AssetsID).WillReturnError(expected)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product: product,
			},
			Expected: ProductExpectedData{
				err: expected,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "duplicate_name"
		expected = fmt.Errorf(ErrSQLDuplicateProductNameEntryString, product.Name)
		mock.ExpectBegin()
		mock.ExpectExec(AddProductQuery).WithArgs(product.ID, product.Name, product.DetailsID, product.AssetsID).WillReturnError(expected)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product: product,
			},
			Expected: ProductExpectedData{
				err: expected,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case addProductUsersTest:
		testCase := "valid_products"
		mock.ExpectBegin()
		for _, userID := range productUsers.UserIDArray {
			privilege := productUsers.UserMap[userID]
			mock.ExpectExec(AddProductUsersQuery).WithArgs(userID, product.ID, privilege).WillReturnResult(sqlmock.NewResult(1, 1))
		}
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product:      product,
				productUsers: productUsers,
			},
			Expected: ProductExpectedData{
				err: nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		expected := errors.New("This is a failure test")
		mock.ExpectBegin()
		for _, userID := range productUsers.UserIDArray {
			privilege := productUsers.UserMap[userID]
			mock.ExpectExec(AddProductUsersQuery).WithArgs(userID, product.ID, privilege).WillReturnError(expected)
		}
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product:      product,
				productUsers: productUsers,
			},
			Expected: ProductExpectedData{
				err: expected,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_to_add"
		mock.ExpectBegin()
		for _, userID := range productUsers.UserIDArray {
			privilege := productUsers.UserMap[userID]
			mock.ExpectExec(AddProductUsersQuery).WithArgs(userID, product.ID, privilege).WillReturnResult(sqlmock.NewResult(1, 0))
		}
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product:      product,
				productUsers: productUsers,
			},
			Expected: ProductExpectedData{
				err: ErrNoProductUserAdded,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case deleteProductUsersByProductIDTest:
		testCase := "valid_id"
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductUsersByProductIDQuery).WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product: product,
			},
			Expected: ProductExpectedData{
				err: nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_id"
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductUsersByProductIDQuery).WithArgs(product.ID).WillReturnError(ErrNoUserWithProduct)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product: product,
			},
			Expected: ProductExpectedData{
				err: ErrNoUserWithProduct,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case getProductByIDTest:

		testCase := "valid_id"
		rows := sqlmock.NewRows([]string{"id", "name", "details", "product_assets_id"}).
			AddRow(binaryProductID, product.Name, product.DetailsID, product.AssetsID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByIDQuery).WithArgs(product.ID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product: product,
			},
			Expected: ProductExpectedData{
				product: product,
				err:     nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_id"
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByIDQuery).WithArgs(product.ID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product: product,
			},
			Expected: ProductExpectedData{
				product: nil,
				err:     sql.ErrNoRows,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case getProductByNameTest:
		testCase := "valid_name"
		rows := sqlmock.NewRows([]string{"id", "name", "details", "product_assets_id"}).
			AddRow(binaryProductID, product.Name, product.DetailsID, binaryAssetID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByNameQuery).WithArgs(product.Name).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product: product,
			},
			Expected: ProductExpectedData{
				product: product,
				err:     nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_name"
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductByNameQuery).WithArgs(product.Name).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product: product,
			},
			Expected: ProductExpectedData{
				product: nil,
				err:     sql.ErrNoRows,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case GetUserProductIDsTest:
		testCase := "valid_id"
		rows := sqlmock.NewRows([]string{"products_id", "privilege"})
		for _, productID := range userProducts.ProductIDArray {
			rows.AddRow(productID, userProducts.ProductMap[productID])
		}
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				userID: &userID,
			},
			Expected: ProductExpectedData{
				userProducts: userProducts,
				err:          nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_products"
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProductIDsQuery).WithArgs(userID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				userID: &userID,
			},
			Expected: ProductExpectedData{
				userProducts: nil,
				err:          sql.ErrNoRows,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case deleteProductTest:
		testCase := "valid_id"
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductQuery).WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product: product,
			},
			Expected: ProductExpectedData{
				err: nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "no_product"
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductQuery).WithArgs(product.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				product: product,
			},
			Expected: ProductExpectedData{
				err: ErrNoProductDeleted,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case updateUsersProductsTest:
		testCase := "valid_id"
		privilege := 1
		mock.ExpectBegin()
		mock.ExpectExec(UpdateUsersProductsQuery).WithArgs(privilege, userID, product.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				userID:    &userID,
				product:   product,
				privilege: privilege,
			},
			Expected: ProductExpectedData{
				err: nil,
			},
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "no_users_products"
		privilege = 1
		mock.ExpectBegin()
		mock.ExpectExec(UpdateUsersProductsQuery).WithArgs(privilege, userID, product.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = tests.Data{
			Data: ProductInputData{
				userID:    &userID,
				product:   product,
				privilege: privilege,
			},
			Expected: ProductExpectedData{
				err: ErrNoUsersProductUpdate,
			},
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
			expectedData := testCase.Expected.(ProductExpectedData)
			inputData := testCase.Data.(ProductInputData)

			err = DBFunctions.AddProduct(inputData.product, tx)
			tests.CheckResult(nil, nil, err, expectedData.err, testCaseString, t)
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
			expectedData := testCase.Expected.(ProductExpectedData)
			inputData := testCase.Data.(ProductInputData)
			err = DBFunctions.AddProductUsers(&inputData.product.ID, inputData.productUsers, tx)
			tests.CheckResult(nil, nil, err, expectedData.err, testCaseString, t)
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
			expectedData := testCase.Expected.(ProductExpectedData)
			inputData := testCase.Data.(ProductInputData)

			err = DBFunctions.UpdateUsersProducts(inputData.userID, &inputData.product.ID, inputData.privilege, tx)
			tests.CheckResult(nil, nil, err, expectedData.err, testCaseString, t)
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
			expectedData := testCase.Expected.(ProductExpectedData)
			inputData := testCase.Data.(ProductInputData)

			err = DBFunctions.DeleteProductUsersByProductID(&inputData.product.ID, tx)
			tests.CheckResult(nil, nil, err, expectedData.err, testCaseString, t)
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
			expectedData := testCase.Expected.(ProductExpectedData)
			inputData := testCase.Data.(ProductInputData)

			output, err := DBFunctions.GetProductByID(&inputData.product.ID, tx)
			tests.CheckResult(output, expectedData.product, err, expectedData.err, testCaseString, t)
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
			expectedData := testCase.Expected.(ProductExpectedData)
			inputData := testCase.Data.(ProductInputData)

			output, err := DBFunctions.GetProductByName(inputData.product.Name, tx)
			tests.CheckResult(output, expectedData.product, err, expectedData.err, testCaseString, t)
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
			expectedData := testCase.Expected.(ProductExpectedData)
			inputData := testCase.Data.(ProductInputData)

			output, err := DBFunctions.GetUserProductIDs(inputData.userID, tx)
			tests.CheckResult(output, expectedData.userProducts, err, expectedData.err, testCaseString, t)
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
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}

			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(ProductExpectedData)
			inputData := testCase.Data.(ProductInputData)

			err = DBFunctions.DeleteProduct(&inputData.product.ID, tx)
			tests.CheckResult(nil, nil, err, expectedData.err, testCaseString, t)
		})
	}
}
