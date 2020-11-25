package controllers

import (
	"errors"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/google/uuid"
)

var ErrUserExists = errors.New("User with this email already exists")

func CreateUser(
	name string,
	email string,
	passwd string,
	generateAssetPath func(assetID *uuid.UUID) string,
	encryptPassword func(password string) ([]byte, error)) (*models.User, error) {

	references := make(models.References)
	asset, err := models.Interface.NewAsset(references, generateAssetPath)
	if err != nil {
		return nil, err
	}

	settings := make(models.Settings)
	userSettings, err := models.Interface.NewUserSettings(settings)
	if err != nil {
		return nil, err
	}

	password, err := encryptPassword(passwd)
	if err != nil {
		return nil, err
	}

	user, err := models.Interface.NewUser(name, email, password, userSettings.ID, asset.ID)
	if err != nil {
		return nil, err
	}

	// Start a DB transaction and do all inserts within the same transaction to improve consistency.
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	existingUser, err := mysqldb.FunctionInterface.GetUserByEmail(email, tx)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrUserExists
	}

	if err := mysqldb.FunctionInterface.AddAsset(mysqldb.UserAssets, asset, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.FunctionInterface.AddSettings(userSettings, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.FunctionInterface.AddUser(user, tx); err != nil {
		return nil, err
	}

	return user, nil
}
