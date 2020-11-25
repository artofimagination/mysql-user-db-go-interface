package mysqldb

import (
	"database/sql"
	"encoding/json"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

var AddProductUsersQuery = "INSERT INTO users_products (users_id, products_id, privilege) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?)"

func (MYSQLFunctionInterface) AddProductUsers(productID *uuid.UUID, productUsers models.ProductUsers, tx *sql.Tx) error {
	for userID, privilege := range productUsers {
		_, err := tx.Exec(AddProductUsersQuery, productID, userID, privilege)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteProductUsersByProductID(productID *uuid.UUID, tx *sql.Tx) error {
	queryString := "DELETE FROM users_products where products_id = UUID_TO_BIN(?)"

	_, err := tx.Exec(queryString, productID)
	if err != nil {
		return err
	}

	return nil
}

var AddProductQuery = "INSERT INTO products (id, name, public, details, product_assets_id) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?)"

func (MYSQLFunctionInterface) AddProduct(product *models.Product, tx *sql.Tx) error {
	// Prepare data
	jsonRaw, err := ConvertToJSONRaw(product.Details)
	if err != nil {
		return err
	}

	// Execute transaction
	_, err = tx.Exec(AddProductQuery, product.ID, product.Name, product.Public, jsonRaw, product.AssetsID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	return tx.Commit()
}

var GetProductByIDQuery = "SELECT BIN_TO_UUID(id), name, public, details, BIN_TO_UUID(product_assets_id) FROM products WHERE id = UUID_TO_BIN(?)"

func getProductByID(ID uuid.UUID, tx *sql.Tx) (*models.Product, error) {
	details := json.RawMessage{}
	product := models.Product{}

	query, err := tx.Query(GetProductByIDQuery, ID)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&product.ID, &product.Name, &product.Public, &details, &product.AssetsID); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(details, &product.Details); err != nil {
		return nil, err
	}

	return &product, nil
}

func GetProductByID(ID uuid.UUID) (*models.Product, error) {
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	product, err := getProductByID(ID, tx)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return product, tx.Commit()
}

var GetUserProductIDsQuery = "SELECT BIN_TO_UUID(products_id), privilege FROM users_products where users_id = UUID_TO_BIN(?)"

func getUserProductIDs(userID uuid.UUID, tx *sql.Tx) (*models.UserProducts, error) {
	rows, err := tx.Query(GetUserProductIDsQuery, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	data := make(models.UserProducts)
	for rows.Next() {
		productID := uuid.UUID{}
		privilege := -1
		err := rows.Scan(&productID, &privilege)
		if err != nil {
			return nil, err
		}
		data[productID] = privilege
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetProductsByUserID returns all products belonging to the selected user.
// The function first gets all rows matching with the user DI from users_products table,
// then gets all products based on the product ids from the first query result.
func GetProductsByUserID(userID uuid.UUID) (*[]models.Product, error) {
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	ownershipMap, err := getUserProductIDs(userID, tx)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	products := []models.Product{}
	for productID := range *ownershipMap {
		product, err := getProductByID(productID, tx)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		products = append(products, *product)
	}

	return &products, tx.Commit()
}

var GetProductByNameQuery = "SELECT BIN_TO_UUID(id), name, public, details, BIN_TO_UUID(product_assests_id) FROM products WHERE name = ?"

func (MYSQLFunctionInterface) GetProductByName(name string, tx *sql.Tx) (*models.Product, error) {
	details := json.RawMessage{}
	product := models.Product{}

	query, err := tx.Query(GetProductByNameQuery, name)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&product.ID, &product.Name, &product.Public, &details, &product.AssetsID); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(details, &product.Details); err != nil {
		return nil, err
	}

	return &product, nil
}

func DeleteProduct(productID *uuid.UUID) error {
	queryString := "DELETE FROM products where id = UUID_TO_BIN(?)"

	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	if err := deleteProductUsersByProductID(productID, tx); err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	_, err = tx.Exec(queryString, productID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	return tx.Commit()
}

var GetPrivilegesQuery = "SELECT id, name, description from privileges"

func (MYSQLFunctionInterface) GetPrivileges() (models.Privileges, error) {
	tx, err := DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(GetPrivilegesQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	privileges := make(models.Privileges, 0)
	for rows.Next() {
		privilege := models.Privilege{}
		err := rows.Scan(&privilege.ID, &privilege.Name, &privilege.Description)
		if err != nil {
			return nil, err
		}
		privileges = append(privileges, privilege)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return privileges, nil
}
