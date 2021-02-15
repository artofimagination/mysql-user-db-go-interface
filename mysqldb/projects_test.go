package mysqldb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	addProjectTest = iota
	addProjectUsersTest
	deleteProjectUsersByProjectIDTest
	deleteProjectsByProductIDTest
	getProjectByIDTest
	getUserProjectIDsTest
	getProductProjectsTest
	deleteProjectTest
	updateUsersProjectsTest
)

func createTestProjectData() (*models.Project, error) {
	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	projectID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	detailsID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	project := &models.Project{
		ID:        projectID,
		ProductID: productID,
		DetailsID: detailsID,
		AssetsID:  assetID,
	}

	return project, nil
}

func createTestUserProjectsData(quantity int) (*models.UserProjectIDs, error) {
	userProjects := &models.UserProjectIDs{
		ProjectMap:     make(map[uuid.UUID]int),
		ProjectIDArray: make([]uuid.UUID, 0),
	}

	for ; quantity > 0; quantity-- {
		projectID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		userProjects.ProjectMap[projectID] = 1
		userProjects.ProjectIDArray = append(userProjects.ProjectIDArray, projectID)
	}
	return userProjects, nil
}

type ProjectExpectedData struct {
	project      *models.Project
	projects     []models.Project
	userProjects *models.UserProjectIDs
	err          error
}

type ProjectInputData struct {
	userID       *uuid.UUID
	productID    *uuid.UUID
	project      *models.Project
	projectUsers *models.ProjectUserIDs
	privilege    int
}

