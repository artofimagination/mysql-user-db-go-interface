package dbcontrollers

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/google/uuid"
)

var ErrDuplicateEmailEntry = errors.New("User with this email already exists")
var ErrDuplicateNameEntry = errors.New("User with this name already exists")
var ErrUserNotFound = errors.New("The selected user not found")
var ErrInvalidEmailOrPasswd = errors.New("Invalid email or password")

func (MYSQLController) CreateUser(
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
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingUser != nil {
		if err := mysqldb.DBConnector.Rollback(tx); err != nil {
			return nil, err
		}
		return nil, ErrDuplicateEmailEntry
	}

	if err := mysqldb.Functions.AddAsset(mysqldb.UserAssets, asset, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddAsset(mysqldb.UserSettings, userSettings, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddUser(user, tx); err != nil {
		errDuplicateName := fmt.Errorf(mysqldb.ErrSQLDuplicateUserNameEntryString, user.Name)
		if err.Error() == errDuplicateName.Error() {
			return nil, ErrDuplicateNameEntry
		}
		return nil, err
	}

	return user, mysqldb.DBConnector.Commit(tx)
}

func (MYSQLController) DeleteUser(ID *uuid.UUID, nominatedOwners map[uuid.UUID]uuid.UUID) error {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	// Valid user
	user, err := mysqldb.Functions.GetUser(mysqldb.GetUserByIDQuery, ID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := mysqldb.DBConnector.Rollback(tx); err != nil {
				return err
			}
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

func (MYSQLController) GetUser(userID *uuid.UUID) (*models.UserData, error) {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	user, err := mysqldb.Functions.GetUser(mysqldb.GetUserByIDQuery, *userID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := mysqldb.DBConnector.Rollback(tx); err != nil {
				return nil, err
			}
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	settings, err := mysqldb.GetAsset(mysqldb.UserSettings, &user.SettingsID)
	if err != nil {
		return nil, err
	}

	assets, err := mysqldb.GetAsset(mysqldb.UserAssets, &user.AssetsID)
	if err != nil {
		return nil, err
	}

	userData := models.UserData{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Settings: *settings,
		Assets:   *assets,
	}

	return &userData, mysqldb.DBConnector.Commit(tx)
}

func (MYSQLController) UpdateUserSettings(settings *models.Asset) error {
	return mysqldb.UpdateAsset(mysqldb.UserSettings, settings)
}

func (MYSQLController) UpdateUserAssets(assets *models.Asset) error {
	return mysqldb.UpdateAsset(mysqldb.UserAssets, assets)
}

func (MYSQLController) Authenticate(email string, passwd []byte, authenticate func(string, []byte, *models.User) error) error {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	user, err := mysqldb.Functions.GetUser(mysqldb.GetUserByEmailQuery, email, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := mysqldb.DBConnector.Rollback(tx); err != nil {
				return err
			}
			return ErrInvalidEmailOrPasswd
		}
		return err
	}

	if err := authenticate(email, passwd, user); err != nil {
		return err
	}

	return mysqldb.DBConnector.Commit(tx)
}
