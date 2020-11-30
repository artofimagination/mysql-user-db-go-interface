package mysqldb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	// Need to register mysql drivers with database/sql
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
)

// Common interface for DB connection. Needed in order to allow mock and custom DB interface implementation.
type DBConnectorCommon interface {
	BootstrapSystem() error
	ConnectSystem() (*sql.Tx, error)
}

// MYSQL Interface implementation
type MYSQLConnector struct {
}

// Data handling common function interface. Needed in order to allow mock and custom functionality implementations.
type FunctionCommonInterface interface {
	GetUser(keyName string, keyValue interface{}, tx *sql.Tx) (*models.User, error)
	AddUser(user *models.User, tx *sql.Tx) error
	DeleteUser(userID *uuid.UUID, tx *sql.Tx) error

	AddAsset(assetType string, asset *models.Asset, tx *sql.Tx) error
	DeleteAsset(assetType string, assetID *uuid.UUID, tx *sql.Tx) error

	GetProductsByUserID(userID uuid.UUID) ([]models.Product, error)
	UpdateUsersProducts(userID *uuid.UUID, productID *uuid.UUID, privilege int, tx *sql.Tx) error
	AddProductUsers(productID *uuid.UUID, productUsers models.ProductUsers, tx *sql.Tx) error
	DeleteProductUsersByProductID(productID *uuid.UUID, tx *sql.Tx) error
	GetUserProductIDs(userID uuid.UUID, tx *sql.Tx) (*models.UserProducts, error)

	GetProductByID(ID uuid.UUID) (*models.Product, error)
	GetProductByName(name string, tx *sql.Tx) (*models.Product, error)
	AddProduct(product *models.Product, tx *sql.Tx) error
	DeleteProduct(productID *uuid.UUID, tx *sql.Tx) error

	GetPrivileges() (models.Privileges, error)
}

// MYSQL Interface implementation
type MYSQLFunctions struct {
}

var DBConnection = ""
var DBConnector DBConnectorCommon
var Functions FunctionCommonInterface
var MigrationDirectory = ""

func (MYSQLConnector) BootstrapSystem() error {
	fmt.Printf("Executing MYSQL migration\n")
	migrations := &migrate.FileMigrationSource{
		Dir: MigrationDirectory,
	}
	fmt.Printf("Getting migration files\n")

	db, err := sql.Open("mysql", DBConnection)
	if err != nil {
		return err
	}
	fmt.Printf("DB connection open\n")

	n := 0
	for retryCount := 10; retryCount > 0; retryCount-- {
		n, err = migrate.Exec(db, "mysql", migrations, migrate.Up)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
		log.Printf("Failed to execute migration %s. Retrying...\n", err.Error())
	}

	if err != nil {
		return errors.Wrap(errors.WithStack(err), "Migration failed after multiple retries.")
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}

func RollbackWithErrorStack(tx *sql.Tx, errorStack error) error {
	if err := tx.Rollback(); err != nil {
		errorString := fmt.Sprintf("%s\n%s\n", errorStack.Error(), err.Error())
		return errors.Wrap(errors.WithStack(errors.New(errorString)), "Failed to rollback changes")
	}
	return errorStack
}

func (MYSQLConnector) ConnectSystem() (*sql.Tx, error) {
	db, err := sql.Open("mysql", DBConnection)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}
