package models

import (
	"errors"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/tests"
	"github.com/google/uuid"
)

const (
	NewProject = iota
)

type ProjectExpectedData struct {
	project *Project
	err     error
}

type ProjectMockData struct {
	projectID uuid.UUID
	err       error
}

type ProjectInputData struct {
	project *Project
}

func createProjectTestData(testID int) (*tests.OrderedTests, error) {
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

	projectID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	project := &Project{
		ID:        projectID,
		ProductID: productID,
		AssetsID:  assetsID,
		DetailsID: detailsID,
	}

	switch testID {
	case NewProject:
		testCase := "valid_product"
		inputData := ProjectInputData{
			project: project,
		}
		expectedData := ProjectExpectedData{
			project: project,
			err:     nil,
		}
		mockData := ProjectMockData{
			projectID: projectID,
			err:       nil,
		}
		dataSet.TestDataSet[testCase] = tests.Data{
			Data:     inputData,
			Expected: expectedData,
			Mock:     mockData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "failure_case"
		err := errors.New("Failed with error")
		expectedData = ProjectExpectedData{
			project: nil,
			err:     err,
		}
		mockData = ProjectMockData{
			projectID: projectID,
			err:       err,
		}
		dataSet.TestDataSet[testCase] = tests.Data{
			Data:     inputData,
			Expected: expectedData,
			Mock:     mockData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	}

	ModelFunctions = &RepoFunctions{}

	return dataSet, nil
}

func TestNewProject(t *testing.T) {
	// Create test data
	dataSet, err := createProjectTestData(NewProject)
	if err != nil {
		t.Errorf("Failed to create test data %s", err)
		return
	}

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			expectedData := dataSet.TestDataSet[testCaseString].Expected.(ProjectExpectedData)
			inputData := dataSet.TestDataSet[testCaseString].Data.(ProjectInputData)
			mockData := dataSet.TestDataSet[testCaseString].Mock.(ProjectMockData)

			ModelFunctions.UUIDImpl = &UUIDImplMock{
				uuidMock: mockData.projectID,
				err:      mockData.err,
			}

			output, err := ModelFunctions.NewProject(
				&inputData.project.ProductID,
				&inputData.project.DetailsID,
				&inputData.project.AssetsID,
			)
			tests.CheckResult(output, expectedData.project, err, expectedData.err, testCaseString, t)
		})
	}
}
