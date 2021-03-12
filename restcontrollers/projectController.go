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

func (c *RESTController) addProject(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding project")
	data, err := decodePostData(w, r)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}
	// Parse product info
	projectJSON, ok := data["project"]
	if !ok {
		writeError("Missing 'project' element", w, http.StatusBadRequest)
		return
	}

	name, ok := projectJSON.(map[string]interface{})["name"].(string)
	if !ok {
		writeError("Missing 'name' element", w, http.StatusBadRequest)
		return
	}

	visibility, ok := projectJSON.(map[string]interface{})["visibility"].(string)
	if !ok {
		writeError("Missing 'visibility' element", w, http.StatusBadRequest)
		return
	}

	// Get user ID
	userIDString, ok := data["owner_id"]
	if !ok {
		writeError("Missing 'user' element", w, http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		writeError("Invalid 'userID'", w, http.StatusBadRequest)
		return
	}

	productIDString, ok := data["product_id"]
	if !ok {
		writeError("Missing 'product' element", w, http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(productIDString.(string))
	if err != nil {
		writeError("Invalid 'product id'", w, http.StatusBadRequest)
		return
	}

	project, err := c.DBController.CreateProject(name, visibility, &userID, &productID,
		func(*uuid.UUID) (string, error) {
			return testPath, nil
		})
	if err == nil {
		b, err := json.Marshal(project)
		if err != nil {
			writeError(err.Error(), w, http.StatusInternalServerError)
			return
		}

		writeData(string(b), w, http.StatusCreated)
		return
	}

	duplicateProject := fmt.Errorf(dbcontrollers.ErrProjectExistsString, name)
	if err.Error() == duplicateProject.Error() || err.Error() == dbcontrollers.ErrEmptyUsersList.Error() {
		writeError(err.Error(), w, http.StatusAccepted)
		return
	}

	writeError(err.Error(), w, http.StatusInternalServerError)
}

func (c *RESTController) getProject(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting project")
	if err := checkRequestType(GET, w, r); err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		writeError("Url Param 'id' is missing", w, http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(ids[0])
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
	}

	projectData, err := c.DBController.GetProject(&id)
	if err == nil {
		b, err := json.Marshal(projectData)
		if err != nil {
			writeError(err.Error(), w, http.StatusInternalServerError)
			return
		}

		writeData(string(b), w, http.StatusOK)
		return
	}

	if err.Error() == dbcontrollers.ErrProjectNotFound.Error() {
		writeError(err.Error(), w, http.StatusAccepted)
		return
	}

	writeError(err.Error(), w, http.StatusInternalServerError)
}

func (c *RESTController) getProductProjects(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting projects belonging to a product")
	if err := checkRequestType(GET, w, r); err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	ids, ok := r.URL.Query()["product_id"]
	if !ok || len(ids[0]) < 1 {
		writeError("Url Param 'product_id' is missing", w, http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(ids[0])
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
	}

	projectList, err := c.DBController.GetProjectsByProductID(&id)
	if err == nil {
		b, err := json.Marshal(projectList)
		if err != nil {
			writeError(err.Error(), w, http.StatusInternalServerError)
			return
		}

		writeData(string(b), w, http.StatusOK)
		return
	}

	if err.Error() == dbcontrollers.ErrNoProjectForProduct.Error() {
		writeError(err.Error(), w, http.StatusAccepted)
		return
	}
	writeError(err.Error(), w, http.StatusInternalServerError)
}

func (c *RESTController) getProjects(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting multiple projects")
	if err := checkRequestType(GET, w, r); err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	idList, err := parseIDList(w, r)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	projectData, err := c.DBController.GetProjects(idList)
	if err == nil {
		b, err := json.Marshal(projectData)
		if err != nil {
			writeError(err.Error(), w, http.StatusInternalServerError)
			return
		}

		writeData(string(b), w, http.StatusOK)
		return
	}

	if err.Error() == dbcontrollers.ErrProjectNotFound.Error() {
		writeError(err.Error(), w, http.StatusAccepted)
		return
	}
	err = errors.Wrap(errors.WithStack(err), "Failed to get projects")
	writeError(err.Error(), w, http.StatusInternalServerError)
}

func (c *RESTController) deleteProject(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting project")
	data, err := decodePostData(w, r)
	if err != nil {
		writeError(err.Error(), w, http.StatusBadRequest)
		return
	}

	projectIDString, ok := data["id"]
	if !ok {
		writeError("Missing 'project_id' element", w, http.StatusBadRequest)
		return
	}

	projectID, err := uuid.Parse(projectIDString.(string))
	if err != nil {
		writeError("Invalid 'productID'", w, http.StatusBadRequest)
		return
	}

	err = c.DBController.DeleteProject(&projectID)
	if err == nil {
		_, err = c.DBController.GetProject(&projectID)
		if err != nil && err.Error() != dbcontrollers.ErrProjectNotFound.Error() {
			writeError(err.Error(), w, http.StatusInternalServerError)
			return
		}

		writeData(DataOK, w, http.StatusOK)
		return
	}

	if err.Error() == dbcontrollers.ErrProjectNotFound.Error() {
		writeError(err.Error(), w, http.StatusAccepted)
		return
	}
	writeError(err.Error(), w, http.StatusInternalServerError)
}
