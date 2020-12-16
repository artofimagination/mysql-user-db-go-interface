package models

import (
	"errors"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/google/uuid"
	"github.com/kr/pretty"
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

func createProjectTestData(testID int) (*test.OrderedTests, error) {
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
		dataSet.TestDataSet[testCase] = test.Data{
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
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     inputData,
			Expected: expectedData,
			Mock:     mockData,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	}

	Interface = &RepoInterface{}

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

			UUIDImpl = &UUIDImplMock{
				uuidMock: mockData.projectID,
				err:      mockData.err,
			}

			output, err := Interface.NewProject(
				&inputData.project.ProductID,
				&inputData.project.DetailsID,
				&inputData.project.AssetsID,
			)
			if diff := pretty.Diff(output, expectedData.project); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData.project, diff)
				return
			}

			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
				return
			}
		})
	}
}
