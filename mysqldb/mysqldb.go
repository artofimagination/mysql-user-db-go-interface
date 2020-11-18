package mysqldb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// Need to register mysql drivers with database/sql
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
)

type DBInterfaceCommon interface {
	BootstrapSystem() error
	ConnectSystem() (*sql.DB, error)
}

type MYSQLInterface struct {
}

var DBConnection = ""
var DBInterface DBInterfaceCommon
var MigrationDirectory = ""

func (MYSQLInterface) BootstrapSystem() error {
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

func (MYSQLInterface) ConnectSystem() (*sql.DB, error) {
	db, err := sql.Open("mysql", DBConnection)

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	return db, nil
}
