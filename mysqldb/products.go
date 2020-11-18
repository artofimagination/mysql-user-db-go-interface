package mysqldb

import (
	"database/sql"
	"encoding/json"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

func addProductUsers(productID *uuid.UUID, productUsers *models.ProductUsers, tx *sql.Tx) error {
	queryString := "INSERT INTO users_products (users_id, products_id, privilege) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?)"

	for userID, privilege := range *productUsers {
		_, err := tx.Exec(queryString, productID, userID, privilege)
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

func AddProduct(product *models.Product, productUsers *models.ProductUsers) error {
	// Prepare data
	queryString := "INSERT INTO products (id, name, public, details) VALUES (UUID_TO_BIN(?), ?, ?, ?)"

	jsonRaw, err := ConvertToJSONRaw(product.Details)
	if err != nil {
		return err
	}

	// Execute transaction
	tx, err := DBInterface.ConnectSystem()
	if err != nil {
		return err
	}

	_, err = tx.Exec(queryString, product.ID, product.Name, product.Public, jsonRaw)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if err := addProductUsers(&product.ID, productUsers, tx); err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	return tx.Commit()
}

func getProductByID(ID uuid.UUID, tx *sql.Tx) (*models.Product, error) {
	queryString := "SELECT BIN_TO_UUID(id), name, public, details FROM products WHERE id = UUID_TO_BIN(?)"
	details := json.RawMessage{}
	product := models.Product{}

	query, err := tx.Query(queryString, ID)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&product.ID, &product.Name, &product.Public, &details); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(details, &product.Details); err != nil {
		return nil, err
	}

	return &product, nil
}

func GetProductByID(ID uuid.UUID) (*models.Product, error) {
	tx, err := DBInterface.ConnectSystem()
	if err != nil {
		return nil, err
	}

	product, err := getProductByID(ID, tx)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return product, tx.Commit()
}

func getUserProductIDs(userID uuid.UUID, tx *sql.Tx) (*models.UserProducts, error) {
	queryString := "SELECT BIN_TO_UUID(products_id), privilege FROM users_products where users_id = UUID_TO_BIN(?)"

	rows, err := tx.Query(queryString, userID)
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
	tx, err := DBInterface.ConnectSystem()
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

func DeleteProduct(productID *uuid.UUID) error {
	queryString := "DELETE FROM products where id = UUID_TO_BIN(?)"

	tx, err := DBInterface.ConnectSystem()
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
