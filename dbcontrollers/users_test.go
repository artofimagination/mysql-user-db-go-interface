package dbcontrollers

import (
	"database/sql"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

const (
	CreateUserTest = iota
	DeleteUserTest
)

type UserExpectedData struct {
	userData             *models.UserData
	userDeleted          bool
	productDeleted       bool
	usersProductsUpdated bool
	err                  error
}

type UserMockData struct {
	user          *models.User
	product       *models.Product
	usersProducts *models.UserProductIDs
	privileges    models.Privileges
	err           error
}

type UserInputData struct {
	userData *models.UserData
	password []byte
	nominees map[uuid.UUID]uuid.UUID
}

func createUserTestData(testID int) (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
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

	user := &models.User{
		Name:       "testName",
		Email:      "testEmail",
		Password:   []byte{},
		ID:         userID,
		SettingsID: assetID,
		AssetsID:   assetID,
	}

	nomineeID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	usersProducts := &models.UserProductIDs{
		ProductMap: make(map[uuid.UUID]int),
	}

	product := &models.Product{
		ID:     productID,
	}

	dataMap := make(models.DataMap)
	assets := &models.Asset{
		ID:      assetID,
		DataMap: dataMap,
	}

	userData := &models.UserData{
		ID:       user.ID,
		Name:     user.Name,
		Settings: assets,
		Assets:   assets,
	}

	privileges := make(models.Privileges, 2)
	privilege := &models.Privilege{
		ID:          0,
		Name:        "Owner",
		Description: "description0",
	}
	privileges[0] = privilege
	privilege = &models.Privilege{
		ID:          1,
		Name:        "User",
		Description: "description1",
	}
	privileges[1] = privilege

	dbController = &MYSQLController{
		DBFunctions: &DBFunctionMock{},
		DBConnector: &DBConnectorMock{},
		ModelFunctions: &ModelMock{
			assetID:    assetID,
			settingsID: settingsID,
			userID:     userID,
			asset:      assets,
		},
	}

	switch testID {
	case CreateUserTest:
		testCase := "no_existing_user"
		expected := UserExpectedData{
			userData: userData,
			err:      nil,
		}
		input := UserInputData{
			userData: userData,
			password: user.Password,
		}
		mock := UserMockData{
			user: nil,
		}
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     input,
			Mock:     mock,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "existing_user"
		expected = UserExpectedData{
			userData: nil,
			err:      ErrDuplicateEmailEntry,
		}
		input = UserInputData{
			userData: userData,
			password: user.Password,
		}
		mock = UserMockData{
			user: user,
		}
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     input,
			Mock:     mock,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case DeleteUserTest:
		testCase := "valid_data_has_nominee"
		usersProducts.ProductMap[productID] = 0
		nominees := make(map[uuid.UUID]uuid.UUID)
		nominees[productID] = nomineeID
		expected := UserExpectedData{
			userDeleted:          true,
			productDeleted:       false,
			usersProductsUpdated: true,
			err:                  nil,
		}
		input := UserInputData{
			userData: userData,
			nominees: nominees,
		}
		mock := UserMockData{
			user:          user,
			product:       product,
			usersProducts: usersProducts,
			privileges:    privileges,
		}
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     input,
			Mock:     mock,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "valid_data_has_no_nominee"
		usersProducts.ProductMap[productID] = 0
		expected = UserExpectedData{
			userDeleted:          true,
			productDeleted:       true,
			usersProductsUpdated: false,
			err:                  nil,
		}
		input = UserInputData{
			userData: userData,
			nominees: nil,
		}
		mock = UserMockData{
			user:          user,
			product:       product,
			usersProducts: usersProducts,
			privileges:    privileges,
		}
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     input,
			Mock:     mock,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_user"
		usersProducts.ProductMap[productID] = 0
		expected = UserExpectedData{
			userDeleted:          false,
			productDeleted:       false,
			usersProductsUpdated: false,
			err:                  ErrUserNotFound,
		}
		input = UserInputData{
			userData: userData,
			nominees: nil,
		}
		mock = UserMockData{
			user:          user,
			product:       product,
			usersProducts: usersProducts,
			privileges:    privileges,
			err:           sql.ErrNoRows,
		}
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     input,
			Mock:     mock,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "has_no_products"
		usersProducts = &models.UserProductIDs{
			ProductMap: make(map[uuid.UUID]int),
		}
		expected = UserExpectedData{
			userDeleted:          true,
			productDeleted:       false,
			usersProductsUpdated: false,
			err:                  nil,
		}
		input = UserInputData{
			userData: userData,
			nominees: nil,
		}
		mock = UserMockData{
			user:          user,
			product:       product,
			usersProducts: usersProducts,
			privileges:    privileges,
			err:           nil,
		}
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     input,
			Mock:     mock,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	}

	return dataSet, nil
}

func TestCreateUser(t *testing.T) {
	// Create test data
	dataSet, err := createUserTestData(CreateUserTest)
	if err != nil {
		t.Errorf("Failed to create test data %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(UserExpectedData)
			inputData := testCase.Data.(UserInputData)
			mockData := testCase.Mock.(UserMockData)

			dbController.DBFunctions = &DBFunctionMock{
				user:         mockData.user,
				userAdded:    false,
				productAdded: false,
			}

			output, err := dbController.CreateUser(
				inputData.userData.Name,
				inputData.userData.Email,
				inputData.password,
				func(*uuid.UUID) (string, error) {
					return "testPath", nil
				}, func([]byte) ([]byte, error) {
					return []byte{}, nil
				})

			if diff := pretty.Diff(output, expectedData.userData); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData.userData, diff)
				return
			}

			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
				return
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	// Create test data
	dataSet, err := createUserTestData(DeleteUserTest)
	if err != nil {
		t.Errorf("Failed to create test data %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(UserExpectedData)
			inputData := testCase.Data.(UserInputData)
			mockData := testCase.Mock.(UserMockData)

			dbController.DBFunctions = &DBFunctionMock{
				user:                 mockData.user,
				product:              mockData.product,
				userProducts:         mockData.usersProducts,
				err:                  mockData.err,
				userDeleted:          false,
				productDeleted:       false,
				usersProductsUpdated: false,
				privileges:           mockData.privileges,
			}

			err := dbController.DeleteUser(&inputData.userData.ID, inputData.nominees)
			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
				return
			}

			if diff := pretty.Diff(dbController.DBFunctions.(*DBFunctionMock).userDeleted, expectedData.userDeleted); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, dbController.DBFunctions.(*DBFunctionMock).userDeleted, expectedData.userDeleted)
				return
			}

			if diff := pretty.Diff(dbController.DBFunctions.(*DBFunctionMock).productDeleted, expectedData.productDeleted); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, dbController.DBFunctions.(*DBFunctionMock).productDeleted, expectedData.productDeleted)
				return
			}

			if diff := pretty.Diff(dbController.DBFunctions.(*DBFunctionMock).usersProductsUpdated, expectedData.usersProductsUpdated); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, dbController.DBFunctions.(*DBFunctionMock).usersProductsUpdated, expectedData.usersProductsUpdated)
				return
			}
		})
	}
}
