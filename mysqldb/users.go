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
	ByEmail = "email"
	ByID    = "id"
)

var ErrNoUserWithEmail = errors.New("There is no user associated with this email")

var ErrSQLDuplicateUserNameEntryString = "Duplicate entry '%s' for key 'users.name'"
var ErrSQLDuplicateEmailEntryString = "Duplicate entry '%s' for key 'users.email'"
var ErrDuplicateUserNameEntry = errors.New("User with this name already exists")

var GetUserQuery = "select BIN_TO_UUID(id), name, email, password, BIN_TO_UUID(user_settings_id), BIN_TO_UUID(user_assets_id) from users where ? = ?"

// GetUser returns the user defined by the key name and key value.
// Key name can be either id or email.
func (MYSQLFunctions) GetUser(keyName string, keyValue interface{}, tx *sql.Tx) (*models.User, error) {
	switch keyName {
	case ByEmail:
		keyValue = strings.ReplaceAll(keyValue.(string), " ", "")
	case ByID:
		keyValue = keyValue.(uuid.UUID)
	}

	var user models.User
	query, err := tx.Query(GetUserQuery, keyName, keyValue)
	switch {
	case err == sql.ErrNoRows:
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, err
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.SettingsID, &user.AssetsID); err != nil {
		return nil, err
	}
	return &user, tx.Commit()
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
	_, err := tx.Exec(InsertUserQuery, user.ID, user.Name, user.Email, user.Password, user.SettingsID, user.AssetsID)
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
			return tx.Commit()
		}
	}
	return tx.Commit()
}

func deleteUserEntry(email string, tx *sql.Tx) error {
	query := "DELETE FROM users WHERE email=?"

	_, err := tx.Exec(query, email)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(email string) error {
	email = strings.ReplaceAll(email, " ", "")

	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	user, err := Functions.GetUser(ByEmail, email, tx)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if err := deleteUserEntry(email, tx); err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if err := DeleteAsset(UserAssets, &user.SettingsID, tx); err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	return tx.Commit()
}
