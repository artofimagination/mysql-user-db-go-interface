package models

import (
	"fmt"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/tests"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

const (
	SetFilePathTest = iota
	GetFilePathTest
	GetFieldTest
	NewAssetTest
)

type AssetExpectedData struct {
	asset *Asset
	data  string
	err   error
}

type AssetInputData struct {
	asset          Asset
	assetType      string
	assetExtension string
	defaultData    string
}

func createAssetTestData(testID int) (*tests.OrderedTests, error) {
	dataSet := &tests.OrderedTests{
		OrderedList: make(tests.OrderedTestList, 0),
		TestDataSet: make(tests.DataSet),
	}

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	asset := &Asset{
		ID:      assetID,
		DataMap: make(DataMap),
	}

	referenceID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	ModelFunctions = &RepoFunctions{
		UUIDImpl: &UUIDImplMock{
			uuidMock: referenceID,
		},
	}

	baseAssetPath := "test/path"
	asset.DataMap[BaseAssetPath] = baseAssetPath

	switch testID {
	case SetFilePathTest:
		testCase := "valid"
		expected := AssetExpectedData{
			data: fmt.Sprintf("%s/%s.jpg", baseAssetPath, referenceID.String()),
			err:  nil,
		}
		input := AssetInputData{
			asset:          *asset,
			assetType:      "testType",
			assetExtension: ".jpg",
		}

		dataSet.TestDataSet[testCase] = tests.Data{
			Data:     input,
			Mock:     nil,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetFilePathTest:
		testCase := "valid_key_uuid"
		defaultPath := "default/default.jpg"
		assetType := "testType"
		returnPath := fmt.Sprintf("%s/%s.jpg", baseAssetPath, referenceID.String())
		asset.DataMap[assetType] = returnPath
		expected := AssetExpectedData{
			data: returnPath,
			err:  nil,
		}
		input := AssetInputData{
			asset:       *asset,
			assetType:   assetType,
			defaultData: defaultPath,
		}

		dataSet.TestDataSet[testCase] = tests.Data{
			Data:     input,
			Mock:     nil,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_asset_type"
		expected = AssetExpectedData{
			data: defaultPath,
			err:  nil,
		}
		input = AssetInputData{
			asset:       *asset,
			assetType:   "testType2",
			defaultData: defaultPath,
		}

		dataSet.TestDataSet[testCase] = tests.Data{
			Data:     input,
			Mock:     nil,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetFieldTest:
		testCase := "valid_url"
		url := "https://success.com"
		assetType := "testType"
		asset.DataMap[assetType] = url
		expected := AssetExpectedData{
			data: url,
		}
		input := AssetInputData{
			asset:       *asset,
			assetType:   assetType,
			defaultData: url,
		}

		dataSet.TestDataSet[testCase] = tests.Data{
			Data:     input,
			Mock:     nil,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_url"
		defaultURL := "https://default.com"
		assetType = "testType2"
		asset.DataMap[assetType] = defaultURL
		expected = AssetExpectedData{
			data: defaultURL,
		}
		input = AssetInputData{
			asset:       *asset,
			assetType:   assetType,
			defaultData: defaultURL,
		}

		dataSet.TestDataSet[testCase] = tests.Data{
			Data:     input,
			Mock:     nil,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case NewAssetTest:
		testCase := "valid"
		ModelFunctions.UUIDImpl = &UUIDImplMock{
			uuidMock: asset.ID,
		}
		expected := AssetExpectedData{
			asset: &Asset{
				ID:      assetID,
				DataMap: make(DataMap),
			},
			err: nil,
		}
		expected.asset.DataMap[BaseAssetPath] = baseAssetPath
		input := AssetInputData{
			asset: *asset,
		}

		dataSet.TestDataSet[testCase] = tests.Data{
			Data:     input,
			Mock:     nil,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "nil_reference"
		ModelFunctions.UUIDImpl = &UUIDImplMock{
			uuidMock: asset.ID,
		}
		expected = AssetExpectedData{
			asset: nil,
			err:   ErrAssetRefNotInitialised,
		}
		asset.DataMap = nil
		input = AssetInputData{
			asset: *asset,
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

func TestSetFilePath(t *testing.T) {
	// Create test data
	dataSet, err := createAssetTestData(SetFilePathTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(AssetExpectedData)
			inputData := testCase.Data.(AssetInputData)

			err = ModelFunctions.SetFilePath(&inputData.asset, inputData.assetType, inputData.assetExtension)
			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(tests.TestResultString, testCaseString, err, expectedData.err)
				return
			}

			if diff := pretty.Diff(inputData.asset.DataMap[inputData.assetType], expectedData.data); len(diff) != 0 {
				t.Errorf(tests.TestResultString, testCaseString, inputData.asset.DataMap[inputData.assetType], expectedData.data, diff)
				return
			}
		})
	}
}

func TestGetFilePath(t *testing.T) {
	// Create test data
	dataSet, err := createAssetTestData(GetFilePathTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(AssetExpectedData)
			inputData := testCase.Data.(AssetInputData)

			output := ModelFunctions.GetFilePath(&inputData.asset, inputData.assetType, inputData.defaultData)
			if diff := pretty.Diff(output, expectedData.data); len(diff) != 0 {
				t.Errorf(tests.TestResultString, testCaseString, output, expectedData.data, diff)
				return
			}
		})
	}
}

func TestGetField(t *testing.T) {
	// Create test data
	dataSet, err := createAssetTestData(GetFieldTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(AssetExpectedData)
			inputData := testCase.Data.(AssetInputData)

			output := ModelFunctions.GetField(&inputData.asset, inputData.assetType, inputData.defaultData)
			if diff := pretty.Diff(output, expectedData.data); len(diff) != 0 {
				t.Errorf(tests.TestResultString, testCaseString, output, expectedData.data, diff)
				return
			}
		})
	}
}

func TestNewAsset(t *testing.T) {
	// Create test data
	dataSet, err := createAssetTestData(NewAssetTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(AssetExpectedData)
			inputData := testCase.Data.(AssetInputData)

			output, err := ModelFunctions.NewAsset(
				inputData.asset.DataMap,
				func(*uuid.UUID) (string, error) {
					return "test/path", nil
				})

			tests.CheckResult(output, expectedData.asset, err, expectedData.err, testCaseString, t)
		})
	}
}
