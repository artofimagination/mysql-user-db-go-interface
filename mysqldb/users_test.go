package mysqldb

import (
	"database/sql"
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
		Password:   []byte{},
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

	// Prepare mock
	binaryUserID, err := json.Marshal(user.ID)
	if err != nil {
		t.Errorf("Failed to marshal user uuid %s", err)
		return
	}

	binarySettingsID, err := json.Marshal(user.SettingsID)
	if err != nil {
		t.Errorf("Failed to marshal settings uuid %s", err)
		return
	}

	binaryAssetsID, err := json.Marshal(user.AssetsID)
	if err != nil {
		t.Errorf("Failed to marshal asset test data %s", err)
		return
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate user test data %s", err)
		return
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "user_settings_id", "user_assets_id"}).
		AddRow(binaryUserID, user.Name, user.Email, user.Password, binarySettingsID, binaryAssetsID)
	mock.ExpectBegin()
	mock.ExpectQuery(GetUserByEmailQuery).WithArgs(user.Email).WillReturnRows(rows)
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

	data, err := FunctionInterface.GetUserByEmail(user.Email, tx)
	if err != nil {
		t.Errorf("Failed to get user %s", err)
		return
	}

	if !cmp.Equal(*data, *user) {
		t.Errorf("Test returned:\n %+v\nExpected:\n %+v", *data, *user)
		return
	}
}

func TestGetUserByEmail_InvalidEmail(t *testing.T) {
	// Create test data
	email := "test@test.com"

	// Prepare mock
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to create DB mock %s", err)
		return
	}

	mock.ExpectBegin()
	mock.ExpectQuery(GetUserByEmailQuery).WithArgs(email).WillReturnError(sql.ErrNoRows)
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

	_, err = FunctionInterface.GetUserByEmail(email, tx)
	if err == nil || (err != nil && err != ErrNoUserWithEmail) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrNoUserWithEmail)
		return
	}
}

func TestAddUser_Valid(t *testing.T) {
	// Create test data
	user, err := createUserTestData()
	if err != nil {
		t.Errorf("Failed to generate user test data %s", err)
		return
	}
	password := []byte{}

	// Prepare mock
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Failed to generate user test data %s", err)
		return
	}

	mock.ExpectBegin()
	mock.ExpectExec(InsertUserQuery).WithArgs(user.ID, user.Name, user.Email, password, user.SettingsID, user.AssetsID).WillReturnResult(sqlmock.NewResult(1, 1))
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

	err = FunctionInterface.AddUser(user, tx)
	if err != nil {
		t.Errorf("Failed to add user %s", err)
		return
	}
}
