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
var ErrNoProductsForUser = errors.New("This user has no products")
var ErrNoProjectsForUser = errors.New("This user has no projects")
var ErrProductUserNotAssociated = errors.New("Unable to associate the product with the selected user")
var ErrNoUserSetttingsUpdate = errors.New("Settings for the selected user not found or no change happened")
var ErrNoUserAssetsUpdate = errors.New("Assets for the selected user not found or no change happened")
var ErrEmptyUserIDList = errors.New("Request does not contain any user identifiers")

var ErrMissingUserDBString = "Error 1452: Cannot add or update a child row: a foreign key constraint fails (`user_database`.`users_products`, CONSTRAINT `users_products_ibfk_2` FOREIGN KEY (`users_id`) REFERENCES `users` (`id`))"

func (c *MYSQLController) CreateUser(
	name string,
	email string,
	passwd []byte,
	generateAssetPath func(assetID *uuid.UUID) (string, error),
	encryptPassword func(password []byte) ([]byte, error)) (*models.UserData, error) {

	references := make(models.DataMap)
	asset, err := c.ModelFunctions.NewAsset(references, generateAssetPath)
	if err != nil {
		return nil, err
	}

	settings := make(models.DataMap)
	userSettings, err := c.ModelFunctions.NewAsset(settings, generateAssetPath)
	if err != nil {
		return nil, err
	}

	password, err := encryptPassword(passwd)
	if err != nil {
		return nil, err
	}

	user, err := c.ModelFunctions.NewUser(name, email, password, userSettings.ID, asset.ID)
	if err != nil {
		return nil, err
	}

	// Start a DB transaction and do all inserts within the same transaction to improve consistency.
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	existingUser, err := c.DBFunctions.GetUser(mysqldb.ByEmail, email, tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingUser != nil {
		if err := c.DBConnector.Rollback(tx); err != nil {
			return nil, err
		}
		return nil, ErrDuplicateEmailEntry
	}

	if err := c.DBFunctions.AddAsset(mysqldb.UserAssets, asset, tx); err != nil {
		return nil, err
	}

	if err := c.DBFunctions.AddAsset(mysqldb.UserSettings, userSettings, tx); err != nil {
		return nil, err
	}

	if err := c.DBFunctions.AddUser(user, tx); err != nil {
		errDuplicateName := fmt.Errorf(mysqldb.ErrSQLDuplicateUserNameEntryString, user.Name)
		if err.Error() == errDuplicateName.Error() {
			return nil, ErrDuplicateNameEntry
		}
		if err.Error() == ErrMissingUserDBString {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	userData := models.UserData{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Settings: userSettings,
		Assets:   asset,
	}

	return &userData, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) DeleteUser(ID *uuid.UUID, nominatedOwners map[uuid.UUID]uuid.UUID) error {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	// Valid user
	user, err := c.DBFunctions.GetUser(mysqldb.ByID, ID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := c.DBConnector.Rollback(tx); err != nil {
				return err
			}
			return ErrUserNotFound
		}
		return err
	}

	// Has products?
	userProducts, err := c.DBFunctions.GetUserProductIDs(&user.ID, tx)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if userProducts != nil {
		// Handle all products
		for productID, privilege := range userProducts.ProductMap {
			privileges, err := c.DBFunctions.GetPrivileges()
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
				if err := c.DBFunctions.DeleteProjectsByProductID(&productID, tx); err != nil && err != mysqldb.ErrNoProjectDeleted {
					return err
				}

				if err := c.deleteProduct(&productID, tx); err != nil {
					return err
				}
			} else {
				// Transfer ownership of the product
				if err := c.DBFunctions.UpdateUsersProducts(&nominated, &productID, 0, tx); err != nil {
					return err
				}
				if err := c.DBFunctions.DeleteProductUser(&productID, ID, tx); err != nil {
					if err == mysqldb.ErrNoUserWithProduct {
						return ErrProductUserNotAssociated
					}
					return err
				}
			}
		}
	}

	if err := c.DBFunctions.DeleteUser(&user.ID, tx); err != nil {
		return err
	}

	if err := c.DBFunctions.DeleteAsset(mysqldb.UserAssets, &user.AssetsID, tx); err != nil {
		return err
	}

	if err := c.DBFunctions.DeleteAsset(mysqldb.UserSettings, &user.SettingsID, tx); err != nil {
		return err
	}

	return c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetUser(userID *uuid.UUID) (*models.UserData, error) {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	user, err := c.DBFunctions.GetUser(mysqldb.ByID, *userID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := c.DBConnector.Rollback(tx); err != nil {
				return nil, err
			}
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	settings, err := c.DBFunctions.GetAsset(mysqldb.UserSettings, &user.SettingsID)
	if err != nil {
		return nil, err
	}

	assets, err := c.DBFunctions.GetAsset(mysqldb.UserAssets, &user.AssetsID)
	if err != nil {
		return nil, err
	}

	userData := models.UserData{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Settings: settings,
		Assets:   assets,
	}

	return &userData, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetUserByEmail(email string) (*models.UserData, error) {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	user, err := c.DBFunctions.GetUser(mysqldb.ByEmail, email, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := c.DBConnector.Rollback(tx); err != nil {
				return nil, err
			}
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	settings, err := c.DBFunctions.GetAsset(mysqldb.UserSettings, &user.SettingsID)
	if err != nil {
		return nil, err
	}

	assets, err := c.DBFunctions.GetAsset(mysqldb.UserAssets, &user.AssetsID)
	if err != nil {
		return nil, err
	}

	userData := models.UserData{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Settings: settings,
		Assets:   assets,
	}

	return &userData, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetUsers(userIDs []uuid.UUID) ([]models.UserData, error) {
	if len(userIDs) == 0 {
		return nil, ErrEmptyUserIDList
	}

	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	users, err := c.DBFunctions.GetUsersByIDs(userIDs, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := c.DBConnector.Rollback(tx); err != nil {
				return nil, err
			}
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	assetIDs := make([]uuid.UUID, 0)
	settingsIDs := make([]uuid.UUID, 0)
	for _, user := range users {
		assetIDs = append(assetIDs, user.AssetsID)
		settingsIDs = append(settingsIDs, user.SettingsID)
	}

	settings, err := c.DBFunctions.GetAssets(mysqldb.UserSettings, settingsIDs, tx)
	if err != nil {
		return nil, err
	}

	assets, err := c.DBFunctions.GetAssets(mysqldb.UserAssets, assetIDs, tx)
	if err != nil {
		return nil, err
	}

	userDataList := make([]models.UserData, 0)
	for index, user := range users {
		userData := models.UserData{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Settings: &settings[index],
			Assets:   &assets[index],
		}
		userDataList = append(userDataList, userData)
	}

	return userDataList, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) UpdateUserSettings(userData *models.UserData) error {
	if err := c.DBFunctions.UpdateAsset(mysqldb.UserSettings, userData.Settings); err != nil {
		if fmt.Errorf(mysqldb.ErrAssetMissing, mysqldb.UserSettings).Error() == err.Error() {
			return ErrNoUserSetttingsUpdate
		}
		return err
	}
	return nil
}

func (c *MYSQLController) UpdateUserAssets(userData *models.UserData) error {
	if err := c.DBFunctions.UpdateAsset(mysqldb.UserAssets, userData.Assets); err != nil {
		if fmt.Errorf(mysqldb.ErrAssetMissing, mysqldb.UserAssets).Error() == err.Error() {
			return ErrNoUserAssetsUpdate
		}
		return err
	}
	return nil
}

func (c *MYSQLController) Authenticate(
	userID *uuid.UUID,
	email string,
	password string,
	authenticate func(string, string, *models.User) error) error {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	user, err := c.DBFunctions.GetUser(mysqldb.ByID, userID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := c.DBConnector.Rollback(tx); err != nil {
				return err
			}
			return ErrUserNotFound
		}
		return err
	}

	if err := authenticate(email, password, user); err != nil {
		return err
	}

	return c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetUsersByProductID(productID *uuid.UUID) ([]models.ProductUser, error) {
	users := make([]models.ProductUser, 0)
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	ownershipMap, err := c.DBFunctions.GetProductUserIDs(productID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoProductsForUser
		}
		return nil, err
	}

	for userID, privilege := range ownershipMap.UserMap {
		userID := userID
		user, err := c.GetUser(&userID)
		if err != nil {
			return nil, err
		}

		productUser := models.ProductUser{
			UserData:  *user,
			Privilege: privilege,
		}

		users = append(users, productUser)
	}

	return users, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) AddProductUser(productID *uuid.UUID, userID *uuid.UUID, privilege int) error {
	productUsers := models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	productUsers.UserMap[*userID] = privilege

	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	if err := c.DBFunctions.AddProductUsers(productID, &productUsers, tx); err != nil {
		if err == mysqldb.ErrNoProductUserAdded {
			return ErrProductUserNotAssociated
		}
		return err
	}

	return c.DBConnector.Commit(tx)
}

func (c *MYSQLController) DeleteProductUser(productID *uuid.UUID, userID *uuid.UUID) error {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	if err := c.DBFunctions.DeleteProductUser(productID, userID, tx); err != nil {
		if err == mysqldb.ErrNoUserWithProduct {
			return ErrProductUserNotAssociated
		}
		return err
	}

	return c.DBConnector.Commit(tx)
}
