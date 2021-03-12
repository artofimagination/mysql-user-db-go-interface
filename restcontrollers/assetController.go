package restcontrollers

import (
	"log"
	"net/http"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
)

func (c *RESTController) updateUserSettings(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user settings")
	data, err := decodePostData(w, r)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	userData, err := parseUserData(data)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	err = c.DBController.UpdateUserSettings(userData)
	if err == nil {
		statusCode, err := c.validateUser(userData)
		if err != nil {
			writeError(err.Error(), w, statusCode)
			return
		}

		writeData(DataOK, w, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoUserSetttingsUpdate.Error() {
		writeError(err.Error(), w, http.StatusAccepted)
		return
	}
	writeError(err.Error(), w, http.StatusInternalServerError)
}

func (c *RESTController) updateUserAssets(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user assets")
	data, err := decodePostData(w, r)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	userData, err := parseUserData(data)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
	}

	err = c.DBController.UpdateUserAssets(userData)
	if err == nil {
		statusCode, err := c.validateUser(userData)
		if err != nil {
			writeError(err.Error(), w, statusCode)
			return
		}

		writeData(DataOK, w, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoUserAssetsUpdate.Error() {
		writeError(err.Error(), w, http.StatusAccepted)
		return
	}
	writeError(err.Error(), w, http.StatusInternalServerError)
}

func (c *RESTController) updateProductDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("Update product details")
	data, err := decodePostData(w, r)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	productData, err := parseProductData(data)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	err = c.DBController.UpdateProductDetails(productData)
	if err == nil {
		statusCode, err := c.validateProduct(productData)
		if err != nil {
			writeError(err.Error(), w, statusCode)
			return
		}

		writeData(DataOK, w, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoProductDetailUpdate.Error() {
		writeError(err.Error(), w, http.StatusAccepted)
		return
	}
	writeError(err.Error(), w, http.StatusInternalServerError)
}

func (c *RESTController) updateProductAssets(w http.ResponseWriter, r *http.Request) {
	log.Println("Update product assets")
	data, err := decodePostData(w, r)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	productData, err := parseProductData(data)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	err = c.DBController.UpdateProductAssets(productData)
	if err == nil {
		statusCode, err := c.validateProduct(productData)
		if err != nil {
			writeError(err.Error(), w, statusCode)
			return
		}

		writeData(DataOK, w, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoProductAssetUpdate.Error() {
		writeError(err.Error(), w, http.StatusAccepted)
		return
	}
	writeError(err.Error(), w, http.StatusInternalServerError)
}

func (c *RESTController) updateProjectDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("Update project details")
	data, err := decodePostData(w, r)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	projectData, err := parseProjectData(data)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	err = c.DBController.UpdateProjectDetails(projectData)
	if err == nil {
		statusCode, err := c.validateProject(projectData)
		if err != nil {
			writeError(err.Error(), w, statusCode)
			return
		}

		writeData(DataOK, w, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoProjectDetailsUpdate.Error() {
		writeError(err.Error(), w, http.StatusAccepted)
		return
	}
	writeError(err.Error(), w, http.StatusInternalServerError)
}

func (c *RESTController) updateProjectAssets(w http.ResponseWriter, r *http.Request) {
	log.Println("Update project assets")
	data, err := decodePostData(w, r)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	projectData, err := parseProjectData(data)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
	}

	err = c.DBController.UpdateProjectAssets(projectData)
	if err == nil {
		statusCode, err := c.validateProject(projectData)
		if err != nil {
			writeError(err.Error(), w, statusCode)
			return
		}

		writeData(DataOK, w, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoProjectAssetsUpdate.Error() {
		writeError(err.Error(), w, http.StatusAccepted)
		return
	}
	writeError(err.Error(), w, http.StatusInternalServerError)
}
