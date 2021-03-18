package restcontrollers

import (
	"log"
	"net/http"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
)

func (c *RESTController) updateUserSettings(w ResponseWriter, r *Request) {
	log.Println("Update user settings")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	userData, err := parseUserData(data)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	err = c.DBController.UpdateUserSettings(userData)
	if err == nil {
		statusCode, err := c.validateUser(userData)
		if err != nil {
			w.writeError(err.Error(), statusCode)
			return
		}

		w.writeData(DataOK, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoUserSetttingsUpdate.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) updateUserAssets(w ResponseWriter, r *Request) {
	log.Println("Update user assets")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	userData, err := parseUserData(data)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
	}

	err = c.DBController.UpdateUserAssets(userData)
	if err == nil {
		statusCode, err := c.validateUser(userData)
		if err != nil {
			w.writeError(err.Error(), statusCode)
			return
		}

		w.writeData(DataOK, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoUserAssetsUpdate.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) updateProductDetails(w ResponseWriter, r *Request) {
	log.Println("Update product details")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	productData, err := parseProductData(data)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	err = c.DBController.UpdateProductDetails(productData)
	if err == nil {
		statusCode, err := c.validateProduct(productData)
		if err != nil {
			w.writeError(err.Error(), statusCode)
			return
		}

		w.writeData(DataOK, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoProductDetailUpdate.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) updateProductAssets(w ResponseWriter, r *Request) {
	log.Println("Update product assets")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	productData, err := parseProductData(data)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	err = c.DBController.UpdateProductAssets(productData)
	if err == nil {
		statusCode, err := c.validateProduct(productData)
		if err != nil {
			w.writeError(err.Error(), statusCode)
			return
		}

		w.writeData(DataOK, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoProductAssetUpdate.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) updateProjectDetails(w ResponseWriter, r *Request) {
	log.Println("Update project details")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	projectData, err := parseProjectData(data)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	err = c.DBController.UpdateProjectDetails(projectData)
	if err == nil {
		statusCode, err := c.validateProject(projectData)
		if err != nil {
			w.writeError(err.Error(), statusCode)
			return
		}

		w.writeData(DataOK, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoProjectDetailsUpdate.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) updateProjectAssets(w ResponseWriter, r *Request) {
	log.Println("Update project assets")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	projectData, err := parseProjectData(data)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
	}

	err = c.DBController.UpdateProjectAssets(projectData)
	if err == nil {
		statusCode, err := c.validateProject(projectData)
		if err != nil {
			w.writeError(err.Error(), statusCode)
			return
		}

		w.writeData(DataOK, statusCode)
		return
	}

	if err.Error() == dbcontrollers.ErrNoProjectAssetsUpdate.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}
