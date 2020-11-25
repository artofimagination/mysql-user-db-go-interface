package controllers

import (
	"errors"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
	"github.com/google/uuid"
)

var ErrProductExistsString = "Product with name %s already exists"
var ErrEmptyUsersList = errors.New("At least one product user is required")
var ErrUnknownPrivilegeString = "Unknown privilege %d set for user %s"
var ErrInvalidOwnerCount = errors.New("Product must have a single owner")

func validateUsers(users models.ProductUsers) error {
	if users == nil || (users != nil && len(users) == 0) {
		return ErrEmptyUsersList
	}

	privileges, err := mysqldb.FunctionInterface.GetPrivileges()
	if err != nil {
		return err
	}

	hasOwner := false
	for ID, privilege := range users {
		if !privileges.IsValidPrivilege(privilege) {
			return fmt.Errorf(ErrUnknownPrivilegeString, privilege, ID.String())
		}

		if privileges.IsOwnerPrivilege(privilege) {
			if hasOwner {
				return ErrInvalidOwnerCount
			}
			hasOwner = true
		}
	}

	if !hasOwner {
		return ErrInvalidOwnerCount
	}

	return nil
}

func CreateProduct(name string, public bool, users models.ProductUsers, generateAssetPath func(assetID *uuid.UUID) string) (*models.Product, error) {
	// Need to check whether the product users list is valid.
	// - is there exactly one owner
	// - are the privilege id-s valid
	if err := validateUsers(users); err != nil {
		return nil, err
	}

	references := make(models.References)
	asset, err := models.Interface.NewAsset(references, generateAssetPath)
	if err != nil {
		return nil, err
	}

	details := make(models.Details)

	product, err := models.Interface.NewProduct(name, public, details, &asset.ID)
	if err != nil {
		return nil, err
	}

	// Start a DB transaction and do all inserts within the same transaction to improve consistency.
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	existingProduct, err := mysqldb.FunctionInterface.GetProductByName(name, tx)
	if err != nil {
		return nil, err
	}

	if existingProduct != nil {
		return nil, fmt.Errorf(ErrProductExistsString, product.Name)
	}

	if err := mysqldb.FunctionInterface.AddProductUsers(&product.ID, users, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.FunctionInterface.AddAsset(mysqldb.ProductAssets, asset, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.FunctionInterface.AddProduct(product, tx); err != nil {
		return nil, err
	}

	return product, nil
}
