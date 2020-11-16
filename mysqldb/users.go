package mysqldb

import (
	"database/sql"
	"strings"

	"github.com/artofimagination/mysql-user-db-go-interface/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// GetUserByEmail returns the user defined by the email.
func GetUserByEmail(email string) (*models.User, error) {
	email = strings.ReplaceAll(email, " ", "")

	var user models.User
	queryString := "select BIN_TO_UUID(id), name, email, password, BIN_TO_UUID(user_settings_id), BIN_TO_UUID(user_assets_id) from users where email = ?"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query, err := db.Query(queryString, email)

	if err != nil {
		return nil, err
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.SettingsID, &user.AssetsID); err != nil {
		return nil, err
	}

	return &user, nil
}

// Adds a new assetID if there is non assigned yet.
// This only can happen if the user was generated before introduction of assets.
func UpdateAssetID(user *models.User) error {
	assetID, err := AddAsset()
	if err != nil {
		return err
	}

	user.AssetsID = *assetID
	if err := addUserAssetID(user); err != nil {
		return err
	}
	return nil
}

// GetUserByID returns the user defined by it uuid.
func GetUserByID(ID uuid.UUID) (*models.User, error) {
	var user models.User
	queryString := "select BIN_TO_UUID(id), name, email, password, BIN_TO_UUID(user_settings_id), BIN_TO_UUID(user_assets_id) from users where id = UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query, err := db.Query(queryString, ID)

	if err != nil {
		return nil, err
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.SettingsID, &user.AssetsID); err != nil {
		return nil, err
	}

	return &user, nil
}

func addUserAssetID(user *models.User) error {
	queryString := "UPDATE users set user_assets_id = UUID_TO_BIN(?) where id = UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}

	defer db.Close()

	query, err := db.Query(queryString, user.AssetsID, user.ID)
	if err != nil {
		return err
	}

	defer query.Close()
	return nil
}

func UserExists(username string) (bool, error) {
	var user models.User
	db, err := ConnectSystem()
	if err != nil {
		return false, err
	}
	defer db.Close()

	queryString := "SELECT name FROM users WHERE name = ?"
	queryUser := db.QueryRow(queryString, username)
	err = queryUser.Scan(&user.Name)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, err
	}
}

func EmailExists(email string) (bool, error) {
	email = strings.ReplaceAll(email, " ", "")

	var user models.User
	queryString := "select email from users where email = ?"
	db, err := ConnectSystem()
	if err != nil {
		return false, err
	}
	defer db.Close()

	queryEmail := db.QueryRow(queryString, email)

	err = queryEmail.Scan(&user.Email)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, err
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
	db, err := ConnectSystem()
	if err != nil {
		return err
	}

	defer db.Close()

	settingsID, err := AddSettings()
	if err != nil {
		return err
	}

	assetID, err := AddAsset()
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwd), 16)
	if err != nil {
		if err := DeleteSettings(settingsID); err != nil {
			return errors.Wrap(errors.WithStack(err), "Failed to revert settings creation")
		}
		return err
	}

	query, err := db.Query(queryString, name, email, hashedPassword, &settingsID, &assetID)
	if err != nil {
		if err := DeleteSettings(settingsID); err != nil {
			return errors.Wrap(errors.WithStack(err), "Failed to revert settings creation")
		}
		return err
	}

	defer query.Close()
	return nil
}

func deleteUserEntry(email string) error {
	query := "DELETE FROM users WHERE email=?"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query, email)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(email string) error {
	email = strings.ReplaceAll(email, " ", "")
	user, err := GetUserByEmail(email)
	if err != nil {
		return err
	}

	if err := deleteUserEntry(email); err != nil {
		return err
	}

	if err := DeleteSettings(&user.SettingsID); err != nil {
		return err
	}

	return nil
}
