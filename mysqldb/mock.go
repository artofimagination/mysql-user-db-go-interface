package mysqldb

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

type DBInterfaceMock struct {
	Mock sqlmock.Sqlmock
	DB   *sql.DB
}

func (i DBInterfaceMock) ConnectSystem() (*sql.DB, error) {
	return i.DB, nil
}

func (DBInterfaceMock) BootstrapSystem() error {
	return nil
}
