package mysqldb

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

func createTestSettings() (*models.UserSetting, error) {
	settings := make(models.Settings)

	settingID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	setting := models.UserSetting{
		ID:       settingID,
		Settings: settings,
	}
	return &setting, nil
}

func TestAddSettings_ValidUserSetting(t *testing.T) {
	// Create test data
	settings, err := createTestSettings()
	if err != nil {
		t.Errorf("Failed to generate settings test data %s", err)
		return
	}

	binary, err := json.Marshal(settings.Settings)
	if err != nil {
		t.Errorf("Failed to marshal settings references map %s", err)
		return
	}

	// Prepare mock
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate user test data %s", err)
		return
	}

	mock.ExpectBegin()
	mock.ExpectExec(AddUserSettingsQuery).WithArgs(settings.ID, binary).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	DBConnector = DBConnectorMock{
		DB:   db,
		Mock: mock,
	}
	defer db.Close()

	FunctionInterface = MYSQLFunctionInterface{}

	// Run test
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to setup DB transaction %s", err)
		return
	}

	err = FunctionInterface.AddSettings(settings, tx)
	if err != nil {
		t.Errorf("Failed to add settings %s", err)
		return
	}
}
