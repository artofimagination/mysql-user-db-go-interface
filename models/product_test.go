package models

import (
	"errors"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/tests"
	"github.com/google/uuid"
)

const (
	NewProduct = iota
)

type ProductExpectedData struct {
	product *Product
	err     error
}

type ProductInputData struct {
	product *Product
}

func createTestData(testID int) (*tests.OrderedTests, error) {
	dataSet := &tests.OrderedTests{
		OrderedList: make(tests.OrderedTestList, 0),
		TestDataSet: make(tests.DataSet),
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

	ModelFunctions = &RepoFunctions{
		UUIDImpl: &UUIDImplMock{
			uuidMock: productID,
		},
	}

	product := &Product{
		ID:        productID,
		Name:      "TestProduct",
		AssetsID:  assetsID,
		DetailsID: detailsID,
	}

	switch testID {
	case NewProduct:
		testCase := "valid_product"
		expected := ProductExpectedData{
			product: product,
			err:     nil,
		}
		input := ProductInputData{
			product: product,
		}
		dataSet.TestDataSet[testCase] = tests.Data{
			Data:     input,
			Mock:     nil,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failure_case"
		expected = ProductExpectedData{
			product: nil,
			err:     errors.New("Failed with error"),
		}
		input = ProductInputData{
			product: product,
		}
		dataSet.TestDataSet[testCase] = tests.Data{
			Data:     input,
			Mock:     nil,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	}

	return dataSet, nil
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
			expectedData := testCase.Expected.(ProductExpectedData)
			inputData := testCase.Data.(ProductInputData)
			ModelFunctions.UUIDImpl.(*UUIDImplMock).err = expectedData.err

			output, err := ModelFunctions.NewProduct(
				inputData.product.Name,
				&inputData.product.DetailsID,
				&inputData.product.AssetsID,
			)
			tests.CheckResult(output, expectedData.product, err, expectedData.err, testCaseString, t)
		})
	}
}
