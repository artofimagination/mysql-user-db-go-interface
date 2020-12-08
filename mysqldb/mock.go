package mysqldb

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

type DBConnectorMock struct {
	Mock sqlmock.Sqlmock
	DB   *sql.DB
}

func (i *DBConnectorMock) ConnectSystem() (*sql.Tx, error) {
	tx, err := i.DB.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (*DBConnectorMock) BootstrapSystem() error {
	return nil
}

func (*DBConnectorMock) Commit(tx *sql.Tx) error {
	return nil
}
