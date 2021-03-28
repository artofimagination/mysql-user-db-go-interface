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

var testPath = "testPath"

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

func NewRESTController() (*RESTController, error) {
	dbController, err := dbcontrollers.NewDBController()
	if err != nil {
		return nil, err
	}

	restController := &RESTController{
		DBController: dbController,
	}

	http.HandleFunc("/", sayHello)
	http.HandleFunc(UserPathAdd, makeHandler(restController.addUser))
	http.HandleFunc(UserPathGetByID, makeHandler(restController.getUser))
	http.HandleFunc(UserPathGetByEmail, makeHandler(restController.getUserByEmail))
	http.HandleFunc(UserPathGetMultiple, makeHandler(restController.getUsers))
	http.HandleFunc(UserPathUpdateSettings, makeHandler(restController.updateUserSettings))
	http.HandleFunc(UserPathUpdateAssets, makeHandler(restController.updateUserAssets))
	http.HandleFunc(UserPathDeleteByID, makeHandler(restController.deleteUser))
	http.HandleFunc(UserPathAuthenticate, makeHandler(restController.authenticate))

	http.HandleFunc(UserPathAddProductUser, makeHandler(restController.addProductUser))
	http.HandleFunc(UserPathDeleteProductUser, makeHandler(restController.deleteProductUser))

	http.HandleFunc(ProductPathAdd, makeHandler(restController.addProduct))
	http.HandleFunc(ProductPathGetByID, makeHandler(restController.getProduct))
	http.HandleFunc(ProductPathGetMultiple, makeHandler(restController.getProducts))
	http.HandleFunc(ProductPathUpdateDetails, makeHandler(restController.updateProductDetails))
	http.HandleFunc(ProductPathUpdateAssets, makeHandler(restController.updateProductAssets))
	http.HandleFunc(ProductPathDeleteByID, makeHandler(restController.deleteProduct))

	http.HandleFunc(ProjectPathAdd, makeHandler(restController.addProject))
	http.HandleFunc(ProjectPathGetByID, makeHandler(restController.getProject))
	http.HandleFunc(ProjectPathGetMultiple, makeHandler(restController.getProjects))
	http.HandleFunc(ProjectPathUpdateDetails, makeHandler(restController.updateProjectDetails))
	http.HandleFunc(ProjectPathUpdateAssets, makeHandler(restController.updateProjectAssets))
	http.HandleFunc(ProjectPathGetProductProject, makeHandler(restController.getProductProjects))
	http.HandleFunc(ProjectPathDelete, makeHandler(restController.deleteProject))
	http.HandleFunc(ProjectPathAddViewer, makeHandler(restController.addProjectViewer))
	http.HandleFunc(ProjectPathGetViewerByUser, makeHandler(restController.getProjectViewersByUserID))
	http.HandleFunc(ProjectPathGetViewerByViewer, makeHandler(restController.getProjectViewersByViewerID))
	http.HandleFunc(ProjectPathDeleteViewerByUser, makeHandler(restController.deleteProjectViewerByUserID))
	http.HandleFunc(ProjectPathDeleteViewerByViewer, makeHandler(restController.deleteProjectViewerByViewerID))

	return restController, nil
}
