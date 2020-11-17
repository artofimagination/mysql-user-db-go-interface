package mysqldb

import (
	"encoding/json"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func addOwnerships(productID *uuid.UUID, owners *models.OwnershipMap) error {
	queryString := "INSERT INTO users_products (users_id, products_id, privilege) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?)"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}

	defer db.Close()

	for userID, privilege := range *owners {
		query, err := db.Query(queryString, productID, userID, privilege)
		if err != nil {
			return err
		}

		defer query.Close()
	}
	return nil
}

func deleteOwnershipByProductID(productID *uuid.UUID) error {
	queryString := "DELETE FROM users_products where products_id = UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(queryString, productID)
	if err != nil {
		return err
	}

	return nil
}

func AddProduct(product *models.Product, owners *models.OwnershipMap) error {
	queryString := "INSERT INTO products (id, name, public, details) VALUES (UUID_TO_BIN(UUID()), ?, ?, ?)"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}

	defer db.Close()

	jsonRaw, err := ConvertToJSONRaw(product.Details)
	if err != nil {
		return err
	}

	query, err := db.Query(queryString, product.Name, product.Public, jsonRaw)
	if err != nil {
		return err
	}

	if err := addOwnerships(&product.ID, owners); err != nil {
		if err := DeleteProduct(&product.ID); err != nil {
			return errors.Wrap(errors.WithStack(err), "Failed to revert products creation")
		}
		return err
	}

	defer query.Close()
	return nil
}

func GetProductByID(ID uuid.UUID) (*models.Product, error) {
	queryString := "SELECT BIN_TO_UUID(id), name, public, details FROM products where id = UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	details := json.RawMessage{}
	product := models.Product{}
	query, err := db.Query(queryString, product.ID)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&product.ID, &product.Name, &product.Public, &product.Details); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(details, &product.Details); err != nil {
		return nil, err
	}

	return &product, nil
}

func getUsersProductIDs(userID *uuid.UUID) (*models.OwnershipMap, error) {
	queryString := "SELECT BIN_TO_UUID(products_id), privilege FROM users_products where users_id = UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(queryString, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	data := make(models.OwnershipMap)
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

func GetProductsByUserID(userID *uuid.UUID) (*[]models.Product, error) {
	ownershipMap, err := getUsersProductIDs(userID)
	if err != nil {
		return nil, err
	}

	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	products := []models.Product{}

	for productID := range *ownershipMap {
		product, err := GetProductByID(productID)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}

	return &products, nil
}

func DeleteProduct(productID *uuid.UUID) error {
	queryString := "DELETE FROM products where id = UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}

	defer db.Close()

	query, err := db.Query(queryString, productID)
	if err != nil {
		return err
	}

	if err := deleteOwnershipByProductID(productID); err != nil {
		return err
	}

	defer query.Close()

	return nil
}
