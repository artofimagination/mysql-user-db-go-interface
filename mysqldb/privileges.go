package mysqldb

import (
	"database/sql"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
)

var GetPrivilegesQuery = "SELECT id, name, description from privileges"

func (MYSQLFunctions) GetPrivileges() (models.Privileges, error) {
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	rows, err := tx.Query(GetPrivilegesQuery)
	switch {
	case err == sql.ErrNoRows:
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	defer rows.Close()

	privileges := make(models.Privileges, 0)
	for rows.Next() {
		privilege := models.Privilege{}
		err := rows.Scan(&privilege.ID, &privilege.Name, &privilege.Description)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		privileges = append(privileges, privilege)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return privileges, tx.Commit()
}

var GetPrivilegeQuery = "SELECT id, name, description from privileges where name = ?"

func (MYSQLFunctions) GetPrivilege(name string) (*models.Privilege, error) {
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	rows := tx.QueryRow(GetPrivilegeQuery, name)
	switch {
	case err == sql.ErrNoRows:
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	privilege := models.Privilege{}
	err = rows.Scan(&privilege.ID, &privilege.Name, &privilege.Description)
	switch {
	case err == sql.ErrNoRows:
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	return &privilege, tx.Commit()
}
