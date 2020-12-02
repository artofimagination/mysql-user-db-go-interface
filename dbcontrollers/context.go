package dbcontrollers

import (
	"errors"
	"fmt"
	"os"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/google/uuid"
)

type DBControllerCommon interface {
	CreateProduct(name string, public bool, users models.ProductUsers, generateAssetPath func(assetID *uuid.UUID) string) (*models.Product, error)
	DeleteProduct(productID *uuid.UUID) error
	GetProduct(productID *uuid.UUID) (*models.ProductData, error)
	UpdateProductDetails(details *models.Asset) error
	UpdateProductAssets(assets *models.Asset) error

	CreateUser(
		name string,
		email string,
		passwd []byte,
		generateAssetPath func(assetID *uuid.UUID) string,
		encryptPassword func(password []byte) ([]byte, error)) (*models.User, error)
	DeleteUser(ID *uuid.UUID, nominatedOwners map[uuid.UUID]uuid.UUID) error
	GetUser(userID *uuid.UUID) (*models.UserData, error)
	UpdateUserSettings(settings *models.Asset) error
	UpdateUserAssets(assets *models.Asset) error
	Authenticate(email string, passwd []byte, authenticate func(string, []byte, *models.User) error) error
}

type MYSQLController struct{}

func NewDBController() (*MYSQLController, error) {
	address := os.Getenv("MYSQL_DB_ADDRESS")
	if address == "" {
		return nil, errors.New("MYSQL DB address not defined")
	}
	port := os.Getenv("MYSQL_DB_PORT")
	if address == "" {
		return nil, errors.New("MYSQL DB port not defined")
	}
	username := os.Getenv("MYSQL_DB_USER")
	if address == "" {
		return nil, errors.New("MYSQL DB username not defined")
	}
	pass := os.Getenv("MYSQL_DB_PASSWORD")
	if address == "" {
		return nil, errors.New("MYSQL DB password not defined")
	}
	dbName := os.Getenv("MYSQL_DB_NAME")
	if address == "" {
		return nil, errors.New("MYSQL DB name not defined")
	}

	mysqldb.MigrationDirectory = os.Getenv("MYSQL_DB_PASSWORD")
	if mysqldb.MigrationDirectory == "" {
		return nil, errors.New("MYSQL DB migration folder not defined")
	}
	mysqldb.DBConnection = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username,
		pass,
		address,
		port,
		dbName)
	mysqldb.Functions = mysqldb.MYSQLFunctions{}
	mysqldb.DBConnector = mysqldb.MYSQLConnector{}

	models.Interface = models.RepoInterface{}
	models.UUIDImpl = models.RepoUUIDInterface{}
	controller := MYSQLController{}
	return &controller, nil
}
