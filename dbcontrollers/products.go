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

func (c *MYSQLController) validateOwnership(users *models.ProductUserIDs) error {
	if users == nil || (users != nil && len(users.UserMap) == 0) {
		return ErrEmptyUsersList
	}

	privileges, err := c.DBFunctions.GetPrivileges()
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

func (c *MYSQLController) CreateProduct(name string, public bool, owner *uuid.UUID, generateAssetPath func(assetID *uuid.UUID) (string, error)) (*models.ProductData, error) {
	references := make(models.DataMap)
	asset, err := c.ModelFunctions.NewAsset(references, generateAssetPath)
	if err != nil {
		return nil, err
	}

	details := make(models.DataMap)
	productDetails, err := c.ModelFunctions.NewAsset(details, generateAssetPath)
	if err != nil {
		return nil, err
	}

	product, err := c.ModelFunctions.NewProduct(name, public, &productDetails.ID, &asset.ID)
	if err != nil {
		return nil, err
	}

	// Start a DB transaction and do all inserts within the same transaction to improve consistency.
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	existingProduct, err := c.DBFunctions.GetProductByName(name, tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingProduct != nil {
		if err := c.DBConnector.Rollback(tx); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(ErrProductExistsString, product.Name)
	}

	users := models.ProductUserIDs{
		UserIDArray: make([]uuid.UUID, 0),
		UserMap:     make(map[uuid.UUID]int),
	}
	privilege, err := c.DBFunctions.GetPrivilege("Owner")
	if err != nil {
		return nil, err
	}

	if err := c.DBFunctions.AddAsset(mysqldb.ProductDetails, productDetails, tx); err != nil {
		return nil, err
	}

	if err := c.DBFunctions.AddAsset(mysqldb.ProductAssets, asset, tx); err != nil {
		return nil, err
	}

	if err := c.DBFunctions.AddProduct(product, tx); err != nil {
		return nil, err
	}

	users.UserMap[*owner] = privilege.ID
	if err := c.DBFunctions.AddProductUsers(&product.ID, &users, tx); err != nil {
		return nil, err
	}

	productData := models.ProductData{
		ID:      product.ID,
		Name:    product.Name,
		Public:  product.Public,
		Details: productDetails,
		Assets:  asset,
	}

	return &productData, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) DeleteProduct(productID *uuid.UUID) error {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	if err := c.deleteProduct(productID, tx); err != nil {
		return err
	}

	return c.DBConnector.Commit(tx)
}

func (c *MYSQLController) deleteProduct(productID *uuid.UUID, tx *sql.Tx) error {
	product, err := c.DBFunctions.GetProductByID(productID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrProductNotFound
		}
	}

	if err := c.DBFunctions.DeleteProductUsersByProductID(productID, tx); err != nil {
		if err == mysqldb.ErrNoProductDeleted {
			return ErrProductNotFound
		}
	}

	if err := c.DBFunctions.DeleteProduct(productID, tx); err != nil {
		if err == mysqldb.ErrNoProductDeleted {
			return ErrProductNotFound
		}
		return err
	}

	if err := c.DBFunctions.DeleteAsset(mysqldb.ProductAssets, &product.AssetsID, tx); err != nil {
		return err
	}

	if err := c.DBFunctions.DeleteAsset(mysqldb.ProductDetails, &product.DetailsID, tx); err != nil {
		return err
	}

	return nil
}

func (c *MYSQLController) GetProduct(productID *uuid.UUID) (*models.ProductData, error) {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	product, err := c.DBFunctions.GetProductByID(productID, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := c.DBConnector.Rollback(tx); err != nil {
				return nil, err
			}
			return nil, ErrProductNotFound
		}
	}

	details, err := c.DBFunctions.GetAsset(mysqldb.ProductDetails, &product.DetailsID)
	if err != nil {
		return nil, err
	}

	assets, err := c.DBFunctions.GetAsset(mysqldb.ProductAssets, &product.AssetsID)
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

	return &productData, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) UpdateProductDetails(productData *models.ProductData) error {
	if err := c.DBFunctions.UpdateAsset(mysqldb.ProductDetails, productData.Details); err != nil {
		if fmt.Errorf(mysqldb.ErrAssetMissing, mysqldb.ProductDetails).Error() == err.Error() {
			return ErrMissingProductDetail
		}
		return err
	}
	return nil
}

func (c *MYSQLController) UpdateProductAssets(productData *models.ProductData) error {
	if err := c.DBFunctions.UpdateAsset(mysqldb.ProductAssets, productData.Assets); err != nil {
		if fmt.Errorf(mysqldb.ErrAssetMissing, mysqldb.ProductAssets).Error() == err.Error() {
			return ErrMissingProductAsset
		}
		return err
	}
	return nil
}

func (c *MYSQLController) UpdateProductUser(productID *uuid.UUID, userID *uuid.UUID, privilege int) error {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	return c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetProductsByUserID(userID *uuid.UUID) ([]models.UserProduct, error) {
	products := make([]models.UserProduct, 0)
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	ownershipMap, err := c.DBFunctions.GetUserProductIDs(userID, tx)
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

	return products, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetProducts(productIDs []uuid.UUID) ([]models.ProductData, error) {
	if len(productIDs) == 0 {
		return nil, ErrEmptyProductIDList
	}

	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	products, err := c.DBFunctions.GetProductsByIDs(productIDs, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := c.DBConnector.Rollback(tx); err != nil {
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

	details, err := c.DBFunctions.GetAssets(mysqldb.ProductDetails, detailsIDs, tx)
	if err != nil {
		return nil, err
	}

	assets, err := c.DBFunctions.GetAssets(mysqldb.ProductAssets, assetIDs, tx)
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

	return productDataList, c.DBConnector.Commit(tx)
}
