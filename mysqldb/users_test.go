package mysqldb

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func createUserTestData() (*models.User, error) {
	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	settingsID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	assetsID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:         userID,
		Name:       "testName",
		Email:      "test@test.com",
		Password:   "testPass",
		SettingsID: settingsID,
		AssetsID:   assetsID,
	}
	return &user, nil
}

func TestGetUserByEmail_ValidEmail(t *testing.T) {
	// Create test data
	user, err := createUserTestData()
	if err != nil {
		t.Errorf("Failed to generate user test data %s", err)
		return
	}

	binaryUserID, err := json.Marshal(user.ID)
	if err != nil {
		t.Errorf("Failed to generate binary from user UUID %s", err)
		return
	}

	binarySettingsID, err := json.Marshal(user.SettingsID)
	if err != nil {
		t.Errorf("Failed to generate binary from settings UUID %s", err)
		return
	}

	binaryAssetsID, err := json.Marshal(user.AssetsID)
	if err != nil {
		t.Errorf("Failed to generate binary from settings UUID %s", err)
		return
	}

	// Create mock conditions
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate DB mock %s", err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "user_settings_id", "user_assets_id"}).
		AddRow(binaryUserID, user.Name, user.Email, user.Password, binarySettingsID, binaryAssetsID)
	mock.ExpectBegin()
	mock.ExpectQuery("select BIN_TO_UUID(id), name, email, password, BIN_TO_UUID(user_settings_id), BIN_TO_UUID(user_assets_id) from users where email = ?").WithArgs(user.Email).WillReturnRows(rows)
	mock.ExpectCommit()

	DBInterface = DBInterfaceMock{
		DB:   db,
		Mock: mock,
	}

	// Run test
	data, err := GetUserByEmail(user.Email)
	if err != nil {
		t.Errorf("Failed to get user %s", err)
		return
	}

	if !cmp.Equal(*data, *user) {
		t.Errorf("Test returned:\n %+v\nExpected:\n %+v", *data, *user)
	}
}
