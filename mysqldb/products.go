package mysqldb

import (
	"database/sql"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrNoProductsForUser = errors.New("This user has no products")
var ErrSQLDuplicateProductNameEntryString = "Duplicate entry '%s' for key 'products.name'"
var ErrDuplicateProductNameEntry = errors.New("Product with this name already exists")
var ErrNoUserWithProduct = errors.New("No user is associated to this product")
var ErrNoProductDeleted = errors.New("No product was deleted")
var ErrNoUsersProductUpdate = errors.New("No users product was updated")

var AddProductUsersQuery = "INSERT INTO users_products (users_id, products_id, privilege) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?)"

func (MYSQLFunctions) AddProductUsers(productID *uuid.UUID, productUsers models.ProductUsers, tx *sql.Tx) error {
	for userID, privilege := range productUsers {
		_, err := tx.Exec(AddProductUsersQuery, productID, userID, privilege)
		if err != nil {
			return RollbackWithErrorStack(tx, err)
		}
	}
	return tx.Commit()
}

var DeleteProductUsersByProductIDQuery = "DELETE FROM users_products where products_id = UUID_TO_BIN(?)"

func (MYSQLFunctions) DeleteProductUsersByProductID(productID *uuid.UUID, tx *sql.Tx) error {
	result, err := tx.Exec(DeleteProductUsersByProductIDQuery, productID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		if errRb := tx.Rollback(); errRb != nil {
			return err
		}
		return ErrNoUserWithProduct
	}

	return nil
}

var UpdateUsersProductsQuery = "UPDATE users_products set privilege = ? where users_id = UUID_TO_BIN(?) AND products_id = UUID_TO_BIN(?)"

func (MYSQLFunctions) UpdateUsersProducts(userID *uuid.UUID, productID *uuid.UUID, privilege int, tx *sql.Tx) error {
	result, err := tx.Exec(UpdateUsersProductsQuery, privilege, userID, productID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		if errRb := tx.Rollback(); errRb != nil {
			return err
		}
		return ErrNoUsersProductUpdate
	}

	return nil
}

var AddProductQuery = "INSERT INTO products (id, name, public, product_details_id, product_assets_id) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?)"

func (*MYSQLFunctions) AddProduct(product *models.Product, tx *sql.Tx) error {
	// Execute transaction
	_, err := tx.Exec(AddProductQuery, product.ID, product.Name, product.Public, product.DetailsID, product.AssetsID)
	errDuplicateName := fmt.Errorf(ErrSQLDuplicateProductNameEntryString, product.Name)
	if err != nil {
		switch {
		case err.Error() == errDuplicateName.Error():
			if errRb := tx.Rollback(); errRb != nil {
				return err
			}
			return errDuplicateName
		case err != nil:
			return RollbackWithErrorStack(tx, err)
		default:
			return tx.Commit()
		}
	}
	return tx.Commit()
}

var GetProductByIDQuery = "SELECT BIN_TO_UUID(id), name, public, details, BIN_TO_UUID(product_assets_id) FROM products WHERE id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) GetProductByID(ID *uuid.UUID, tx *sql.Tx) (*models.Product, error) {
	product := models.Product{}

	query, err := tx.Query(GetProductByIDQuery, ID)
	switch {
	case err == sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&product.ID, &product.Name, &product.Public, &product.DetailsID, &product.AssetsID); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return &product, nil
}

var GetUserProductIDsQuery = "SELECT BIN_TO_UUID(products_id), privilege FROM users_products where users_id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) GetUserProductIDs(userID *uuid.UUID, tx *sql.Tx) (*models.UserProducts, error) {
	rows, err := tx.Query(GetUserProductIDsQuery, userID)
	switch {
	case err == sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	defer rows.Close()
	userProducts := models.UserProducts{
		ProductMap:     make(map[uuid.UUID]int),
		ProductIDArray: make([]uuid.UUID, 0),
	}
	for rows.Next() {
		productID := uuid.UUID{}
		privilege := -1
		err := rows.Scan(&productID, &privilege)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		userProducts.ProductMap[productID] = privilege
		userProducts.ProductIDArray = append(userProducts.ProductIDArray, productID)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}
	return &userProducts, nil
}

// GetProductsByUserID returns all products belonging to the selected user.
// The function first gets all rows matching with the user DI from users_products table,
// then gets all products based on the product ids from the first query result.
func (*MYSQLFunctions) GetProductsByUserID(userID *uuid.UUID) ([]models.Product, error) {
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	ownershipMap, err := Functions.GetUserProductIDs(userID, tx)
	switch {
	case err == sql.ErrNoRows:
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	products := []models.Product{}
	for _, productID := range ownershipMap.ProductIDArray {
		productID := productID
		product, err := Functions.GetProductByID(&productID, tx)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		products = append(products, *product)
	}

	return products, tx.Commit()
}

var GetProductByNameQuery = "SELECT BIN_TO_UUID(id), name, public, details, BIN_TO_UUID(product_assests_id) FROM products WHERE name = ?"

func (*MYSQLFunctions) GetProductByName(name string, tx *sql.Tx) (*models.Product, error) {
	product := models.Product{}

	query, err := tx.Query(GetProductByNameQuery, name)
	switch {
	case err == sql.ErrNoRows:
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&product.ID, &product.Name, &product.Public, &product.DetailsID, &product.AssetsID); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return &product, tx.Commit()
}

var DeleteProductQuery = "DELETE FROM products where id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) DeleteProduct(productID *uuid.UUID, tx *sql.Tx) error {
	result, err := tx.Exec(DeleteProductQuery, productID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		if errRb := tx.Rollback(); errRb != nil {
			return err
		}
		return ErrNoProductDeleted
	}

	return nil
}

var GetPrivilegesQuery = "SELECT id, name, description from privileges"

func (*MYSQLFunctions) GetPrivileges() (models.Privileges, error) {
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	rows, err := tx.Query(GetPrivilegesQuery)
	switch {
	case err == sql.ErrNoRows:
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	defer rows.Close()

	privileges := make(models.Privileges, 0)
	for rows.Next() {
		privilege := models.Privilege{}
		err := rows.Scan(&privilege.ID, &privilege.Name, &privilege.Description)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		privileges = append(privileges, privilege)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return privileges, tx.Commit()
}
