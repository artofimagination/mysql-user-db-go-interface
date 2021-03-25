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

func parseProjectData(data map[string]interface{}) (*models.ProjectData, error) {
	projectData := &models.ProjectData{}
	projectDataMap, ok := data["project"]
	if !ok {
		return nil, errors.New("Missing 'project'")
	}

	projectDataByte, err := json.Marshal(projectDataMap)
	if err != nil {
		return nil, errors.New("Invalid 'project json'")
	}

	if err := json.Unmarshal(projectDataByte, &projectData); err != nil {
		return nil, errors.New("Invalid 'project'")
	}

	return projectData, nil
}

func (c *RESTController) validateProject(expected *models.ProjectData) (int, error) {
	project, err := c.DBController.GetProject(&expected.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if diff := pretty.Diff(project, expected); len(diff) != 0 {
		return http.StatusAccepted, errors.New("Failed to update project")
	}
	return http.StatusOK, nil
}

func (c *RESTController) addProject(w ResponseWriter, r *Request) {
	log.Println("Adding project")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}
	// Parse product info
	projectJSON, ok := data["project"]
	if !ok {
		w.writeError("Missing 'project' element", http.StatusBadRequest)
		return
	}

	name, ok := projectJSON.(map[string]interface{})["name"].(string)
	if !ok {
		w.writeError("Missing 'name' element", http.StatusBadRequest)
		return
	}

	visibility, ok := projectJSON.(map[string]interface{})["visibility"].(string)
	if !ok {
		w.writeError("Missing 'visibility' element", http.StatusBadRequest)
		return
	}

	// Get user ID
	userIDString, ok := data["owner_id"]
	if !ok {
		w.writeError("Missing 'user' element", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		w.writeError("Invalid 'userID'", http.StatusBadRequest)
		return
	}

	productIDString, ok := data["product_id"]
	if !ok {
		w.writeError("Missing 'product' element", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(productIDString.(string))
	if err != nil {
		w.writeError("Invalid 'product id'", http.StatusBadRequest)
		return
	}

	project, err := c.DBController.CreateProject(name, visibility, &userID, &productID,
		func(*uuid.UUID) (string, error) {
			return testPath, nil
		})
	if err != nil {
		duplicateProject := fmt.Errorf(dbcontrollers.ErrProjectExistsString, name)
		if err.Error() == duplicateProject.Error() || err.Error() == dbcontrollers.ErrEmptyUsersList.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}

		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData(project, http.StatusCreated)
}

func (c *RESTController) getProject(w ResponseWriter, r *Request) {
	log.Println("Getting project")
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
	}

	projectData, err := c.DBController.GetProject(&id)
	if err != nil {
		if err.Error() == dbcontrollers.ErrProjectNotFound.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}

		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData(projectData, http.StatusOK)
}

func (c *RESTController) getProductProjects(w ResponseWriter, r *Request) {
	log.Println("Getting projects belonging to a product")
	if err := checkRequestType(GET, w, r); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	ids, ok := r.URL.Query()["product_id"]
	if !ok || len(ids[0]) < 1 {
		w.writeError("Url Param 'product_id' is missing", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(ids[0])
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
	}

	projectList, err := c.DBController.GetProjectsByProductID(&id)
	if err != nil {
		if err.Error() == dbcontrollers.ErrNoProjectForProduct.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData(projectList, http.StatusOK)
}

func (c *RESTController) getProjects(w ResponseWriter, r *Request) {
	log.Println("Getting multiple projects")
	if err := checkRequestType(GET, w, r); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	idList, err := parseIDList(r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	projectData, err := c.DBController.GetProjects(idList)
	if err != nil {
		if err.Error() == dbcontrollers.ErrProjectNotFound.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		err = errors.Wrap(errors.WithStack(err), "Failed to get projects")
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData(projectData, http.StatusOK)
}

func (c *RESTController) deleteProject(w ResponseWriter, r *Request) {
	log.Println("Deleting project")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	projectIDString, ok := data["id"]
	if !ok {
		w.writeError("Missing 'project_id' element", http.StatusBadRequest)
		return
	}

	projectID, err := uuid.Parse(projectIDString.(string))
	if err != nil {
		w.writeError("Invalid 'productID'", http.StatusBadRequest)
		return
	}

	err = c.DBController.DeleteProject(&projectID)
	if err != nil {
		if err.Error() == dbcontrollers.ErrProjectNotFound.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = c.DBController.GetProject(&projectID)
	if err != nil && err.Error() != dbcontrollers.ErrProjectNotFound.Error() {
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData(DataOK, http.StatusOK)
}
