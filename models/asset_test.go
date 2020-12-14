package models

import (
	"fmt"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

const (
	SetFilePathTest = iota
	GetFilePathTest
	GetFieldTest
	NewAssetTest
)

func createAssetTestData(testID int) (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	asset := Asset{
		ID:      assetID,
		DataMap: make(DataMap),
	}

	referenceID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	UUIDImpl = &UUIDImplMock{
		uuidMock: referenceID,
	}
	baseAssetPath := "test/path"
	asset.DataMap[BaseAssetPath] = baseAssetPath

	switch testID {
	case SetFilePathTest:
		testCase := "valid"
		data := make(map[string]interface{})
		data["asset"] = asset
		data["asset_type"] = "testType"
		data["asset_extension"] = ".jpg"
		expected := make(map[string]interface{})
		expected["data"] = fmt.Sprintf("%s/%s.jpg", baseAssetPath, referenceID.String())
		expected["error"] = nil
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetFilePathTest:
		testCase := "valid_key_uuid"
		defaultPath := "default/default.jpg"
		expected := make(map[string]interface{})
		expected["data"] = fmt.Sprintf("%s/%s.jpg", baseAssetPath, referenceID.String())
		expected["error"] = nil
		data := make(map[string]interface{})
		data["asset_type"] = "testType"
		data["default_path"] = defaultPath
		asset.DataMap[data["asset_type"].(string)] = expected["data"].(string)
		data["asset"] = asset
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_asset_type"
		data = make(map[string]interface{})
		data["asset"] = asset
		data["default_path"] = defaultPath
		data["asset_type"] = "testType2"
		expected = make(map[string]interface{})
		expected["data"] = defaultPath
		expected["error"] = nil
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetFieldTest:
		testCase := "valid_url"
		expected := "https://success.com"
		data := make(map[string]interface{})
		data["asset_type"] = "testType"
		asset.DataMap[data["asset_type"].(string)] = expected
		data["asset"] = asset
		data["default_url"] = expected
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_url"
		defaultURL := "https://default.com"
		data = make(map[string]interface{})
		data["asset_type"] = "testType2"
		asset.DataMap[data["asset_type"].(string)] = defaultURL
		data["asset"] = asset
		data["default_url"] = defaultURL
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: defaultURL,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case NewAssetTest:
		testCase := "valid"
		UUIDImpl = &UUIDImplMock{
			uuidMock: asset.ID,
		}
		expected := make(map[string]interface{})
		expected["data"] = &asset
		expected["error"] = nil
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     asset.DataMap,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "nil_reference"
		var nilRef DataMap
		UUIDImpl = &UUIDImplMock{
			uuidMock: asset.ID,
		}
		expected = make(map[string]interface{})
		expected["data"] = nil
		expected["error"] = ErrAssetRefNotInitialised
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     nilRef,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	}

	Interface = &RepoInterface{}

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
			expectedData := testCase.Expected.(map[string]interface{})["data"].(string)
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}
			assetType := testCase.Data.(map[string]interface{})["asset_type"].(string)
			assetExtension := testCase.Data.(map[string]interface{})["asset_extension"].(string)
			asset := testCase.Data.(map[string]interface{})["asset"].(Asset)

			err = asset.SetFilePath(assetType, assetExtension)
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, testCase.Expected)
				return
			}

			if diff := pretty.Diff(asset.DataMap[assetType], expectedData); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, asset.DataMap[assetType], expectedData, diff)
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
			expectedData := testCase.Expected.(map[string]interface{})["data"].(string)
			defaultPath := testCase.Data.(map[string]interface{})["default_path"].(string)
			assetType := testCase.Data.(map[string]interface{})["asset_type"].(string)
			asset := testCase.Data.(map[string]interface{})["asset"].(Asset)

			output := asset.GetFilePath(assetType, defaultPath)
			if diff := pretty.Diff(output, expectedData); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData, diff)
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
			expectedData := testCase.Expected.(string)
			assetType := testCase.Data.(map[string]interface{})["asset_type"].(string)
			defaultURL := testCase.Data.(map[string]interface{})["default_url"].(string)
			asset := testCase.Data.(map[string]interface{})["asset"].(Asset)

			output := asset.GetField(assetType, defaultURL)
			if diff := pretty.Diff(output, expectedData); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData, diff)
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
			var expectedData *Asset
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*Asset)
			}
			references := testCase.Data.(DataMap)
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}

			output, err := Interface.NewAsset(
				references,
				func(*uuid.UUID) (string, error) {
					return "test/path", nil
				})
			if diff := pretty.Diff(output, expectedData); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData, diff)
				return
			}

			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError, diff)
				return
			}
		})
	}
}
