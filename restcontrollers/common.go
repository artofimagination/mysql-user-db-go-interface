package restcontrollers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RESTController struct {
	DBController *dbcontrollers.MYSQLController
}

const (
	POST = "POST"
	GET  = "GET"
)

var testPath = "testPath"

func checkRequestType(requestTypeString string, w http.ResponseWriter, r *http.Request) error {
	if r.Method != requestTypeString {
		w.WriteHeader(http.StatusBadRequest)
		errorString := fmt.Sprintf("Invalid request type %s", r.Method)
		fmt.Fprint(w, errorString)
		return errors.New(errorString)
	}
	return nil
}

func decodePostData(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	if err := checkRequestType(POST, w, r); err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(errors.WithStack(err), "Failed to decode request json")
		fmt.Fprint(w, err.Error())
		return nil, err
	}

	return data, nil
}

func parseIDList(w http.ResponseWriter, r *http.Request) ([]uuid.UUID, error) {
	ids, ok := r.URL.Query()["ids"]
	if !ok || len(ids[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New("Missing 'ids'")
	}

	idList := make([]uuid.UUID, 0)
	for _, idString := range ids {
		id, err := uuid.Parse(idString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil, errors.New("Invalid 'ids'")
		}
		idList = append(idList, id)
	}

	return idList, nil
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi! I am an example server!")
}

func NewRESTController() (*RESTController, error) {
	dbController, err := dbcontrollers.NewDBController()
	if err != nil {
		return nil, err
	}

	restController := &RESTController{
		DBController: dbController,
	}

	http.HandleFunc("/", sayHello)
	http.HandleFunc("/add-user", restController.addUser)
	http.HandleFunc("/get-user", restController.getUser)
	http.HandleFunc("/get-user-by-email", restController.getUserByEmail)
	http.HandleFunc("/get-users", restController.getUsers)
	http.HandleFunc("/update-user-settings", restController.updateUserSettings)
	http.HandleFunc("/update-user-assets", restController.updateUserAssets)
	http.HandleFunc("/delete-user", restController.deleteUser)
	http.HandleFunc("/authenticate", restController.authenticate)

	http.HandleFunc("/add-product-user", restController.addProductUser)
	http.HandleFunc("/delete-product-user", restController.deleteProductUser)

	http.HandleFunc("/add-product", restController.addProduct)
	http.HandleFunc("/get-product", restController.getProduct)
	http.HandleFunc("/get-products", restController.getProducts)
	http.HandleFunc("/update-product-details", restController.updateProductDetails)
	http.HandleFunc("/update-product-assets", restController.updateProductAssets)
	http.HandleFunc("/delete-product", restController.deleteProduct)

	http.HandleFunc("/add-project", restController.addProject)
	http.HandleFunc("/get-project", restController.getProject)
	http.HandleFunc("/get-projects", restController.getProjects)
	http.HandleFunc("/update-project-details", restController.updateProjectDetails)
	http.HandleFunc("/update-project-assets", restController.updateProjectAssets)
	http.HandleFunc("/get-product-projects", restController.getProductProjects)
	http.HandleFunc("/delete-project", restController.deleteProject)

	return restController, nil
}
