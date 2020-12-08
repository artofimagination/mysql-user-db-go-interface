package models

import (
	"fmt"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

const (
	SetImagePathTest = 0
	GetImagePathTest = 1
	GetURLTest       = 2
	NewAssetTest     = 3
)

func createAssetTestData(testID int) (*test.OrderedTests, error) {
	dataSet := test.OrderedTests{
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

	UUIDImpl = UUIDImplMock{
		uuidMock: referenceID,
	}
	baseAssetPath := "test/path"
	asset.DataMap[BaseAssetPath] = baseAssetPath

	switch testID {
	case SetImagePathTest:
		testCase := "valid"

		data := test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["asset"] = asset
		data.Data.(map[string]interface{})["asset_type"] = "testType"
		data.Expected.(map[string]interface{})["data"] = fmt.Sprintf("%s/%s.jpg", baseAssetPath, referenceID.String())
		data.Expected.(map[string]interface{})["error"] = nil
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetImagePathTest:
		testCase := "valid_key_uuid"

		data := test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		defaultPath := "default/default.jpg"
		data.Data.(map[string]interface{})["default_path"] = defaultPath
		data.Expected.(map[string]interface{})["data"] = fmt.Sprintf("%s/%s.jpg", baseAssetPath, referenceID.String())
		data.Expected.(map[string]interface{})["error"] = nil
		data.Data.(map[string]interface{})["asset_type"] = "testType"
		asset.DataMap[data.Data.(map[string]interface{})["asset_type"].(string)] = data.Expected.(map[string]interface{})["data"].(string)
		data.Data.(map[string]interface{})["asset"] = asset
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_asset_type"

		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: make(map[string]interface{}),
		}

		data.Data.(map[string]interface{})["asset"] = asset
		data.Data.(map[string]interface{})["default_path"] = defaultPath
		data.Data.(map[string]interface{})["asset_type"] = "testType2"
		data.Expected.(map[string]interface{})["data"] = defaultPath
		data.Expected.(map[string]interface{})["error"] = nil
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetURLTest:
		testCase := "valid_url"

		data := test.Data{
			Data:     make(map[string]interface{}),
			Expected: "https://success.com",
		}

		data.Data.(map[string]interface{})["asset_type"] = "testType"
		asset.DataMap[data.Data.(map[string]interface{})["asset_type"].(string)] = data.Expected.(string)
		data.Data.(map[string]interface{})["asset"] = asset
		data.Data.(map[string]interface{})["default_url"] = data.Expected.(string)
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_url"
		data = test.Data{
			Data:     make(map[string]interface{}),
			Expected: "https://default.com",
		}

		data.Data.(map[string]interface{})["asset_type"] = "testType2"
		asset.DataMap[data.Data.(map[string]interface{})["asset_type"].(string)] = data.Expected.(string)
		data.Data.(map[string]interface{})["asset"] = asset
		data.Data.(map[string]interface{})["default_url"] = data.Expected.(string)
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case NewAssetTest:
		testCase := "valid"
		data := test.Data{
			Data:     asset.DataMap,
			Expected: make(map[string]interface{}),
		}
		UUIDImpl = UUIDImplMock{
			uuidMock: asset.ID,
		}
		data.Expected.(map[string]interface{})["data"] = &asset
		data.Expected.(map[string]interface{})["error"] = nil
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "nil_reference"
		var nilRef DataMap
		data = test.Data{
			Data:     nilRef,
			Expected: make(map[string]interface{}),
		}
		UUIDImpl = UUIDImplMock{
			uuidMock: asset.ID,
		}
		data.Expected.(map[string]interface{})["data"] = nil
		data.Expected.(map[string]interface{})["error"] = ErrAssetRefNotInitialised
		dataSet.TestDataSet[testCase] = data
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	}

	Interface = RepoInterface{}

	return &dataSet, nil
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
			if !test.ErrEqual(err, expectedError) {
				t.Errorf(test.TestResultString, testCaseString, err, testCase.Expected)
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
			defaultPath := testCase.Data.(map[string]interface{})["default_path"].(string)
			assetType := testCase.Data.(map[string]interface{})["asset_type"].(string)
			asset := testCase.Data.(map[string]interface{})["asset"].(Asset)

			output := asset.GetImagePath(assetType, defaultPath)
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
			defaultURL := testCase.Data.(map[string]interface{})["default_url"].(string)
			asset := testCase.Data.(map[string]interface{})["asset"].(Asset)

			output := asset.GetURL(assetType, defaultURL)
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
				func(*uuid.UUID) (string, error) {
					return "test/path", nil
				})
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
