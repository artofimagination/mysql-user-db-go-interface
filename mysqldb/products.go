package mysqldb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrSQLDuplicateProductNameEntryString = "Duplicate entry '%s' for key 'products.name'"
var ErrDuplicateProductNameEntry = errors.New("Product with this name already exists")
var ErrNoUserWithProduct = errors.New("No user is associated to this product")
var ErrNoProductUserAdded = errors.New("No product user relation has been added")
var ErrNoProductDeleted = errors.New("No product was deleted")
var ErrNoUsersProductUpdate = errors.New("No users product was updated")

var AddProductUsersQuery = "INSERT INTO users_products (users_id, products_id, privileges_id) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?)"

func (*MYSQLFunctions) AddProductUsers(productID *uuid.UUID, productUsers *models.ProductUserIDs, tx *sql.Tx) error {
	for userID, privilege := range productUsers.UserMap {
		result, err := tx.Exec(AddProductUsersQuery, userID, productID, privilege)
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
			return ErrNoProductUserAdded
		}
	}
	return nil
}

var DeleteProductUsersByProductIDQuery = "DELETE FROM users_products where products_id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) DeleteProductUsersByProductID(productID *uuid.UUID, tx *sql.Tx) error {
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

var UpdateUsersProductsQuery = "UPDATE users_products set privileges_id = ? where users_id = UUID_TO_BIN(?) AND products_id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) UpdateUsersProducts(userID *uuid.UUID, productID *uuid.UUID, privilege int, tx *sql.Tx) error {
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

var AddProductQuery = "INSERT INTO products (id, name, public, product_details_id, product_assets_id) VALUES (UUID_TO_BIN(?), ?, ?, UUID_TO_BIN(?), UUID_TO_BIN(?))"

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
			return nil
		}
	}
	return nil
}

var GetProductByIDQuery = "SELECT BIN_TO_UUID(id), name, public, BIN_TO_UUID(product_details_id), BIN_TO_UUID(product_assets_id) FROM products WHERE id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) GetProductByID(ID *uuid.UUID, tx *sql.Tx) (*models.Product, error) {
	product := models.Product{}
	query := tx.QueryRow(GetProductByIDQuery, ID)
	err := query.Scan(&product.ID, &product.Name, &product.Public, &product.DetailsID, &product.AssetsID)
	switch {
	case err == sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	return &product, nil
}

var GetProductsByIDsQuery = "SELECT BIN_TO_UUID(id), name, public, BIN_TO_UUID(product_details_id), BIN_TO_UUID(product_assets_id) FROM products WHERE id IN (UUID_TO_BIN(?)"

func (*MYSQLFunctions) GetProductsByIDs(IDs []uuid.UUID, tx *sql.Tx) ([]models.Product, error) {
	query := GetProductsByIDsQuery + strings.Repeat(",UUID_TO_BIN(?)", len(IDs)-1) + ")"
	interfaceList := make([]interface{}, len(IDs))
	for i := range IDs {
		interfaceList[i] = IDs[i]
	}
	rows, err := tx.Query(query, interfaceList...)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		product := models.Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Public, &product.DetailsID, &product.AssetsID)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	if len(products) == 0 {
		return nil, sql.ErrNoRows
	}

	return products, nil
}

var GetUserProductIDsQuery = "SELECT BIN_TO_UUID(products_id), privileges_id FROM users_products where users_id = UUID_TO_BIN(?)"

func (*MYSQLFunctions) GetUserProductIDs(userID *uuid.UUID, tx *sql.Tx) (*models.UserProductIDs, error) {
	rows, err := tx.Query(GetUserProductIDsQuery, userID)
	switch {
	case err == sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	defer rows.Close()
	userProducts := models.UserProductIDs{
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

var GetProductByNameQuery = "SELECT BIN_TO_UUID(id), name, public, BIN_TO_UUID(product_details_id), BIN_TO_UUID(product_assets_id) FROM products WHERE name = ?"

func (*MYSQLFunctions) GetProductByName(name string, tx *sql.Tx) (*models.Product, error) {
	product := models.Product{}

	query := tx.QueryRow(GetProductByNameQuery, name)

	err := query.Scan(&product.ID, &product.Name, &product.Public, &product.DetailsID, &product.AssetsID)
	switch {
	case err == sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	return &product, nil
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
		privilege := &models.Privilege{}
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

var GetPrivilegeQuery = "SELECT id, name, description from privileges where name = ?"

func (*MYSQLFunctions) GetPrivilege(name string) (*models.Privilege, error) {
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	rows := tx.QueryRow(GetPrivilegeQuery, name)
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

	privilege := models.Privilege{}
	err = rows.Scan(&privilege.ID, &privilege.Name, &privilege.Description)
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

	return &privilege, tx.Commit()
}
