package controllers

import (
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func createUserTestData() (*test.OrderedTests, error) {
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	settingsID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	models.Interface = ModelInterfaceMock{
		assetID:    assetID,
		settingsID: settingsID,
		userID:     userID,
	}

	user := models.User{
		Name:       "testName",
		Email:      "testEmail",
		Password:   []byte{},
		ID:         userID,
		SettingsID: assetID,
		AssetsID:   assetID,
	}

	testCase := "no_existing_user"
	data := test.Data{
		Data:     make(map[string]interface{}),
		Expected: make(map[string]interface{}),
	}

	data.Data.(map[string]interface{})["input"] = user
	data.Data.(map[string]interface{})["db_mock"] = nil
	data.Expected.(map[string]interface{})["data"] = &user
	data.Expected.(map[string]interface{})["error"] = nil
	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "existing_user"
	data = test.Data{
		Data:     make(map[string]interface{}),
		Expected: make(map[string]interface{}),
	}

	data.Data.(map[string]interface{})["input"] = user
	data.Data.(map[string]interface{})["db_mock"] = &user
	data.Expected.(map[string]interface{})["data"] = nil
	data.Expected.(map[string]interface{})["error"] = ErrDuplicateEmailEntry
	dataSet.TestDataSet[testCase] = data
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	mysqldb.Functions = DBFunctionInterfaceMock{}
	mysqldb.DBConnector = DBConnectorMock{}
	return &dataSet, nil
}

func TestCreateUser(t *testing.T) {
	// Create test data
	dataSet, err := createUserTestData()
	if err != nil {
		t.Errorf("Failed to create test data %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			var expectedData *models.User
			if testCase.Expected.(map[string]interface{})["data"] != nil {
				expectedData = testCase.Expected.(map[string]interface{})["data"].(*models.User)
			}
			var expectedError error
			if testCase.Expected.(map[string]interface{})["error"] != nil {
				expectedError = testCase.Expected.(map[string]interface{})["error"].(error)
			}
			input := testCase.Data.(map[string]interface{})["input"].(models.User)
			var DBMock *models.User
			if testCase.Data.(map[string]interface{})["db_mock"] != nil {
				DBMock = testCase.Data.(map[string]interface{})["db_mock"].(*models.User)
			}

			mysqldb.Functions = DBFunctionInterfaceMock{
				user: DBMock,
			}

			output, err := CreateUser(
				input.Name,
				input.Email,
				input.Password,
				func(*uuid.UUID) string {
					return "testPath"
				}, func([]byte) ([]byte, error) {
					return []byte{}, nil
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
