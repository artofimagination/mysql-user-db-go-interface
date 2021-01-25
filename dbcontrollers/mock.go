package dbcontrollers

import (
	"database/sql"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

// ModelMock overwrites the models package implementations with mock code.
type ModelMock struct {
	assetID    uuid.UUID
	settingsID uuid.UUID
	userID     uuid.UUID
	productID  uuid.UUID
	projectID  uuid.UUID
	project    *models.Project
	asset      *models.Asset

	err error
}

func (i *ModelMock) NewProduct(name string, detailsID *uuid.UUID, assetsID *uuid.UUID) (*models.Product, error) {
	p := &models.Product{
		Name:      name,
		ID:        i.productID,
		DetailsID: *detailsID,
		AssetsID:  *assetsID,
	}
	return p, i.err
}

func (i *ModelMock) NewAsset(references models.DataMap, generatePath func(assetID *uuid.UUID) (string, error)) (*models.Asset, error) {
	return i.asset, i.err
}

func (i *ModelMock) GetFilePath(asset *models.Asset, typeString string, defaultPath string) string {
	return ""
}
func (i *ModelMock) SetFilePath(asset *models.Asset, typeString string, extension string) error {
	return i.err
}
func (i *ModelMock) GetField(asset *models.Asset, typeString string, defaultURL string) string {
	return ""
}
func (i *ModelMock) SetField(asset *models.Asset, typeString string, field interface{}) {}
func (i *ModelMock) ClearAsset(asset *models.Asset, typeString string) error {
	return i.err
}

func (i *ModelMock) NewProject(productID *uuid.UUID, detailsID *uuid.UUID, assetsID *uuid.UUID) (*models.Project, error) {
	return i.project, i.err
}

func (i *ModelMock) NewUser(
	name string,
	email string,
	password []byte,
	settingsID uuid.UUID,
	assetsID uuid.UUID) (*models.User, error) {
	u := &models.User{
		ID:         i.userID,
		Name:       name,
		Email:      email,
		Password:   password,
		SettingsID: settingsID,
		AssetsID:   assetsID,
	}

	return u, i.err
}

// DBFunctionMock overwrites the mysqldb package function implementations with mock code.
type DBFunctionMock struct {
	user                 *models.User
	userDeleted          bool
	userAdded            bool
	product              *models.Product
	project              *models.Project
	productAdded         bool
	projectAdded         bool
	productDeleted       bool
	usersProductsUpdated bool
	privileges           models.Privileges
	userProducts         *models.UserProductIDs
	userProjects         *models.UserProjectIDs
	productUsers         *models.ProductUserIDs
	err                  error
}

func (i *DBFunctionMock) GetPrivileges() (models.Privileges, error) {
	return i.privileges, i.err
}

func (i *DBFunctionMock) GetPrivilege(name string) (*models.Privilege, error) {
	return i.privileges[0], i.err
}

func (i *DBFunctionMock) GetUser(queryType int, keyValue interface{}, tx *sql.Tx) (*models.User, error) {
	return i.user, i.err
}

func (i *DBFunctionMock) AddUser(user *models.User, tx *sql.Tx) error {
	i.userAdded = true
	return i.err
}

func (i *DBFunctionMock) GetUsersByIDs(IDs []uuid.UUID, tx *sql.Tx) ([]models.User, error) {
	return nil, i.err
}

func (i *DBFunctionMock) GetAssets(assetType string, IDs []uuid.UUID, tx *sql.Tx) ([]models.Asset, error) {
	return nil, i.err
}

func (i *DBFunctionMock) GetProductUserIDs(productID *uuid.UUID, tx *sql.Tx) (*models.ProductUserIDs, error) {
	return i.productUsers, i.err
}

func (i *DBFunctionMock) AddAsset(assetType string, asset *models.Asset, tx *sql.Tx) error {
	return i.err
}

func (i *DBFunctionMock) DeleteAsset(assetType string, assetID *uuid.UUID, tx *sql.Tx) error {
	return i.err
}

func (i *DBFunctionMock) GetAsset(assetType string, assetID *uuid.UUID) (*models.Asset, error) {
	return nil, i.err
}

func (i *DBFunctionMock) UpdateAsset(assetType string, asset *models.Asset) error {
	return i.err
}

func (i *DBFunctionMock) AddProductUsers(productID *uuid.UUID, productUsers *models.ProductUserIDs, tx *sql.Tx) error {
	return i.err
}

func (i *DBFunctionMock) AddProduct(product *models.Product, tx *sql.Tx) error {
	i.productAdded = true
	return i.err
}

func (i *DBFunctionMock) GetProductsByIDs(IDs []uuid.UUID, tx *sql.Tx) ([]models.Product, error) {
	return nil, i.err
}

func (i *DBFunctionMock) GetProductByName(name string, tx *sql.Tx) (*models.Product, error) {
	return i.product, i.err
}

func (i *DBFunctionMock) GetUserProductIDs(userID *uuid.UUID, tx *sql.Tx) (*models.UserProductIDs, error) {
	return i.userProducts, i.err
}

func (i *DBFunctionMock) DeleteProduct(productID *uuid.UUID, tx *sql.Tx) error {
	i.productDeleted = true
	return i.err
}

func (i *DBFunctionMock) DeleteUser(userID *uuid.UUID, tx *sql.Tx) error {
	i.userDeleted = true
	return i.err
}

func (i *DBFunctionMock) DeleteProductUsersByProductID(productID *uuid.UUID, tx *sql.Tx) error {
	return i.err
}

func (i *DBFunctionMock) DeleteProductUser(productID *uuid.UUID, userID *uuid.UUID, tx *sql.Tx) error {
	return i.err
}

func (i *DBFunctionMock) UpdateUsersProducts(userID *uuid.UUID, productID *uuid.UUID, privilege int, tx *sql.Tx) error {
	i.usersProductsUpdated = true
	return i.err
}

func (i *DBFunctionMock) GetProductByID(ID *uuid.UUID, tx *sql.Tx) (*models.Product, error) {
	return i.product, i.err
}

func (i DBFunctionMock) AddProject(project *models.Project, tx *sql.Tx) error {
	return i.err
}

func (i *DBFunctionMock) AddProjectUsers(projectID *uuid.UUID, projectUsers *models.ProjectUserIDs, tx *sql.Tx) error {
	return i.err
}

func (i *DBFunctionMock) GetProjectByID(ID *uuid.UUID, tx *sql.Tx) (*models.Project, error) {
	return i.project, i.err
}

func (i *DBFunctionMock) UpdateUsersProjects(userID *uuid.UUID, projectID *uuid.UUID, privilege int, tx *sql.Tx) error {
	return i.err
}

func (i *DBFunctionMock) DeleteProjectUsersByProjectID(projectID *uuid.UUID, tx *sql.Tx) error {
	return i.err
}

func (i *DBFunctionMock) DeleteProject(projectID *uuid.UUID, tx *sql.Tx) error {
	return i.err
}

func (i *DBFunctionMock) GetUserProjectIDs(userID *uuid.UUID, tx *sql.Tx) (*models.UserProjectIDs, error) {
	return i.userProjects, i.err
}

func (i *DBFunctionMock) GetProjectsByIDs(IDs []uuid.UUID, tx *sql.Tx) ([]models.Project, error) {
	return nil, i.err
}

func (i *DBFunctionMock) DeleteProjectsByProductID(productID *uuid.UUID, tx *sql.Tx) error {
	return i.err
}

// DBConnectorMock overwrites the mysqldb package implementations for DB connectionwith mock code.
type DBConnectorMock struct {
	err error
}

func (i *DBConnectorMock) BootstrapSystem() error {
	return i.err
}

func (i *DBConnectorMock) ConnectSystem() (*sql.Tx, error) {
	return nil, i.err
}

func (i *DBConnectorMock) Commit(tx *sql.Tx) error {
	return i.err
}

func (i DBConnectorMock) Rollback(tx *sql.Tx) error {
	return i.err
}

var dbController *MYSQLController
