package restcontrollers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
)

func parseProductData(data map[string]interface{}) (*models.ProductData, error) {
	productData := &models.ProductData{}
	productDataMap, ok := data["product"]
	if !ok {
		return nil, errors.New("Missing 'product'")
	}

	productDataByte, err := json.Marshal(productDataMap)
	if err != nil {
		return nil, errors.New("Invalid 'product json'")
	}

	if err := json.Unmarshal(productDataByte, &productData); err != nil {
		return nil, errors.New("Invalid 'product'")
	}

	return productData, nil
}

func (c *RESTController) validateProduct(expected *models.ProductData) (int, error) {
	product, err := c.DBController.GetProduct(&expected.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if diff := pretty.Diff(product, expected); len(diff) != 0 {
		return http.StatusAccepted, errors.New("Failed to update product details")
	}
	return http.StatusOK, nil
}

func (c *RESTController) addProduct(w ResponseWriter, r *Request) {
	log.Println("Adding product")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	// Parse product info
	productJSON, ok := data["product"]
	if !ok {
		w.writeError("Missing 'product' element", http.StatusBadRequest)
		return
	}

	name, ok := productJSON.(map[string]interface{})["name"].(string)
	if !ok {
		w.writeError("Missing 'name' element", http.StatusBadRequest)
		return
	}

	// Get user ID
	userIDString, ok := data["user"]
	if !ok {
		w.writeError("Missing 'user' element", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		w.writeError("Invalid 'userId' element", http.StatusBadRequest)
		return
	}

	product, err := c.DBController.CreateProduct(name, &userID,
		func(*uuid.UUID) (string, error) {
			return testPath, nil
		})
	if err == nil {
		b, err := json.Marshal(product)
		if err != nil {
			w.writeError(err.Error(), http.StatusInternalServerError)
			return
		}

		w.writeData(string(b), http.StatusCreated)
		return
	}

	duplicateProduct := fmt.Errorf(dbcontrollers.ErrProductExistsString, name)
	if err.Error() == duplicateProduct.Error() || err.Error() == dbcontrollers.ErrEmptyUsersList.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}

	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) getProduct(w ResponseWriter, r *Request) {
	log.Println("Getting product")
	if err := checkRequestType(GET, w, r); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		w.writeError("Url Param 'id' is missing", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(ids[0])
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	productData, err := c.DBController.GetProduct(&id)
	if err == nil {
		b, err := json.Marshal(productData)
		if err != nil {
			w.writeError(err.Error(), http.StatusInternalServerError)
			return
		}

		w.writeData(string(b), http.StatusOK)
		return
	}

	if err.Error() == dbcontrollers.ErrProductNotFound.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}

	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) getProducts(w ResponseWriter, r *Request) {
	log.Println("Getting multiple products")
	if err := checkRequestType(GET, w, r); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	idList, err := parseIDList(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	productData, err := c.DBController.GetProducts(idList)
	if err == nil {
		b, err := json.Marshal(productData)
		if err != nil {
			w.writeError(err.Error(), http.StatusInternalServerError)
			return
		}

		w.writeData(string(b), http.StatusOK)
		return
	}

	if err.Error() == dbcontrollers.ErrProductNotFound.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) deleteProduct(w ResponseWriter, r *Request) {
	log.Println("Deleting product")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	productIDString, ok := data["product_id"]
	if !ok {
		w.writeError("Missing 'product_id' element", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(productIDString.(string))
	if err != nil {
		w.writeError("Invalid 'product_id' element", http.StatusBadRequest)
		return
	}

	err = c.DBController.DeleteProduct(&productID)
	if err == nil {
		_, err = c.DBController.GetProduct(&productID)
		if err != nil && err.Error() != dbcontrollers.ErrProductNotFound.Error() {
			w.writeError(err.Error(), http.StatusInternalServerError)
			return
		}

		w.writeData(DataOK, http.StatusOK)
		return
	}

	if err.Error() == dbcontrollers.ErrProductNotFound.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}
