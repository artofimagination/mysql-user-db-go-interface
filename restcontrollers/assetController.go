package restcontrollers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
)

func (c *RESTController) updateUserSettings(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user settings")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	userData, err := parseUserData(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to get suer")
	}

	err = c.DBController.UpdateUserSettings(userData)
	if err == nil {
		statusCode, err := c.validateUser(userData)
		if err != nil {
			w.WriteHeader(statusCode)
			fmt.Fprint(w, err.Error())
			return
		}

		w.WriteHeader(statusCode)
		fmt.Fprint(w, "User settings updated")
		return
	}

	if err.Error() == dbcontrollers.ErrNoUserSetttingsUpdate.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) updateUserAssets(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user assets")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	userData, err := parseUserData(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to get suer")
	}

	err = c.DBController.UpdateUserAssets(userData)
	if err == nil {
		statusCode, err := c.validateUser(userData)
		if err != nil {
			w.WriteHeader(statusCode)
			fmt.Fprint(w, err.Error())
			return
		}

		w.WriteHeader(statusCode)
		fmt.Fprint(w, "User assets updated")
		return
	}

	if err.Error() == dbcontrollers.ErrNoUserAssetsUpdate.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) updateProductDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("Update product details")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	productData, err := parseProductData(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to get product")
	}

	err = c.DBController.UpdateProductDetails(productData)
	if err == nil {
		statusCode, err := c.validateProduct(productData)
		if err != nil {
			w.WriteHeader(statusCode)
			fmt.Fprint(w, err.Error())
			return
		}

		w.WriteHeader(statusCode)
		fmt.Fprint(w, "Product details updated")
		return
	}

	if err.Error() == dbcontrollers.ErrNoProductDetailUpdate.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) updateProductAssets(w http.ResponseWriter, r *http.Request) {
	log.Println("Update product assets")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	productData, err := parseProductData(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to get product")
	}

	err = c.DBController.UpdateProductAssets(productData)
	if err == nil {
		statusCode, err := c.validateProduct(productData)
		if err != nil {
			w.WriteHeader(statusCode)
			fmt.Fprint(w, err.Error())
			return
		}

		w.WriteHeader(statusCode)
		fmt.Fprint(w, "Product assets updated")
		return
	}

	if err.Error() == dbcontrollers.ErrNoProductAssetUpdate.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}
