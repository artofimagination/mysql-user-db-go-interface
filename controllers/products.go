package controllers

import (
	"database/sql"
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

func validateOwnership(users models.ProductUsers) error {
	if users == nil || (users != nil && len(users) == 0) {
		return ErrEmptyUsersList
	}

	privileges, err := mysqldb.Functions.GetPrivileges()
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
	if err := validateOwnership(users); err != nil {
		return nil, err
	}

	references := make(models.DataMap)
	asset, err := models.Interface.NewAsset(references, generateAssetPath)
	if err != nil {
		return nil, err
	}

	details := make(models.DataMap)
	productDetails, err := models.Interface.NewAsset(details, generateAssetPath)
	if err != nil {
		return nil, err
	}

	product, err := models.Interface.NewProduct(name, public, &productDetails.ID, &asset.ID)
	if err != nil {
		return nil, err
	}

	// Start a DB transaction and do all inserts within the same transaction to improve consistency.
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	existingProduct, err := mysqldb.Functions.GetProductByName(name, tx)
	if err != nil {
		return nil, err
	}

	if existingProduct != nil {
		return nil, fmt.Errorf(ErrProductExistsString, product.Name)
	}

	if err := mysqldb.Functions.AddProductUsers(&product.ID, users, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddAsset(mysqldb.ProductDetails, productDetails, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddAsset(mysqldb.ProductAssets, asset, tx); err != nil {
		return nil, err
	}

	if err := mysqldb.Functions.AddProduct(product, tx); err != nil {
		return nil, err
	}

	return product, nil
}

func deleteProduct(productID *uuid.UUID, tx *sql.Tx) error {
	// Valid user
	product, err := mysqldb.Functions.GetProductByID(*productID, tx)
	if err != nil {
		return err
	}

	if err := mysqldb.Functions.DeleteProduct(productID, tx); err != nil {
		return err
	}

	if err := mysqldb.Functions.DeleteAsset(mysqldb.ProductAssets, &product.AssetsID, tx); err != nil {
		return err
	}

	if err := mysqldb.Functions.DeleteAsset(mysqldb.ProductDetails, &product.DetailsID, tx); err != nil {
		return err
	}

	return nil
}

func GetProduct(productID *uuid.UUID) (*models.ProductData, error) {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	product, err := mysqldb.Functions.GetProductByID(*productID, tx)
	if err != nil {
		return nil, err
	}

	details, err := mysqldb.GetAsset(mysqldb.ProductDetails, &product.DetailsID)
	if err != nil {
		return nil, err
	}

	assets, err := mysqldb.GetAsset(mysqldb.ProductAssets, &product.AssetsID)
	if err != nil {
		return nil, err
	}

	productData := models.ProductData{
		ID:      product.ID,
		Name:    product.Name,
		Public:  product.Public,
		Details: *details,
		Assets:  *assets,
	}

	return &productData, nil
}

func UpdateProductDetails(details *models.Asset) error {
	return mysqldb.UpdateAsset(mysqldb.ProductDetails, details)
}

func UpdateProductAssets(assets *models.Asset) error {
	return mysqldb.UpdateAsset(mysqldb.ProductAssets, assets)
}
