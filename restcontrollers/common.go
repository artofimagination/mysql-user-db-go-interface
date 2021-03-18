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

var DataOK = "\"OK\""

type ResponseWriter struct {
	http.ResponseWriter
}

type Request struct {
	*http.Request
}

func (w ResponseWriter) writeError(message string, statusCode int) {
	w.writeResponse(fmt.Sprintf("{\"error\":\"%s\"}", message), statusCode)
}

func (w ResponseWriter) writeData(data string, statusCode int) {
	w.writeResponse(fmt.Sprintf("{\"data\": %s}", data), statusCode)
}

func (w ResponseWriter) writeResponse(data string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprint(w, data)
}

func checkRequestType(requestTypeString string, w ResponseWriter, r *Request) error {
	if r.Method != requestTypeString {
		w.WriteHeader(http.StatusBadRequest)
		errorString := fmt.Sprintf("Invalid request type %s", r.Method)
		return errors.New(errorString)
	}
	return nil
}

func decodePostData(w ResponseWriter, r *Request) (map[string]interface{}, error) {
	if err := checkRequestType(POST, w, r); err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(errors.WithStack(err), "Failed to decode request json")
		return nil, err
	}

	return data, nil
}

func parseIDList(w ResponseWriter, r *Request) ([]uuid.UUID, error) {
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
	fmt.Fprintln(w, "Hi! I am a user database server!")
}

func makeHandler(fn func(ResponseWriter, *Request)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		r := &Request{request}
		w := ResponseWriter{writer}
		fn(w, r)
	}
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
	http.HandleFunc("/add-user", makeHandler(restController.addUser))
	http.HandleFunc("/get-user", makeHandler(restController.getUser))
	http.HandleFunc("/get-user-by-email", makeHandler(restController.getUserByEmail))
	http.HandleFunc("/get-users", makeHandler(restController.getUsers))
	http.HandleFunc("/update-user-settings", makeHandler(restController.updateUserSettings))
	http.HandleFunc("/update-user-assets", makeHandler(restController.updateUserAssets))
	http.HandleFunc("/delete-user", makeHandler(restController.deleteUser))
	http.HandleFunc("/authenticate", makeHandler(restController.authenticate))

	http.HandleFunc("/add-product-user", makeHandler(restController.addProductUser))
	http.HandleFunc("/delete-product-user", makeHandler(restController.deleteProductUser))

	http.HandleFunc("/add-product", makeHandler(restController.addProduct))
	http.HandleFunc("/get-product", makeHandler(restController.getProduct))
	http.HandleFunc("/get-products", makeHandler(restController.getProducts))
	http.HandleFunc("/update-product-details", makeHandler(restController.updateProductDetails))
	http.HandleFunc("/update-product-assets", makeHandler(restController.updateProductAssets))
	http.HandleFunc("/delete-product", makeHandler(restController.deleteProduct))

	http.HandleFunc("/add-project", makeHandler(restController.addProject))
	http.HandleFunc("/get-project", makeHandler(restController.getProject))
	http.HandleFunc("/get-projects", makeHandler(restController.getProjects))
	http.HandleFunc("/update-project-details", makeHandler(restController.updateProjectDetails))
	http.HandleFunc("/update-project-assets", makeHandler(restController.updateProjectAssets))
	http.HandleFunc("/get-product-projects", makeHandler(restController.getProductProjects))
	http.HandleFunc("/delete-project", makeHandler(restController.deleteProject))

	return restController, nil
}
