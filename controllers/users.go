package controllers

import (
	"database/sql"
	"errors"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/google/uuid"
)

var ErrDuplicateEmailEntry = errors.New("User with this email already exists")
var ErrUserNotFound = errors.New("The selected user not found")

func CreateUser(
	name string,
	email string,
	passwd []byte,
	generateAssetPath func(assetID *uuid.UUID) string,
	encryptPassword func(password []byte) ([]byte, error)) (*models.User, error) {

	references := make(models.DataMap)
	asset, err := models.Interface.NewAsset(references, generateAssetPath)
	if err != nil {
		return nil, err
	}

	settings := make(models.DataMap)
	userSettings, err := models.Interface.NewAsset(settings, generateAssetPath)
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

	existingUser, err := mysqldb.Functions.GetUser(mysqldb.GetUserByEmailQuery, email, tx)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrDuplicateEmailEntry
	}

	if err := mysqldb.Functions.AddAsset(mysqldb.UserAssets, asset, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddAsset(mysqldb.UserSettings, userSettings, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddUser(user, tx); err != nil {
		return nil, err
	}

	return user, mysqldb.DBConnector.Commit(tx)
}

func DeleteUser(ID *uuid.UUID, nominatedOwners map[uuid.UUID]uuid.UUID) error {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	// Valid user
	user, err := mysqldb.Functions.GetUser(mysqldb.GetUserByIDQuery, ID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return err
	}

	// Has products?
	userProducts, err := mysqldb.Functions.GetUserProductIDs(user.ID, tx)
	if err != nil {
		if err != mysqldb.ErrNoProductsForUser {
			return err
		}
	}

	if userProducts != nil {
		// Handle all products
		for productID, privilege := range userProducts.ProductMap {
			privileges, err := mysqldb.Functions.GetPrivileges()
			if err != nil {
				return err
			}

			if !privileges.IsOwnerPrivilege(privilege) {
				continue
			}

			productID := productID
			// Check nominated owner
			nominated, hasNominatedOwner := nominatedOwners[productID]
			if nominatedOwners == nil || !hasNominatedOwner {
				if err := projectdb.DeleteProjects(productID); err != nil {
					return err
				}

				if err := mysqldb.Functions.DeleteProductUsersByProductID(&productID, tx); err != nil {
					return err
				}

				if err := deleteProduct(&productID, tx); err != nil {
					return err
				}
			} else {
				// Transfer ownership of the product
				if err := mysqldb.Functions.UpdateUsersProducts(&nominated, &productID, 0, tx); err != nil {
					return err
				}
			}
		}
	}

	if err := mysqldb.Functions.DeleteUser(&user.ID, tx); err != nil {
		return err
	}

	if err := mysqldb.Functions.DeleteAsset(mysqldb.UserAssets, &user.AssetsID, tx); err != nil {
		return err
	}

	if err := mysqldb.Functions.DeleteAsset(mysqldb.UserSettings, &user.SettingsID, tx); err != nil {
		return err
	}

	return mysqldb.DBConnector.Commit(tx)
}
