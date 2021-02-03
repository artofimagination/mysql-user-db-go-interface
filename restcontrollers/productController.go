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

func (c *RESTController) addProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding product")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	// Parse product info
	productJSON, ok := data["product"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'product'")
		return
	}

	name, ok := productJSON.(map[string]interface{})["name"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'name'")
		return
	}

	// Get user ID
	userIDString, ok := data["user"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'user'")
		return
	}

	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid 'userID'")
		return
	}

	product, err := c.DBController.CreateProduct(name, &userID,
		func(*uuid.UUID) (string, error) {
			return testPath, nil
		})
	if err == nil {
		b, err := json.Marshal(product)
		if err != nil {
			err = errors.Wrap(errors.WithStack(err), "Failed to encode response")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, string(b))
		return
	}

	duplicateProduct := fmt.Errorf(dbcontrollers.ErrProductExistsString, name)
	if err.Error() == duplicateProduct.Error() || err.Error() == dbcontrollers.ErrEmptyUsersList.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}

	err = errors.Wrap(errors.WithStack(err), "Failed to create product")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) getProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting product")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}

	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("Url Param 'id' is missing"))
		return
	}

	id, err := uuid.Parse(ids[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	productData, err := c.DBController.GetProduct(&id)
	if err == nil {
		b, err := json.Marshal(productData)
		if err != nil {
			err = errors.Wrap(errors.WithStack(err), "Failed to encode response")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(b))
		return
	}

	if err.Error() == dbcontrollers.ErrProductNotFound.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	err = errors.Wrap(errors.WithStack(err), "Failed to get user")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) getProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting multiple products")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}

	idList, err := parseIDList(w, r)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	productData, err := c.DBController.GetProducts(idList)
	if err == nil {
		b, err := json.Marshal(productData)
		if err != nil {
			err = errors.Wrap(errors.WithStack(err), "Failed to encode response")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(b))
		return
	}

	if err.Error() == dbcontrollers.ErrProductNotFound.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	err = errors.Wrap(errors.WithStack(err), "Failed to get products")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) deleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting product")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	productIDString, ok := data["product_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'product_id'")
		return
	}

	productID, err := uuid.Parse(productIDString.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid 'productID'")
		return
	}

	err = c.DBController.DeleteProduct(&productID)
	if err == nil {
		_, err = c.DBController.GetProduct(&productID)
		if err != nil && err.Error() != dbcontrollers.ErrProductNotFound.Error() {
			err = errors.Wrap(errors.WithStack(err), "Failed to get product")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Delete completed")
		return
	}

	if err.Error() == dbcontrollers.ErrProductNotFound.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}
