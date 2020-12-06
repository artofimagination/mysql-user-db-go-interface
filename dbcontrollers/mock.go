package dbcontrollers

import (
	"database/sql"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

// ModelInterfaceMock overwrites the models package implementations with mock code.
type ModelInterfaceMock struct {
	assetID    uuid.UUID
	settingsID uuid.UUID
	userID     uuid.UUID
	productID  uuid.UUID

	err error
}

func (i ModelInterfaceMock) NewProduct(name string, public bool, detailsID *uuid.UUID, assetsID *uuid.UUID) (*models.Product, error) {
	p := models.Product{
		Name:      name,
		Public:    public,
		ID:        i.productID,
		DetailsID: *detailsID,
		AssetsID:  *assetsID,
	}
	return &p, i.err
}

func (i ModelInterfaceMock) NewAsset(references models.DataMap, generatePath func(assetID *uuid.UUID) string) (*models.Asset, error) {
	var a models.Asset
	a.ID = i.assetID
	return &a, i.err
}

func (i ModelInterfaceMock) NewUser(
	name string,
	email string,
	password []byte,
	settingsID uuid.UUID,
	assetsID uuid.UUID) (*models.User, error) {
	u := models.User{
		ID:         i.userID,
		Name:       name,
		Email:      email,
		Password:   password,
		SettingsID: settingsID,
		AssetsID:   assetsID,
	}

	return &u, i.err
}

// DBFunctionInterfaceMock overwrites the mysqldb package function implementations with mock code.
type DBFunctionInterfaceMock struct {
	user                 *models.User
	userDeleted          *bool
	userAdded            *bool
	product              *models.Product
	productAdded         *bool
	productDeleted       *bool
	usersProductsUpdated *bool
	privileges           models.Privileges
	userProducts         *models.UserProductIDs
	productUsers         *models.ProductUserIDs
	err                  error
}

func (i DBFunctionInterfaceMock) GetPrivileges() (models.Privileges, error) {
	return i.privileges, i.err
}

func (i DBFunctionInterfaceMock) GetPrivilege(name string) (*models.Privilege, error) {
	return &i.privileges[0], i.err
}

func (i DBFunctionInterfaceMock) GetUser(queryString string, keyValue interface{}, tx *sql.Tx) (*models.User, error) {
	return i.user, i.err
}

func (i DBFunctionInterfaceMock) AddUser(user *models.User, tx *sql.Tx) error {
	*i.userAdded = true
	return i.err
}

func (i DBFunctionInterfaceMock) GetProductUserIDs(productID *uuid.UUID, tx *sql.Tx) (*models.ProductUserIDs, error) {
	return i.productUsers, i.err
}

func (i DBFunctionInterfaceMock) AddAsset(assetType string, asset *models.Asset, tx *sql.Tx) error {
	return i.err
}

func (i DBFunctionInterfaceMock) DeleteAsset(assetType string, assetID *uuid.UUID, tx *sql.Tx) error {
	return i.err
}

func (i DBFunctionInterfaceMock) AddProductUsers(productID *uuid.UUID, productUsers *models.ProductUserIDs, tx *sql.Tx) error {
	return i.err
}

func (i DBFunctionInterfaceMock) AddProduct(product *models.Product, tx *sql.Tx) error {
	*i.productAdded = true
	return i.err
}

func (i DBFunctionInterfaceMock) GetProductByName(name string, tx *sql.Tx) (*models.Product, error) {
	return i.product, i.err
}

func (i DBFunctionInterfaceMock) GetUserProductIDs(userID *uuid.UUID, tx *sql.Tx) (*models.UserProductIDs, error) {
	return i.userProducts, i.err
}

func (i DBFunctionInterfaceMock) DeleteProduct(productID *uuid.UUID, tx *sql.Tx) error {
	*i.productDeleted = true
	return i.err
}

func (i DBFunctionInterfaceMock) DeleteUser(userID *uuid.UUID, tx *sql.Tx) error {
	*i.userDeleted = true
	return i.err
}

func (i DBFunctionInterfaceMock) DeleteProductUsersByProductID(productID *uuid.UUID, tx *sql.Tx) error {
	return i.err
}

func (i DBFunctionInterfaceMock) DeleteProductUser(productID *uuid.UUID, userID *uuid.UUID, tx *sql.Tx) error {
	return i.err
}

func (i DBFunctionInterfaceMock) UpdateUsersProducts(userID *uuid.UUID, productID *uuid.UUID, privilege int, tx *sql.Tx) error {
	*i.usersProductsUpdated = true
	return i.err
}

func (i DBFunctionInterfaceMock) GetProductByID(ID uuid.UUID, tx *sql.Tx) (*models.Product, error) {
	return i.product, i.err
}

// DBConnectorMock overwrites the mysqldb package implementations for DB connectionwith mock code.
type DBConnectorMock struct {
	err error
}

func (i DBConnectorMock) BootstrapSystem() error {
	return i.err
}

func (i DBConnectorMock) ConnectSystem() (*sql.Tx, error) {
	return nil, i.err
}

func (i DBConnectorMock) Commit(tx *sql.Tx) error {
	return i.err
}

func (i DBConnectorMock) Rollback(tx *sql.Tx) error {
	return i.err
}

var dbController *MYSQLController
