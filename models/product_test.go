package models

import (
	"errors"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/kr/pretty"
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

func createTestData(testID int) (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
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
		dataSet.TestDataSet[testCase] = test.Data{
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
		dataSet.TestDataSet[testCase] = test.Data{
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
			if diff := pretty.Diff(output, expectedData.product); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData.product, diff)
				return
			}

			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
				return
			}
		})
	}
}