func createProjectsTestData(testID int) (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, err
	}
	project, err := createTestProjectData()
	if err != nil {
		return nil, err
	}
	project2, err := createTestProjectData()
	if err != nil {
		return nil, err
	}
	project2.ProductID = project.ProductID
	binaryProjectID, err := json.Marshal(project.ID)
	if err != nil {
		return nil, err
	}
	binaryDetailsID, err := json.Marshal(project.DetailsID)
	if err != nil {
		return nil, err
	}
	binaryAssetID, err := json.Marshal(project.AssetsID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	binaryProductID, err := json.Marshal(project.ProductID)
	if err != nil {
		return nil, err
	}
	userProjects, err := createTestUserProjectsData(2)
	if err != nil {
		return nil, err
	}
	projectUsers, err := createTestProjectUsersData()
	if err != nil {
		return nil, err
	}

	switch testID {
	case addProjectTest:

		testCase := "valid_product"
		inputData := ProjectInputData{
			project: project,
		}
		expectedData := ProjectExpectedData{
			err: nil,
		}
		mock.ExpectBegin()
		mock.ExpectExec(AddProjectQuery).WithArgs(project.ID, project.ProductID, project.DetailsID, project.AssetsID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		expected := errors.New("This is a failure test")
		expectedData = ProjectExpectedData{
			err: expected,
		}
		mock.ExpectBegin()
		mock.ExpectExec(AddProjectQuery).WithArgs(project.ID, project.ProductID, project.DetailsID, project.AssetsID).WillReturnError(expected)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case addProjectUsersTest:
		testCase := "valid_products"
		inputData := ProjectInputData{
			project:      project,
			projectUsers: projectUsers,
		}
		expectedData := ProjectExpectedData{
			err: nil,
		}
		mock.ExpectBegin()
		for _, userID := range projectUsers.UserIDArray {
			privilege := projectUsers.UserMap[userID]
			mock.ExpectExec(AddProjectUsersQuery).WithArgs(userID, project.ID, privilege).WillReturnResult(sqlmock.NewResult(1, 1))
		}
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_query"
		expected := errors.New("This is a failure test")
		expectedData = ProjectExpectedData{
			err: expected,
		}
		mock.ExpectBegin()
		for _, userID := range projectUsers.UserIDArray {
			privilege := projectUsers.UserMap[userID]
			mock.ExpectExec(AddProjectUsersQuery).WithArgs(userID, project.ID, privilege).WillReturnError(expected)
		}
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failed_to_add"
		expectedData = ProjectExpectedData{
			err: ErrNoProjectUserAdded,
		}
		mock.ExpectBegin()
		for _, userID := range projectUsers.UserIDArray {
			privilege := projectUsers.UserMap[userID]
			mock.ExpectExec(AddProjectUsersQuery).WithArgs(userID, project.ID, privilege).WillReturnResult(sqlmock.NewResult(1, 0))
		}
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case deleteProjectUsersByProjectIDTest:
		testCase := "valid_id"
		inputData := ProjectInputData{
			project: project,
		}
		expectedData := ProjectExpectedData{
			err: nil,
		}
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductUsersByProductIDQuery).WithArgs(project.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_id"
		expectedData = ProjectExpectedData{
			err: ErrNoUserWithProject,
		}
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProductUsersByProductIDQuery).WithArgs(project.ID).WillReturnError(expectedData.err)
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case getProjectByIDTest:

		testCase := "valid_id"
		inputData := ProjectInputData{
			project: project,
		}
		expectedData := ProjectExpectedData{
			project: project,
			err:     nil,
		}
		rows := sqlmock.NewRows([]string{"id", "product_id", "product_details_id", "product_assets_id"}).
			AddRow(binaryProjectID, binaryProductID, binaryDetailsID, binaryAssetID)
		mock.ExpectBegin()
		mock.ExpectQuery(GetProjectByIDQuery).WithArgs(project.ID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_id"
		expectedData = ProjectExpectedData{
			project: nil,
			err:     sql.ErrNoRows,
		}
		mock.ExpectBegin()
		mock.ExpectQuery(GetProjectByIDQuery).WithArgs(project.ID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case getUserProjectIDsTest:
		testCase := "valid_id"
		inputData := ProjectInputData{
			userID: &userID,
		}
		expectedData := ProjectExpectedData{
			userProjects: userProjects,
			err:          nil,
		}
		rows := sqlmock.NewRows([]string{"projects_id", "privilege"})
		for _, productID := range userProjects.ProjectIDArray {
			rows.AddRow(productID, userProjects.ProjectMap[productID])
		}
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProjectIDsQuery).WithArgs(userID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_products"
		expectedData = ProjectExpectedData{
			userProjects: nil,
			err:          sql.ErrNoRows,
		}
		mock.ExpectBegin()
		mock.ExpectQuery(GetUserProjectIDsQuery).WithArgs(userID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case getProductProjectsTest:
		testCase := "valid_id"
		inputData := ProjectInputData{
			productID: &project.ProductID,
		}
		projects := make([]models.Project, 2)
		projects = append(projects, *project)
		projects = append(projects, *project2)
		expectedData := ProjectExpectedData{
			projects: projects,
			err:      nil,
		}
		rows := sqlmock.NewRows([]string{"id", "products_id", "project_details_id", "project_assets_id"})
		for _, project := range projects {
			rows.AddRow(project.ID, project.ProductID, project.DetailsID, project.AssetsID)
		}
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductProjectsQuery).WithArgs(project.ProductID).WillReturnRows(rows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "missing_projects"
		expectedData = ProjectExpectedData{
			userProjects: nil,
			err:          sql.ErrNoRows,
		}
		mock.ExpectBegin()
		mock.ExpectQuery(GetProductProjectsQuery).WithArgs(project.ProductID).WillReturnError(sql.ErrNoRows)
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case deleteProjectTest:
		testCase := "valid_id"
		inputData := ProjectInputData{
			project: project,
		}
		expectedData := ProjectExpectedData{
			err: nil,
		}
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProjectQuery).WithArgs(project.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "no_product"
		expectedData = ProjectExpectedData{
			err: ErrNoProjectDeleted,
		}
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProjectQuery).WithArgs(project.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case deleteProjectsByProductIDTest:
		testCase := "valid_id"
		inputData := ProjectInputData{
			project: project,
		}
		expectedData := ProjectExpectedData{
			err: nil,
		}
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProjectsByProductIDQuery).WithArgs(project.ProductID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "no_project"
		expectedData = ProjectExpectedData{
			err: ErrNoProjectDeleted,
		}
		mock.ExpectBegin()
		mock.ExpectExec(DeleteProjectsByProductIDQuery).WithArgs(project.ProductID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	case updateUsersProjectsTest:
		testCase := "valid_id"
		inputData := ProjectInputData{
			project:   project,
			userID:    &userID,
			privilege: 1,
		}
		expectedData := ProjectExpectedData{
			err: nil,
		}
		mock.ExpectBegin()
		mock.ExpectExec(UpdateUsersProjectsQuery).WithArgs(inputData.privilege, userID, project.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "no_users_products"
		expectedData = ProjectExpectedData{
			err: ErrNoUsersProjectUpdate,
		}
		mock.ExpectBegin()
		mock.ExpectExec(UpdateUsersProjectsQuery).WithArgs(inputData.privilege, userID, project.ID).WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	default:
		return nil, fmt.Errorf("Unknown test %d", testID)
	}

	DBFunctions = &MYSQLFunctions{
		DBConnector: &DBConnectorMock{
			DB:   db,
			Mock: mock,
		},
	}

	return dataSet, nil
}

func TestAddProject(t *testing.T) {
	// Create test data
	dataSet, err := createProjectsTestData(addProjectTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			expectedData := dataSet.TestDataSet[testCaseString].Expected.(ProjectExpectedData)
			inputData := dataSet.TestDataSet[testCaseString].Data.(ProjectInputData)

			err = DBFunctions.AddProject(inputData.project, tx)
			test.CheckResult(nil, nil, err, expectedData.err, testCaseString, t)
		})
	}
}

func TestAddProjectUsers(t *testing.T) {
	// Create test data
	dataSet, err := createProjectsTestData(addProjectUsersTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			expectedData := dataSet.TestDataSet[testCaseString].Expected.(ProjectExpectedData)
			inputData := dataSet.TestDataSet[testCaseString].Data.(ProjectInputData)
			err = DBFunctions.AddProjectUsers(&inputData.project.ID, inputData.projectUsers, tx)
			test.CheckResult(nil, nil, err, expectedData.err, testCaseString, t)
		})
	}
}

func TestUpdateUsersProjects(t *testing.T) {
	// Create test data
	dataSet, err := createProjectsTestData(updateUsersProjectsTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			expectedData := dataSet.TestDataSet[testCaseString].Expected.(ProjectExpectedData)
			inputData := dataSet.TestDataSet[testCaseString].Data.(ProjectInputData)
			err = DBFunctions.UpdateUsersProjects(inputData.userID, &inputData.project.ID, inputData.privilege, tx)
			test.CheckResult(nil, nil, err, expectedData.err, testCaseString, t)
		})
	}
}

