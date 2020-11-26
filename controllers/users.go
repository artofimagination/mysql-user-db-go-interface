package controllers

import (
	"errors"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/google/uuid"
)

var ErrProductsStillAssociated = errors.New("There are stillsome products associated to this user")

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

	existingUser, err := mysqldb.Functions.GetUserByEmail(email, tx)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, mysqldb.ErrDuplicateEmailEntry
	}

	if err := mysqldb.Functions.AddAsset(mysqldb.UserAssets, asset, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddSettings(userSettings, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddUser(user, tx); err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUser(email string) error {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	user, err := mysqldb.Functions.GetUserByEmail(email, tx)
	if err != nil {
		return err
	}

	productIDs, err := mysqldb.Functions.GetUserProductIDs(user.ID, tx)
	if err != nil {
		if err != mysqldb.ErrNoProductsForUser {
			return err
		}
	}

	if productIDs != nil {
		return ErrProductsStillAssociated
	}

	return nil
}
