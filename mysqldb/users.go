package mysqldb

import (
	"database/sql"
	"strings"

	"github.com/artofimagination/mysql-user-db-go-interface/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func getUserByEmail(email string, tx *sql.Tx) (*models.User, error) {
	email = strings.ReplaceAll(email, " ", "")

	var user models.User
	queryString := "select BIN_TO_UUID(id), name, email, password, BIN_TO_UUID(user_settings_id), BIN_TO_UUID(user_assets_id) from users where email = ?"

	query, err := tx.Query(queryString, email)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.SettingsID, &user.AssetsID); err != nil {
		return nil, err
	}

	return &user, err
}

// GetUserByEmail returns the user defined by the email.
func GetUserByEmail(email string) (*models.User, error) {
	tx, err := DBInterface.ConnectSystem()
	if err != nil {
		return nil, err
	}

	user, err := getUserByEmail(email, tx)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return user, tx.Commit()
}

// Adds a new assetID if there is non assigned yet.
// This only can happen if the user was generated before introduction of assets.
func UpdateAssetID(user *models.User) error {
	tx, err := DBInterface.ConnectSystem()
	if err != nil {
		return err
	}

	assetID, err := AddAsset(tx)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	user.AssetsID = *assetID
	if err := addUserAssetID(user, tx); err != nil {
		return RollbackWithErrorStack(tx, err)
	}
	return nil
}

// GetUserByID returns the user defined by it uuid.
func GetUserByID(ID uuid.UUID) (*models.User, error) {
	var user models.User
	queryString := "select BIN_TO_UUID(id), name, email, password, BIN_TO_UUID(user_settings_id), BIN_TO_UUID(user_assets_id) from users where id = UUID_TO_BIN(?)"
	tx, err := DBInterface.ConnectSystem()
	if err != nil {
		return nil, err
	}

	query, err := tx.Query(queryString, ID)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.SettingsID, &user.AssetsID); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return &user, tx.Commit()
}

func addUserAssetID(user *models.User, tx *sql.Tx) error {
	queryString := "UPDATE users set user_assets_id = UUID_TO_BIN(?) where id = UUID_TO_BIN(?)"
	query, err := tx.Query(queryString, user.AssetsID, user.ID)
	if err != nil {
		return err
	}

	defer query.Close()
	return err
}

func UserExists(username string) (bool, error) {
	var user models.User
	tx, err := DBInterface.ConnectSystem()
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
	tx, err := DBInterface.ConnectSystem()
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

// AddUser creates a new user entry in the DB.
// Whitespaces in the email are automatically deleted
// Email is a unique attribute, so the function checks for existing email, before adding a new entry
func AddUser(name string, email string, passwd string) error {
	email = strings.ReplaceAll(email, " ", "")

	queryString := "INSERT INTO users (id, name, email, password, user_settings_id, user_assets_id) VALUES (UUID_TO_BIN(UUID()), ?, ?, ?, UUID_TO_BIN(?), UUID_TO_BIN(?))"
	tx, err := DBInterface.ConnectSystem()
	if err != nil {
		return err
	}

	settingsID, err := AddSettings(tx)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	assetID, err := AddAsset(tx)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwd), 16)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	query, err := tx.Query(queryString, name, email, hashedPassword, &settingsID, &assetID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	defer query.Close()
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

	tx, err := DBInterface.ConnectSystem()
	if err != nil {
		return err
	}

	user, err := getUserByEmail(email, tx)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if err := deleteUserEntry(email, tx); err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if err := deleteSettings(&user.SettingsID, tx); err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	return tx.Commit()
}
