package mysqldb

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
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

func (DBConnectorMock) Rollback(tx *sql.Tx) error {
	return nil
}

func createTestProductUsersData() (*models.ProductUserIDs, error) {
	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	owners := models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	owners.UserMap[userID] = 1
	owners.UserIDArray = append(owners.UserIDArray, userID)
	return &owners, nil
}

func createTestProjectUsersData() (*models.ProjectUserIDs, error) {
	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	owners := models.ProjectUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	owners.UserMap[userID] = 1
	owners.UserIDArray = append(owners.UserIDArray, userID)
	return &owners, nil
}