func TestDeleteProjectUsersByProjectID(t *testing.T) {
	// Create test data
	dataSet, err := createProjectsTestData(deleteProjectUsersByProjectIDTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			expectedData := dataSet.TestDataSet[testCaseString].Expected.(ProjectExpectedData)
			inputData := dataSet.TestDataSet[testCaseString].Data.(ProjectInputData)
			err = DBFunctions.DeleteProductUsersByProductID(&inputData.project.ID, tx)
			test.CheckResult(nil, nil, err, expectedData.err, testCaseString, t)
		})
	}
}

func TestGetProjectByID(t *testing.T) {
	// Create test data
	dataSet, err := createProjectsTestData(getProjectByIDTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}

			expectedData := dataSet.TestDataSet[testCaseString].Expected.(ProjectExpectedData)
			inputData := dataSet.TestDataSet[testCaseString].Data.(ProjectInputData)
			output, err := DBFunctions.GetProjectByID(&inputData.project.ID, tx)
			test.CheckResult(output, expectedData.project, err, expectedData.err, testCaseString, t)
		})
	}
}

func TestGetUserProjectIDs(t *testing.T) {
	// Create test data
	dataSet, err := createProjectsTestData(getUserProjectIDsTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			expectedData := dataSet.TestDataSet[testCaseString].Expected.(ProjectExpectedData)
			inputData := dataSet.TestDataSet[testCaseString].Data.(ProjectInputData)
			output, err := DBFunctions.GetUserProjectIDs(inputData.userID, tx)
			test.CheckResult(output, expectedData.userProjects, err, expectedData.err, testCaseString, t)
		})
	}
}

func TestGetProductProjects(t *testing.T) {
	// Create test data
	dataSet, err := createProjectsTestData(getProductProjectsTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
			if err != nil {
				t.Errorf("Failed to setup DB transaction %s", err)
				return
			}
			expectedData := dataSet.TestDataSet[testCaseString].Expected.(ProjectExpectedData)
			inputData := dataSet.TestDataSet[testCaseString].Data.(ProjectInputData)
			output, err := DBFunctions.GetProductProjects(inputData.productID, tx)
			test.CheckResult(output, expectedData.projects, err, expectedData.err, testCaseString, t)
		})
	}
}

func TestDeleteProject(t *testing.T) {
	// Create test data
	dataSet, err := createProjectsTestData(deleteProjectTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
		if err != nil {
			t.Errorf("Failed to setup DB transaction %s", err)
			return
		}

		testCaseString := testCaseString
		expectedData := dataSet.TestDataSet[testCaseString].Expected.(ProjectExpectedData)
		inputData := dataSet.TestDataSet[testCaseString].Data.(ProjectInputData)
		err = DBFunctions.DeleteProject(&inputData.project.ID, tx)
		test.CheckResult(nil, nil, err, expectedData.err, testCaseString, t)
	}
}

func TestDeleteProjectByProductID(t *testing.T) {
	// Create test data
	dataSet, err := createProjectsTestData(deleteProjectsByProductIDTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}
	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		tx, err := DBFunctions.DBConnector.(*DBConnectorMock).DB.Begin()
		if err != nil {
			t.Errorf("Failed to setup DB transaction %s", err)
			return
		}

		testCaseString := testCaseString
		expectedData := dataSet.TestDataSet[testCaseString].Expected.(ProjectExpectedData)
		inputData := dataSet.TestDataSet[testCaseString].Data.(ProjectInputData)
		err = DBFunctions.DeleteProjectsByProductID(&inputData.project.ProductID, tx)
		test.CheckResult(nil, nil, err, expectedData.err, testCaseString, t)
	}
}
