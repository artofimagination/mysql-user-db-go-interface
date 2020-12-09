package dbcontrollers

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
var ErrProductNotFound = errors.New("The selected product not found")
var ErrMissingProductDetail = errors.New("Details for the selected product not found")
var ErrMissingProductAsset = errors.New("Assets for the selected product not found")
var ErrEmptyProductIDList = errors.New("Request does not contain any product identifiers")

func validateOwnership(users *models.ProductUserIDs) error {
	if users == nil || (users != nil && len(users.UserMap) == 0) {
		return ErrEmptyUsersList
	}

	privileges, err := mysqldb.Functions.GetPrivileges()
	if err != nil {
		return err
	}

	hasOwner := false
	for ID, privilege := range users.UserMap {
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

func (*MYSQLController) CreateProduct(name string, public bool, owner *uuid.UUID, generateAssetPath func(assetID *uuid.UUID) string) (*models.ProductData, error) {
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
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingProduct != nil {
		if err := mysqldb.DBConnector.Rollback(tx); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(ErrProductExistsString, product.Name)
	}

	users := models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	privilege, err := mysqldb.Functions.GetPrivilege("Owner")
	if err != nil {
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

	users.UserMap[*owner] = privilege.ID
	if err := mysqldb.Functions.AddProductUsers(&product.ID, &users, tx); err != nil {
		return nil, err
	}

	productData := models.ProductData{
		ID:      product.ID,
		Name:    product.Name,
		Public:  product.Public,
		Details: productDetails,
		Assets:  asset,
	}

	return &productData, mysqldb.DBConnector.Commit(tx)
}

func (*MYSQLController) DeleteProduct(productID *uuid.UUID) error {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	if err := deleteProduct(productID, tx); err != nil {
		return err
	}

	return mysqldb.DBConnector.Commit(tx)
}

func deleteProduct(productID *uuid.UUID, tx *sql.Tx) error {
	product, err := mysqldb.Functions.GetProductByID(productID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrProductNotFound
		}
	}

	if err := mysqldb.Functions.DeleteProductUsersByProductID(productID, tx); err != nil {
		if err == mysqldb.ErrNoProductDeleted {
			return ErrProductNotFound
		}
	}

	if err := mysqldb.Functions.DeleteProduct(productID, tx); err != nil {
		if err == mysqldb.ErrNoProductDeleted {
			return ErrProductNotFound
		}
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

func (*MYSQLController) GetProduct(productID *uuid.UUID) (*models.ProductData, error) {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	product, err := mysqldb.Functions.GetProductByID(productID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := mysqldb.DBConnector.Rollback(tx); err != nil {
				return nil, err
			}
			return nil, ErrProductNotFound
		}
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
		Details: details,
		Assets:  assets,
	}

	return &productData, mysqldb.DBConnector.Commit(tx)
}

func (*MYSQLController) UpdateProductDetails(productData *models.ProductData) error {

	if err := mysqldb.UpdateAsset(mysqldb.ProductDetails, productData.Details); err != nil {
		if fmt.Errorf(mysqldb.ErrAssetMissing, mysqldb.ProductDetails).Error() == err.Error() {
			return ErrMissingProductDetail
		}
		return err
	}
	return nil
}

func (*MYSQLController) UpdateProductAssets(productData *models.ProductData) error {
	if err := mysqldb.UpdateAsset(mysqldb.ProductAssets, productData.Assets); err != nil {
		if fmt.Errorf(mysqldb.ErrAssetMissing, mysqldb.ProductAssets).Error() == err.Error() {
			return ErrMissingProductAsset
		}
		return err
	}
	return nil
}

func (*MYSQLController) UpdateProductUser(productID *uuid.UUID, userID *uuid.UUID, privilege int) error {
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	return mysqldb.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetProductsByUserID(userID *uuid.UUID) ([]models.UserProduct, error) {
	products := make([]models.UserProduct, 0)
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	ownershipMap, err := mysqldb.Functions.GetUserProductIDs(userID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoProductsForUser
		}
		return nil, err
	}

	for productID, privilege := range ownershipMap.ProductMap {
		productID := productID
		product, err := c.GetProduct(&productID)
		if err != nil {
			return nil, err
		}

		userProduct := models.UserProduct{
			ProductData: *product,
			Privilege:   privilege,
		}

		products = append(products, userProduct)
	}

	return products, mysqldb.DBConnector.Commit(tx)
}

func (*MYSQLController) GetProducts(productIDs []uuid.UUID) ([]models.ProductData, error) {
	if len(productIDs) == 0 {
		return nil, ErrEmptyProductIDList
	}

	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	products, err := mysqldb.Functions.GetProductsByIDs(productIDs, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := mysqldb.DBConnector.Rollback(tx); err != nil {
				return nil, err
			}
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	assetIDs := make([]uuid.UUID, 0)
	detailsIDs := make([]uuid.UUID, 0)
	for _, product := range products {
		assetIDs = append(assetIDs, product.AssetsID)
		detailsIDs = append(detailsIDs, product.DetailsID)
	}

	details, err := mysqldb.Functions.GetAssets(mysqldb.ProductDetails, detailsIDs, tx)
	if err != nil {
		return nil, err
	}

	assets, err := mysqldb.Functions.GetAssets(mysqldb.ProductAssets, assetIDs, tx)
	if err != nil {
		return nil, err
	}

	productDataList := make([]models.ProductData, 0)
	for index, product := range products {
		productData := models.ProductData{
			ID:      product.ID,
			Name:    product.Name,
			Public:  product.Public,
			Details: &details[index],
			Assets:  &assets[index],
		}
		productDataList = append(productDataList, productData)
	}

	return productDataList, mysqldb.DBConnector.Commit(tx)
}
