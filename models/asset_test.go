package models

import (
	"fmt"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

const (
	SetImagePathTest = iota
	GetImagePathTest
	GetURLTest
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
	case SetImagePathTest:
		testCase := "valid"
		data := make(map[string]interface{})
		data["asset"] = asset
		data["asset_type"] = "testType"
		expected := make(map[string]interface{})
		expected["data"] = fmt.Sprintf("%s/%s.jpg", baseAssetPath, referenceID.String())
		expected["error"] = nil
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetImagePathTest:
		testCase := "valid_key_uuid"
		expected := make(map[string]interface{})
		expected["data"] = fmt.Sprintf("%s/%s.jpg", baseAssetPath, referenceID.String())
		expected["error"] = nil
		data := make(map[string]interface{})
		data["asset_type"] = "testType"
		asset.DataMap[data["asset_type"].(string)] = expected["data"].(string)
		data["asset"] = asset
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_asset_type"
		DefaultImagePath = "default/default.jpg"
		data = make(map[string]interface{})
		data["asset"] = asset
		data["asset_type"] = "testType2"
		expected = make(map[string]interface{})
		expected["data"] = DefaultImagePath
		expected["error"] = nil
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetURLTest:
		testCase := "valid_url"
		expected := "https://success.com"
		data := make(map[string]interface{})
		data["asset_type"] = "testType"
		asset.DataMap[data["asset_type"].(string)] = expected
		data["asset"] = asset
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_url"
		DefaultURL = "https://default.com"
		data = make(map[string]interface{})
		data["asset_type"] = "testType2"
		asset.DataMap[data["asset_type"].(string)] = DefaultURL
		data["asset"] = asset
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: DefaultURL,
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

func TestSetImagePath(t *testing.T) {
	// Create test data
	dataSet, err := createAssetTestData(SetImagePathTest)
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
			asset := testCase.Data.(map[string]interface{})["asset"].(Asset)

			err = asset.SetImagePath(assetType)
			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}

			if asset.DataMap[assetType] != expectedData {
				t.Errorf(test.TestResultString, testCaseString, asset.DataMap[assetType], expectedData)
				return
			}
		})
	}
}

func TestGetImagePath(t *testing.T) {
	// Create test data
	dataSet, err := createAssetTestData(GetImagePathTest)
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
			assetType := testCase.Data.(map[string]interface{})["asset_type"].(string)
			asset := testCase.Data.(map[string]interface{})["asset"].(Asset)

			output := asset.GetImagePath(assetType)
			if output != expectedData {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData)
				return
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	// Create test data
	dataSet, err := createAssetTestData(GetURLTest)
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
			asset := testCase.Data.(map[string]interface{})["asset"].(Asset)

			output := asset.GetURL(assetType)
			if output != expectedData {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData)
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
				func(*uuid.UUID) string {
					return "test/path"
				})
			if diff := pretty.Diff(output, expectedData); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData)
				return
			}

			if diff := pretty.Diff(err, expectedError); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedError)
				return
			}
		})
	}
}
