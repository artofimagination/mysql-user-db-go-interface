package restcontrollers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *RESTController) addProject(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding project")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}
	// Parse product info
	projectJSON, ok := data["project"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'project'")
		return
	}

	name, ok := projectJSON.(map[string]interface{})["name"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'name'")
		return
	}

	visibility, ok := projectJSON.(map[string]interface{})["visibility"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'visibility'")
		return
	}

	// Get user ID
	userIDString, ok := data["owner_id"]
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

	productIDString, ok := data["product_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'product'")
		return
	}

	productID, err := uuid.Parse(productIDString.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid 'product id'")
		return
	}

	project, err := c.DBController.CreateProject(name, visibility, &userID, &productID,
		func(*uuid.UUID) (string, error) {
			return testPath, nil
		})
	if err == nil {
		b, err := json.Marshal(project)
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

	duplicateProject := fmt.Errorf(dbcontrollers.ErrProjectExistsString, name)
	if err.Error() == duplicateProject.Error() || err.Error() == dbcontrollers.ErrEmptyUsersList.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}

	err = errors.Wrap(errors.WithStack(err), "Failed to create product")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) getProject(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting project")
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
	}

	projectData, err := c.DBController.GetProject(&id)
	if err == nil {
		b, err := json.Marshal(projectData)
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

	if err.Error() == dbcontrollers.ErrProjectNotFound.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	err = errors.Wrap(errors.WithStack(err), "Failed to get project")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) getProjects(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting multiple projects")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}

	idList, err := parseIDList(w, r)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	projectData, err := c.DBController.GetProjects(idList)
	if err == nil {
		b, err := json.Marshal(projectData)
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

	if err.Error() == dbcontrollers.ErrProjectNotFound.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	err = errors.Wrap(errors.WithStack(err), "Failed to get projects")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) deleteProject(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting project")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	projectIDString, ok := data["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'project_id'")
		return
	}

	projectID, err := uuid.Parse(projectIDString.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid 'productID'")
		return
	}

	err = c.DBController.DeleteProject(&projectID)
	if err == nil {
		_, err = c.DBController.GetProject(&projectID)
		if err != nil && err.Error() != dbcontrollers.ErrProjectNotFound.Error() {
			err = errors.Wrap(errors.WithStack(err), "Failed to get project")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Delete completed")
		return
	}

	if err.Error() == dbcontrollers.ErrProjectNotFound.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}
