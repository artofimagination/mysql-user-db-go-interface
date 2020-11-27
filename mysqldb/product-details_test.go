package mysqldb

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func createDetailsTestData() (*test.OrderedTests, DBConnectorMock, error) {
	dbConnector := DBConnectorMock{}
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}
	details := make(models.Details)

	detailsID, err := uuid.NewUUID()
	if err != nil {
		return nil, dbConnector, err
	}

	productDetails := models.ProductDetails{
		ID:      detailsID,
		Details: details,
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, dbConnector, err
	}

	binary, err := json.Marshal(productDetails.Details)
	if err != nil {
		return nil, dbConnector, err
	}

	data := test.Data{
		Data:     productDetails,
		Expected: nil,
	}

	testCase := "valid_product_details"
	mock.ExpectBegin()
	mock.ExpectExec(AddProductDetailsQuery).WithArgs(productDetails.ID, binary).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "failed_query"
	data.Expected = errors.New("This is a failure test")
	mock.ExpectBegin()
	mock.ExpectExec(AddProductDetailsQuery).WithArgs(productDetails.ID, binary).WillReturnError(data.Expected.(error))
	mock.ExpectRollback()
	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	dbConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}
	Functions = MYSQLFunctions{}

	return &dataSet, dbConnector, nil
}

func TestAddProductDetails(t *testing.T) {
	// Create test data
	dataSet, dbConnector, err := createDetailsTestData()
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
				t.Errorf("Failed to setup DB transaction: %s", err)
				return
			}
			testCase := dataSet.TestDataSet[testCaseString]
			var expectedError error
			if testCase.Expected != nil {
				expectedError = testCase.Expected.(error)
			}
			productDetails := testCase.Data.(models.ProductDetails)

			err = Functions.AddDetails(&productDetails, tx)
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, testCase.Expected)
				return
			}
		})
	}
}
