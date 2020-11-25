package controllers

import (
	"fmt"
	"testing"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func createTestProductData() (*models.Product, error) {
	assetID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	details := make(models.Details)

	productID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	models.Interface = ModelInterfaceMock{
		assetID:   assetID,
		productID: productID,
	}

	expected := models.Product{
		Name:     "testProduct",
		Public:   true,
		ID:       productID,
		Details:  details,
		AssetsID: assetID,
	}

	mysqldb.FunctionInterface = DBFunctionInterfaceMock{}
	mysqldb.DBConnector = DBConnectorMock{}
	return &expected, nil
}

func createTestUsersData() (models.ProductUsers, models.Privileges) {
	privileges := make(models.Privileges, 2)
	privileges[0].ID = 0
	privileges[0].Name = "Owner"
	privileges[0].Description = "description0"
	privileges[1].ID = 1
	privileges[1].Name = "User"
	privileges[1].Description = "description1"
	mysqldb.DBConnector = DBConnectorMock{}

	users := make(models.ProductUsers)

	return users, privileges
}

func TestCreateProduct_NoExistingProduct(t *testing.T) {
	// Create test data
	expected, err := createTestProductData()
	if err != nil {
		t.Errorf("Failed to create test data: %s", err)
		return
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}

	users, privileges := createTestUsersData()
	users[userID] = 0

	mysqldb.FunctionInterface = DBFunctionInterfaceMock{
		privileges: privileges,
	}
	users[userID] = 0

	// Execute test
	product, err := CreateProduct(
		expected.Name,
		expected.Public,
		users,
		func(*uuid.UUID) string {
			return "testPath"
		})
	if err != nil {
		t.Errorf("Failed to create user: %s", err)
		return
	}

	if !cmp.Equal(*product, *expected) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", *product, *expected)
		return
	}
}

func TestCreateProduct_IncorrectUsersData(t *testing.T) {
	// Create test data
	expected, err := createTestProductData()
	if err != nil {
		t.Errorf("Failed to create test data: %s", err)
		return
	}
	productUsers := make(models.ProductUsers)

	// Execute test
	_, err = CreateProduct(
		expected.Name,
		expected.Public,
		productUsers,
		func(*uuid.UUID) string {
			return "testPath"
		})
	if err == nil || (err != nil && err != ErrEmptyUsersList) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrEmptyUsersList)
		return
	}
}

func TestCreateProduct_ProductAlreadyExists(t *testing.T) {
	// Create test data
	expected, err := createTestProductData()
	if err != nil {
		t.Errorf("Failed to create test data: %s", err)
		return
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}

	users, privileges := createTestUsersData()
	users[userID] = 0

	mysqldb.FunctionInterface = DBFunctionInterfaceMock{
		product:    expected,
		privileges: privileges,
	}

	expectedError := fmt.Errorf(ErrProductExistsString, expected.Name)
	// Execute test
	_, err = CreateProduct(
		expected.Name,
		expected.Public,
		users,
		func(*uuid.UUID) string {
			return "testPath"
		})
	if err == nil || (err != nil && !cmp.Equal(err.Error(), expectedError.Error())) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, expectedError)
		return
	}
}

func TestValidateUsers_ValidData(t *testing.T) {
	// Create test data
	users, privileges := createTestUsersData()
	mysqldb.FunctionInterface = DBFunctionInterfaceMock{
		privileges: privileges,
	}
	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 0

	// Execute test
	if err := validateUsers(users); err != nil {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, nil)
		return
	}
}

func TestValidateUsers_EmptyUsersList(t *testing.T) {
	// Create test data
	users, privileges := createTestUsersData()
	mysqldb.FunctionInterface = DBFunctionInterfaceMock{
		privileges: privileges,
	}

	// Execute test
	if err := validateUsers(users); err == nil || (err != nil && err != ErrEmptyUsersList) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrEmptyUsersList)
		return
	}
}

func TestValidateUsers_NilUserList(t *testing.T) {
	// Create test data
	var users models.ProductUsers

	// Execute test
	if err := validateUsers(users); err == nil || (err != nil && err != ErrEmptyUsersList) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrEmptyUsersList)
		return
	}
}

func TestValidateUsers_NoOwner(t *testing.T) {
	// Create test data
	users, privileges := createTestUsersData()
	mysqldb.FunctionInterface = DBFunctionInterfaceMock{
		privileges: privileges,
	}
	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 1

	userID, err = uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 1

	// Execute test
	if err := validateUsers(users); err == nil || (err != nil && err != ErrInvalidOwnerCount) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrInvalidOwnerCount)
		return
	}
}

func TestValidateUsers_MultipleOwners(t *testing.T) {
	// Create test data
	users, privileges := createTestUsersData()
	mysqldb.FunctionInterface = DBFunctionInterfaceMock{
		privileges: privileges,
	}
	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 0

	userID, err = uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 0

	// Execute test
	if err := validateUsers(users); err == nil || (err != nil && err != ErrInvalidOwnerCount) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrInvalidOwnerCount)
		return
	}
}

func TestValidateUsers_InvalidPrivilege(t *testing.T) {
	// Create test data
	users, privileges := createTestUsersData()
	mysqldb.FunctionInterface = DBFunctionInterfaceMock{
		privileges: privileges,
	}
	userIDInvalid, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userIDInvalid] = 2

	userID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create user uuid: %s", err)
		return
	}
	users[userID] = 0

	// Execute test
	expected := fmt.Errorf(ErrUnknownPrivilegeString, 2, userIDInvalid.String())
	if err := validateUsers(users); err == nil || (err != nil && !cmp.Equal(err.Error(), expected.Error())) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, expected)
		return
	}
}
