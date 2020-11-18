package mysqldb

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

type DBInterfaceMock struct {
	Mock sqlmock.Sqlmock
	DB   *sql.DB
}

func (i DBInterfaceMock) ConnectSystem() (*sql.Tx, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (DBInterfaceMock) BootstrapSystem() error {
	return nil
}
