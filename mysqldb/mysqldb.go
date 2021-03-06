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
type ConnectorCommon interface {
	BootstrapSystem() error
	ConnectSystem() (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
}

// MYSQL database connector implementation
type MYSQLConnector struct {
	DBConnection       string
	MigrationDirectory string
}

// Data handling common function interface. Needed in order to allow mock and custom functionality implementations.
type FunctionsCommon interface {
	GetUser(queryType int, keyValue interface{}, tx *sql.Tx) (*models.User, error)
	AddUser(user *models.User, tx *sql.Tx) error
	DeleteUser(userID *uuid.UUID, tx *sql.Tx) error
	GetProductUserIDs(productID *uuid.UUID, tx *sql.Tx) (*models.ProductUserIDs, error)
	GetUsersByIDs(IDs []uuid.UUID, tx *sql.Tx) ([]models.User, error)

	AddAsset(assetType string, asset *models.Asset, tx *sql.Tx) error
	DeleteAsset(assetType string, assetID *uuid.UUID, tx *sql.Tx) error
	GetAssets(assetType string, IDs []uuid.UUID, tx *sql.Tx) ([]models.Asset, error)
	GetAsset(assetType string, assetID *uuid.UUID) (*models.Asset, error)
	UpdateAsset(assetType string, asset *models.Asset) error

	UpdateUsersProducts(userID *uuid.UUID, productID *uuid.UUID, privilege int, tx *sql.Tx) error
	AddProductUsers(productID *uuid.UUID, productUsers *models.ProductUserIDs, tx *sql.Tx) error
	DeleteProductUsersByProductID(productID *uuid.UUID, tx *sql.Tx) error
	DeleteProductUser(productID *uuid.UUID, userID *uuid.UUID, tx *sql.Tx) error
	GetUserProductIDs(userID *uuid.UUID, tx *sql.Tx) (*models.UserProductIDs, error)

	GetProductByID(ID *uuid.UUID, tx *sql.Tx) (*models.Product, error)
	GetProductByName(name string, tx *sql.Tx) (*models.Product, error)
	AddProduct(product *models.Product, tx *sql.Tx) error
	DeleteProduct(productID *uuid.UUID, tx *sql.Tx) error
	GetProductsByIDs(IDs []uuid.UUID, tx *sql.Tx) ([]models.Product, error)

	AddProject(project *models.Project, tx *sql.Tx) error
	AddProjectUsers(projectID *uuid.UUID, projectUsers *models.ProjectUserIDs, tx *sql.Tx) error
	GetProjectByID(ID *uuid.UUID, tx *sql.Tx) (*models.Project, error)
	DeleteProjectUsersByProjectID(projectID *uuid.UUID, tx *sql.Tx) error
	DeleteProject(projectID *uuid.UUID, tx *sql.Tx) error
	DeleteProjectsByProductID(productID *uuid.UUID, tx *sql.Tx) error
	GetProjectsByIDs(IDs []uuid.UUID, tx *sql.Tx) ([]models.Project, error)
	GetProductProjects(productID *uuid.UUID, tx *sql.Tx) ([]models.Project, error)

	GetUserProjectIDs(userID *uuid.UUID, tx *sql.Tx) (*models.UserProjectIDs, error)
	UpdateUsersProjects(userID *uuid.UUID, projectID *uuid.UUID, privilege int, tx *sql.Tx) error
	AddProjectViewer(projectViewer *models.ProjectViewer, tx *sql.Tx) error
	DeleteProjectViewerByUserID(userID *uuid.UUID, tx *sql.Tx) error
	GetProjectViewersByUserID(userID *uuid.UUID, tx *sql.Tx) ([]models.ProjectViewer, error)
	GetProjectViewersByViewerID(viewerID *uuid.UUID, tx *sql.Tx) ([]models.ProjectViewer, error)
	DeleteProjectViewerByViewerID(viewerID *uuid.UUID, tx *sql.Tx) error
	DeleteProjectViewerByProjectID(projectID *uuid.UUID, tx *sql.Tx) error

	DeleteViewerByOwnerID(userID *uuid.UUID, tx *sql.Tx) error

	GetPrivileges() (models.Privileges, error)
	GetPrivilege(name string) (*models.Privilege, error)
}

// MYSQLFunctions represents the implementation of MYSQL data manipulation functions.
type MYSQLFunctions struct {
	DBConnector ConnectorCommon
	UUIDImpl    models.UUIDCommon
}

func (*MYSQLConnector) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (*MYSQLConnector) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (c *MYSQLConnector) BootstrapSystem() error {
	fmt.Printf("Executing MYSQL migration\n")
	migrations := &migrate.FileMigrationSource{
		Dir: c.MigrationDirectory,
	}
	fmt.Printf("Getting migration files\n")

	db, err := sql.Open("mysql", c.DBConnection)
	if err != nil {
		return err
	}
	fmt.Printf("DB connection open\n")

	n := 0
	for retryCount := 20; retryCount > 0; retryCount-- {
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

func (c *MYSQLConnector) ConnectSystem() (*sql.Tx, error) {
	db, err := sql.Open("mysql", c.DBConnection)
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
