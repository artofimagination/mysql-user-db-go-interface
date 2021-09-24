package restcontrollers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type RESTController struct {
	DBController *dbcontrollers.MYSQLController
}

const (
	UserPathAdd               = "/add-user"
	UserPathGetByID           = "/get-user-by-id"
	UserPathGetByEmail        = "/get-user-by-email"
	UserPathGetMultiple       = "/get-users"
	UserPathUpdateSettings    = "/update-user-settings"
	UserPathUpdateAssets      = "/update-user-assets"
	UserPathDeleteByID        = "/delete-user"
	UserPathAuthenticate      = "/authenticate"
	UserPathAddProductUser    = "/add-product-user"
	UserPathDeleteProductUser = "/delete-product-user"
)

const (
	ProductPathAdd           = "/add-product"
	ProductPathGetByID       = "/get-product-by-id"
	ProductPathGetMultiple   = "/get-products"
	ProductPathUpdateDetails = "/update-product-details"
	ProductPathUpdateAssets  = "/update-product-assets"
	ProductPathDeleteByID    = "/delete-product"
)

const (
	ProjectPathAdd                  = "/add-project"
	ProjectPathGetByID              = "/get-project"
	ProjectPathGetMultiple          = "/get-projects"
	ProjectPathUpdateDetails        = "/update-project-details"
	ProjectPathUpdateAssets         = "/update-project-assets"
	ProjectPathGetProductProject    = "/get-product-projects"
	ProjectPathDelete               = "/delete-project"
	ProjectPathAddViewer            = "/add-project-viewer"
	ProjectPathGetViewerByUser      = "/get-project-viewer-by-user"
	ProjectPathGetViewerByViewer    = "/get-project-viewer-by-viewer"
	ProjectPathDeleteViewerByUser   = "/delete-project-viewer-by-user"
	ProjectPathDeleteViewerByViewer = "/delete-project-viewer-by-viewer"
)

const (
	POST = "POST"
	GET  = "GET"
)

var DataOK = "OK"

type ResponseWriter struct {
	http.ResponseWriter
}

type Request struct {
	*http.Request
}

type ResponseData struct {
	Error string      `json:"error" validation:"required"`
	Data  interface{} `json:"data" validation:"required"`
}

func (w ResponseWriter) writeError(message string, statusCode int) {
	response := &ResponseData{
		Error: message,
	}
	w.writeResponse(response, statusCode)
}

func (w ResponseWriter) writeData(data interface{}, statusCode int) {
	response := &ResponseData{
		Data: data,
	}
	w.writeResponse(response, statusCode)
}

func (w ResponseWriter) writeResponse(response *ResponseData, statusCode int) {
	b, err := json.Marshal(response)
	if err != nil {
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	fmt.Fprint(w, string(b))
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

func parseIDList(r *Request) ([]uuid.UUID, error) {
	ids, ok := r.URL.Query()["ids"]
	if !ok || len(ids[0]) < 1 {
		return nil, errors.New("Missing 'ids'")
	}

	idList := make([]uuid.UUID, 0)
	for _, idString := range ids {
		id, err := uuid.Parse(idString)
		if err != nil {
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

func NewRESTController() (*mux.Router, error) {
	dbController, err := dbcontrollers.NewDBController()
	if err != nil {
		return nil, err
	}

	restController := &RESTController{
		DBController: dbController,
	}
	r := mux.NewRouter()
	r.HandleFunc("/", sayHello)
	r.HandleFunc(UserPathAdd, makeHandler(restController.addUser))
	r.HandleFunc(UserPathGetByID, makeHandler(restController.getUser))
	r.HandleFunc(UserPathGetByEmail, makeHandler(restController.getUserByEmail))
	r.HandleFunc(UserPathGetMultiple, makeHandler(restController.getUsers))
	r.HandleFunc(UserPathUpdateSettings, makeHandler(restController.updateUserSettings))
	r.HandleFunc(UserPathUpdateAssets, makeHandler(restController.updateUserAssets))
	r.HandleFunc(UserPathDeleteByID, makeHandler(restController.deleteUser))
	r.HandleFunc(UserPathAuthenticate, makeHandler(restController.authenticate))

	r.HandleFunc(UserPathAddProductUser, makeHandler(restController.addProductUser))
	r.HandleFunc(UserPathDeleteProductUser, makeHandler(restController.deleteProductUser))

	r.HandleFunc(ProductPathAdd, makeHandler(restController.addProduct))
	r.HandleFunc(ProductPathGetByID, makeHandler(restController.getProduct))
	r.HandleFunc(ProductPathGetMultiple, makeHandler(restController.getProducts))
	r.HandleFunc(ProductPathUpdateDetails, makeHandler(restController.updateProductDetails))
	r.HandleFunc(ProductPathUpdateAssets, makeHandler(restController.updateProductAssets))
	r.HandleFunc(ProductPathDeleteByID, makeHandler(restController.deleteProduct))

	r.HandleFunc(ProjectPathAdd, makeHandler(restController.addProject))
	r.HandleFunc(ProjectPathGetByID, makeHandler(restController.getProject))
	r.HandleFunc(ProjectPathGetMultiple, makeHandler(restController.getProjects))
	r.HandleFunc(ProjectPathUpdateDetails, makeHandler(restController.updateProjectDetails))
	r.HandleFunc(ProjectPathUpdateAssets, makeHandler(restController.updateProjectAssets))
	r.HandleFunc(ProjectPathGetProductProject, makeHandler(restController.getProductProjects))
	r.HandleFunc(ProjectPathDelete, makeHandler(restController.deleteProject))
	r.HandleFunc(ProjectPathAddViewer, makeHandler(restController.addProjectViewer))
	r.HandleFunc(ProjectPathGetViewerByUser, makeHandler(restController.getProjectViewersByUserID))
	r.HandleFunc(ProjectPathGetViewerByViewer, makeHandler(restController.getProjectViewersByViewerID))
	r.HandleFunc(ProjectPathDeleteViewerByUser, makeHandler(restController.deleteProjectViewerByUserID))
	r.HandleFunc(ProjectPathDeleteViewerByViewer, makeHandler(restController.deleteProjectViewerByViewerID))

	return r, nil
}
