package mysqldb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

// Defines the possible user query key names
const (
	ByEmail = iota
	ByID
)

var ErrNoUserWithEmail = errors.New("There is no user associated with this email")

var ErrSQLDuplicateUserNameEntryString = "Error 1062: Duplicate entry '%s' for key 'users.name'"
var ErrSQLDuplicateEmailEntryString = "Error 1062: Duplicate entry '%s' for key 'users.email'"
var ErrDuplicateUserNameEntry = errors.New("User with this name already exists")
var ErrNoUserDeleted = errors.New("No user was deleted")

var GetUserByEmailQuery = "select BIN_TO_UUID(id), name, email, password, BIN_TO_UUID(user_settings_id), BIN_TO_UUID(user_assets_id) from users where email = ?"
var GetUserByIDQuery = "select BIN_TO_UUID(id), name, email, password, BIN_TO_UUID(user_settings_id), BIN_TO_UUID(user_assets_id) from users where id = UUID_TO_BIN(?)"

// GetUser returns the user defined by the key name and key value.
// Key name can be either id or email.
func (MYSQLFunctions) GetUser(queryType int, keyValue interface{}, tx *sql.Tx) (*models.User, error) {
	queryString := GetUserByIDQuery
	if queryType == ByEmail {
		queryString = GetUserByEmailQuery
		keyValue = strings.ReplaceAll(keyValue.(string), " ", "")
	}

	var user models.User
	query := tx.QueryRow(queryString, keyValue)
	password := ""
	err := query.Scan(&user.ID, &user.Name, &user.Email, &password, &user.SettingsID, &user.AssetsID)
	switch {
	case err == sql.ErrNoRows:
		return nil, err
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}
	user.Password = []byte(password)
	return &user, nil
}

var GetUsersByIDsQuery = "select BIN_TO_UUID(id), name, email, password, BIN_TO_UUID(user_settings_id), BIN_TO_UUID(user_assets_id) from users where id IN (UUID_TO_BIN(?)"

func (MYSQLFunctions) GetUsersByIDs(IDs []uuid.UUID, tx *sql.Tx) ([]models.User, error) {
	query := GetUsersByIDsQuery + strings.Repeat(",UUID_TO_BIN(?)", len(IDs)-1) + ")"
	interfaceList := make([]interface{}, len(IDs))
	for i := range IDs {
		interfaceList[i] = IDs[i]
	}
	rows, err := tx.Query(query, interfaceList...)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user := models.User{}
		password := ""
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &password, &user.SettingsID, &user.AssetsID)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		user.Password = []byte(password)
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users, nil
}

func UserExists(username string) (bool, error) {
	var user models.User
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return false, err
	}

	queryString := "SELECT name FROM users WHERE name = ?"
	queryUser := tx.QueryRow(queryString, username)
	err = queryUser.Scan(&user.Name)
	switch {
	case err == sql.ErrNoRows:
		return false, tx.Commit()
	case err != nil:
		return false, RollbackWithErrorStack(tx, err)
	default:
		return true, RollbackWithErrorStack(tx, err)
	}
}

func EmailExists(email string) (bool, error) {
	email = strings.ReplaceAll(email, " ", "")

	var user models.User
	queryString := "select email from users where email = ?"
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return false, err
	}

	queryEmail := tx.QueryRow(queryString, email)
	if err != nil {
		return false, RollbackWithErrorStack(tx, err)
	}

	err = queryEmail.Scan(&user.Email)
	switch {
	case err == sql.ErrNoRows:
		return false, tx.Commit()
	case err != nil:
		return false, RollbackWithErrorStack(tx, err)
	default:
		return true, RollbackWithErrorStack(tx, err)
	}
}

// CheckPassword compares the password entered by the user with the stored password.
func IsPasswordCorrect(password string, user *models.User) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

var InsertUserQuery = "INSERT INTO users (id, name, email, password, user_settings_id, user_assets_id) VALUES (UUID_TO_BIN(?), ?, ?, ?, UUID_TO_BIN(?), UUID_TO_BIN(?))"

// AddUser creates a new user entry in the DB.
// Whitespaces in the email are automatically deleted
// Email/Name are unique in DB. Duplicates will return error.
func (MYSQLFunctions) AddUser(user *models.User, tx *sql.Tx) error {
	_, err := tx.Exec(InsertUserQuery, user.ID, user.Name, user.Email, string(user.Password), user.SettingsID, user.AssetsID)
	errDuplicateName := fmt.Errorf(ErrSQLDuplicateUserNameEntryString, user.Name)
	errDuplicateEmail := fmt.Errorf(ErrSQLDuplicateEmailEntryString, user.Email)
	if err != nil {
		switch {
		case err.Error() == errDuplicateName.Error():
			if errRb := tx.Rollback(); errRb != nil {
				return err
			}
			return errDuplicateName
		case err.Error() == errDuplicateEmail.Error():
			if errRb := tx.Rollback(); errRb != nil {
				return err
			}
			return errDuplicateEmail
		case err != nil:
			return RollbackWithErrorStack(tx, err)
		default:
			return nil
		}
	}
	return nil
}

var DeleteUserQuery = "DELETE FROM users WHERE id=UUID_TO_BIN(?)"

func (MYSQLFunctions) DeleteUser(ID *uuid.UUID, tx *sql.Tx) error {
	result, err := tx.Exec(DeleteUserQuery, ID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		if errRb := tx.Rollback(); errRb != nil {
			return err
		}
		return ErrNoUserDeleted
	}

	return nil
}

var GetProductUserIDsQuery = "SELECT BIN_TO_UUID(users_id), privileges_id FROM users_products where products_id = UUID_TO_BIN(?)"

func (MYSQLFunctions) GetProductUserIDs(productID *uuid.UUID, tx *sql.Tx) (*models.ProductUserIDs, error) {
	rows, err := tx.Query(GetProductUserIDsQuery, productID)
	switch {
	case err == sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	defer rows.Close()
	productUsers := models.ProductUserIDs{
		UserMap:     make(map[uuid.UUID]int),
		UserIDArray: make([]uuid.UUID, 0),
	}
	for rows.Next() {
		userID := uuid.UUID{}
		privilege := -1
		err := rows.Scan(&userID, &privilege)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		productUsers.UserMap[userID] = privilege
		productUsers.UserIDArray = append(productUsers.UserIDArray, userID)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}
	return &productUsers, nil
}

var DeleteProductUserQuery = "DELETE FROM users_products where products_id = UUID_TO_BIN(?) AND users_id = UUID_TO_BIN(?)"

func (MYSQLFunctions) DeleteProductUser(productID *uuid.UUID, userID *uuid.UUID, tx *sql.Tx) error {
	result, err := tx.Exec(DeleteProductUserQuery, productID, userID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		if errRb := tx.Rollback(); errRb != nil {
			return err
		}
		return ErrNoUserWithProduct
	}

	return nil
}
