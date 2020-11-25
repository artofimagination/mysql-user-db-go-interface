package controllers

import (
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func createTestUserData() (*models.User, error) {
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

	expected := models.User{
		Name:       "testName",
		Email:      "testEmail",
		Password:   []byte{},
		ID:         userID,
		SettingsID: settingsID,
		AssetsID:   assetID,
	}

	mysqldb.FunctionInterface = DBFunctionInterfaceMock{}
	mysqldb.DBConnector = DBConnectorMock{}
	return &expected, nil
}

func TestCreateUser_NoExistingUser(t *testing.T) {
	// Create test data
	expected, err := createTestUserData()
	if err != nil {
		t.Errorf("Failed to create test data %s", err)
		return
	}

	// Execute test
	user, err := CreateUser(
		expected.Name,
		expected.Email,
		"",
		func(*uuid.UUID) string {
			return "testPath"
		}, func(string) ([]byte, error) {
			return []byte{}, nil
		})
	if err != nil {
		t.Errorf("Failed to create user %s", err)
		return
	}

	if !cmp.Equal(*user, *expected) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", *user, *expected)
		return
	}
}
