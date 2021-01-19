package mysqldb

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/test"
	"github.com/kr/pretty"
)

const (
	GetPrivilegesTest = iota
	GetPrivilegeTest
)

type PrivilegeExpectedData struct {
	privilege  *models.Privilege
	privileges models.Privileges
	err        error
}

func createPrivilegesTestData(testID int) (*test.OrderedTests, error) {
	dataSet := &test.OrderedTests{
		OrderedList: make(test.OrderedTestList, 0),
		TestDataSet: make(test.DataSet),
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, err
	}

	privileges := make(models.Privileges, 2)
	privilege := &models.Privilege{
		ID:          0,
		Name:        "test0",
		Description: "description0",
	}
	privileges[0] = privilege
	privilege = &models.Privilege{
		ID:          1,
		Name:        "test1",
		Description: "description1",
	}
	privileges[1] = privilege

	switch testID {
	case GetPrivilegesTest:
		testCase := "valid_id"
		expected := PrivilegeExpectedData{
			privileges: privileges,
			err:        nil,
		}
		rows := sqlmock.NewRows([]string{"id", "name", "description"})
		for _, privilege := range privileges {
			rows.AddRow(privilege.ID, privilege.Name, privilege.Description)
		}

		mock.ExpectBegin()
		mock.ExpectQuery(GetPrivilegesQuery).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     nil,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_id"
		expected = PrivilegeExpectedData{
			err: sql.ErrNoRows,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(GetPrivilegesQuery).WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     nil,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)
	case GetPrivilegeTest:
		testCase := "valid_name"
		expected := PrivilegeExpectedData{
			privilege: privileges[0],
			err:       nil,
		}
		data := "Owner"
		rows := sqlmock.NewRows([]string{"id", "name", "description"})
		for _, privilege := range privileges {
			rows.AddRow(privilege.ID, privilege.Name, privilege.Description)
		}

		mock.ExpectBegin()
		mock.ExpectQuery(GetPrivilegeQuery).WithArgs(data).WillReturnRows(rows)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

		testCase = "invalid_name"
		expected = PrivilegeExpectedData{
			privileges: nil,
			err:        sql.ErrNoRows,
		}
		data = "TestName"

		mock.ExpectBegin()
		mock.ExpectQuery(GetPrivilegeQuery).WithArgs(data).WillReturnError(expected.err)
		mock.ExpectCommit()
		dataSet.TestDataSet[testCase] = test.Data{
			Data:     data,
			Expected: expected,
		}
		dataSet.OrderedList = append(dataSet.OrderedList, testCase)

	}

	DBFunctions = &MYSQLFunctions{
		DBConnector: &DBConnectorMock{
			DB:   db,
			Mock: mock,
		},
	}

	return dataSet, nil
}

func TestGetPrivileges(t *testing.T) {
	// Create test data
	dataSet, err := createPrivilegesTestData(GetPrivilegesTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			expectedData := dataSet.TestDataSet[testCaseString].Expected.(PrivilegeExpectedData)

			output, err := DBFunctions.GetPrivileges()
			if diff := pretty.Diff(output, expectedData.privileges); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData.privileges, diff)
				return
			}

			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err, diff)
				return
			}
		})
	}
}

func TestGetPrivilege(t *testing.T) {
	// Create test data
	dataSet, err := createPrivilegesTestData(GetPrivilegeTest)
	if err != nil {
		t.Errorf("Failed to generate test data: %s", err)
		return
	}

	defer DBFunctions.DBConnector.(*DBConnectorMock).DB.Close()

	// Run tests
	for _, testCaseString := range dataSet.OrderedList {
		testCaseString := testCaseString
		t.Run(testCaseString, func(t *testing.T) {
			expectedData := dataSet.TestDataSet[testCaseString].Expected.(PrivilegeExpectedData)
			inputData := dataSet.TestDataSet[testCaseString].Data.(string)

			output, err := DBFunctions.GetPrivilege(inputData)
			if diff := pretty.Diff(output, expectedData.privilege); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, output, expectedData.privilege)
				return
			}

			if diff := pretty.Diff(err, expectedData.err); len(diff) != 0 {
				t.Errorf(test.TestResultString, testCaseString, err, expectedData.err)
				return
			}
		})
	}
}
