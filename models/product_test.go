package models

import (
	"errors"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

const (
	NewProduct = iota
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

		data := make(map[string]interface{})
		data["product"] = product
		data["uuid_mock"] = UUIDImplMock{
			uuidMock: productID,
		}
		expected := make(map[string]interface{})
		expected["data"] = &product
		expected["error"] = nil
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failure_case"

		err := errors.New("Failed with error")
		data = make(map[string]interface{})
		data["product"] = product
		data["uuid_mock"] = UUIDImplMock{
			err: err,
		}
		expected = make(map[string]interface{})
		expected["data"] = nil
		expected["error"] = err
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
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
