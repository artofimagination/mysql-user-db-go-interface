package controllers

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

	productDetails *models.ProductDetails

	err error
}

func (i ModelInterfaceMock) NewProduct(name string, public bool, details models.Details, assetsID *uuid.UUID) (*models.Product, error) {
	p := models.Product{
		Name:     name,
		Public:   public,
		ID:       i.productID,
		Details:  details,
		AssetsID: *assetsID,
	}
	return &p, i.err
}

func (i ModelInterfaceMock) NewAsset(references models.References, generatePath func(assetID *uuid.UUID) string) (*models.Asset, error) {
	var a models.Asset
	a.ID = i.assetID
	return &a, i.err
}

func (i ModelInterfaceMock) NewUserSettings(settings models.Settings) (*models.UserSettings, error) {
	var s models.UserSettings
	s.ID = i.settingsID
	return &s, i.err
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

func (i ModelInterfaceMock) NewProductDetails(details models.Details) (*models.ProductDetails, error) {
	return i.productDetails, i.err
}

// DBFunctionInterfaceMock overwrites the mysqldb package function implementations with mock code.
type DBFunctionInterfaceMock struct {
	user         *models.User
	product      *models.Product
	privileges   models.Privileges
	userProducts models.UserProducts
	err          error
}

func (i DBFunctionInterfaceMock) GetPrivileges() (models.Privileges, error) {
	return i.privileges, i.err
}

func (i DBFunctionInterfaceMock) AddSettings(settings *models.UserSettings, tx *sql.Tx) error {
	return i.err
}

func (i DBFunctionInterfaceMock) GetUserByEmail(email string, tx *sql.Tx) (*models.User, error) {
	return i.user, i.err
}

func (i DBFunctionInterfaceMock) AddUser(user *models.User, tx *sql.Tx) error {
	return i.err
}

func (i DBFunctionInterfaceMock) AddAsset(assetType string, asset *models.Asset, tx *sql.Tx) error {
	return i.err
}

func (i DBFunctionInterfaceMock) AddProductUsers(productID *uuid.UUID, productUsers models.ProductUsers, tx *sql.Tx) error {
	return i.err
}

func (i DBFunctionInterfaceMock) AddProduct(product *models.Product, tx *sql.Tx) error {
	return i.err
}

func (i DBFunctionInterfaceMock) GetProductByName(name string, tx *sql.Tx) (*models.Product, error) {
	return i.product, i.err
}

func (i DBFunctionInterfaceMock) GetUserProductIDs(userID uuid.UUID, tx *sql.Tx) (models.UserProducts, error) {
	return i.userProducts, i.err
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
