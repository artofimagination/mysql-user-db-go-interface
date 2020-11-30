package models

import (
	"errors"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

const (
	NewProduct = 0
)

func createTestData(testID int) (*test.OrderedTests, error) {
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	assetsID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	detailsID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	UUIDImpl = UUIDImplMock{
		uuidMock: productID,
	}

	product := Product{
		ID:        productID,
		Name:      "TestProduct",
		AssetsID:  assetsID,
		DetailsID: detailsID,
		Public:    true,
	}

	switch testID {
	case NewProduct:
		testCase := "valid_product"

		data := test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["product"] = product
		data.Data.(map[string]interface{})["uuid_mock"] = UUIDImplMock{
			uuidMock: productID,
		}
		data.Expected.(map[string]interface{})["data"] = &product
		data.Expected.(map[string]interface{})["error"] = nil
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failure_case"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		err := errors.New("Failed with error")
		data.Data.(map[string]interface{})["product"] = product
		data.Data.(map[string]interface{})["uuid_mock"] = UUIDImplMock{
			err: err,
		}
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = err
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	}

	Interface = RepoInterface{}

	return &dataSet, nil
}

func TestNewProduct(t *testing.T) {
	// Create test data
	dataSet, err := createTestData(NewProduct)
	if err != nil {
		t.Errorf("Failed to create test data %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}
			var expectedData *Product
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*Product)
			}

			UUIDImpl = testCase.Data.(map[string]interface{})["uuid_mock"].(UUIDImplMock)
			inputData := testCase.Data.(map[string]interface{})["product"].(Product)

			output, err := Interface.NewProduct(
				inputData.Name,
				inputData.Public,
				&inputData.DetailsID,
				&inputData.AssetsID,
			)
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
