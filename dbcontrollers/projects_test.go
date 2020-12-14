package dbcontrollers

import (
	"fmt"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

func createTestProjectUsersData() (*models.ProjectUserIDs, models.Privileges) {
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
	mysqldb.DBConnector = &DBConnectorMock{}

	users := models.ProjectUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}

	return &users, privileges
}

type ProjectExpectedData struct {
	projectData *models.ProjectData
	err         error
}

type ProjectMockData struct {
	projectData *models.ProjectData
	project     *models.Project
	privileges  models.Privileges
	err         error
}

type ProjectInputData struct {
	projectData *models.ProjectData
	userID      uuid.UUID
}

func createProjectTestData() (*test.OrderedTests, error) {
	dataSet := test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	dbController = &MYSQLController{}

	_, privileges := createTestProjectUsersData()

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	projectID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	project := &models.Project{
		ID:        projectID,
		ProductID: productID,
		DetailsID: assetID,
		AssetsID:  assetID,
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	dataMap := make(models.DataMap)
	dataMap["name"] = "testProject"
	dataMap["visibility"] = models.Protected
	assets := &models.Asset{
		ID:      assetID,
		DataMap: dataMap,
	}

	projectData := &models.ProjectData{
		ID:        project.ID,
		ProductID: productID,
		Details:   assets,
		Assets:    assets,
	}

	testCase := "no_existing_project"
	expected := ProjectExpectedData{
		projectData: projectData,
		err:         nil,
	}
	input := ProjectInputData{
		projectData: projectData,
		userID:      userID,
	}
	mock := ProjectMockData{
		projectData: projectData,
		project:     project,
		privileges:  privileges,
	}

	dataSet.TestDataSet[testCase] = test.Data{
		Data:     input,
		Mock:     mock,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	testCase = "existing_project"
	expected = ProjectExpectedData{
		projectData: nil,
		err:         fmt.Errorf(ErrProjectExistsString, assets.DataMap["name"]),
	}

	mock = ProjectMockData{
		projectData: projectData,
		project:     project,
		privileges:  privileges,
		err:         fmt.Errorf(ErrProjectExistsString, assets.DataMap["name"]),
	}

	dataSet.TestDataSet[testCase] = test.Data{
		Data:     input,
		Mock:     mock,
		Expected: expected,
	}
	dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	mysqldb.Functions = &DBFunctionInterfaceMock{}
	mysqldb.DBConnector = &DBConnectorMock{}
	return &dataSet, nil
}

func TestCreateProject(t *testing.T) {
	// Create test data
	dataSet, err := createProjectTestData()
	if err != nil {
		t.Errorf("Failed to create test data: %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			testCase := dataSet.TestDataSet[testCaseString]
			expectedData := testCase.Expected.(ProjectExpectedData)
			inputData := testCase.Data.(ProjectInputData)
			mockData := testCase.Mock.(ProjectMockData)

			models.Interface = &ModelInterfaceMock{
				assetID:   mockData.project.AssetsID,
				projectID: mockData.project.ID,
				asset:     mockData.projectData.Assets,
				project:   mockData.project,
				err:       mockData.err,
			}
			mysqldb.Functions = &DBFunctionInterfaceMock{
				project:      mockData.project,
				privileges:   mockData.privileges,
				projectAdded: false,
			}

			output, err := dbController.CreateProject(
				inputData.projectData.Assets.DataMap["name"].(string),
				inputData.projectData.Assets.DataMap["visibility"].(int),
				&inputData.userID,
				func(*uuid.UUID) (string, error) {
					return "testPath", nil
				})

			if diff := pretty.Diff(output, expectedData.projectData); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData.projectData, diff)
				return
			}

			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
				return
			}
		})
	}
}
