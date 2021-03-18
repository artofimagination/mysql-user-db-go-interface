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

func (c *RESTController) validateUser(expected *models.UserData) (int, error) {
	user, err := c.DBController.GetUser(&expected.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if diff := pretty.Diff(user, expected); len(diff) != 0 {
		return http.StatusAccepted, errors.New("Failed to update user details")
	}
	return http.StatusOK, nil
}

func parseUserData(data map[string]interface{}) (*models.UserData, error) {
	userData := models.UserData{}
	userDataMap, ok := data["user"]
	if !ok {
		return nil, errors.New("Missing 'user'")
	}

	userDataByte, err := json.Marshal(userDataMap)
	if err != nil {
		return nil, errors.New("Invalid 'user json'")
	}

	if err := json.Unmarshal(userDataByte, &userData); err != nil {
		return nil, errors.New("Invalid 'user'")
	}
	return &userData, nil
}

func (c *RESTController) addUser(w ResponseWriter, r *Request) {
	log.Println("Adding user")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	name, ok := data["name"].(string)
	if !ok {
		w.writeError("Missing 'name' element", http.StatusBadRequest)
		return
	}

	email, ok := data["email"].(string)
	if !ok {
		w.writeError("Missing 'email' element", http.StatusBadRequest)
		return
	}

	password, ok := data["password"].(string)
	if !ok {
		w.writeError("Missing 'password' element", http.StatusBadRequest)
		return
	}

	// Execute function
	user, err := c.DBController.CreateUser(name, email, []byte(password),
		func(*uuid.UUID) (string, error) {
			return testPath, nil
		}, func(password []byte) ([]byte, error) {
			return password, nil
		})
	if err == nil {
		b, err := json.Marshal(user)
		if err != nil {
			w.writeError(err.Error(), http.StatusInternalServerError)
			return
		}

		w.writeData(string(b), http.StatusCreated)
		return
	}

	if err.Error() == dbcontrollers.ErrDuplicateEmailEntry.Error() ||
		err.Error() == dbcontrollers.ErrDuplicateNameEntry.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	err = errors.Wrap(errors.WithStack(err), "Failed to create user")
	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) getUser(w ResponseWriter, r *Request) {
	log.Println("Getting user")
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
		return
	}

	userData, err := c.DBController.GetUser(&id)
	if err == nil {
		b, err := json.Marshal(userData)
		if err != nil {
			w.writeError(err.Error(), http.StatusInternalServerError)
			return
		}

		w.writeData(string(b), http.StatusOK)
		return
	}

	if err.Error() == dbcontrollers.ErrUserNotFound.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) getUserByEmail(w ResponseWriter, r *Request) {
	log.Println("Getting user by email")
	if err := checkRequestType(GET, w, r); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		w.writeError("Url Param 'email' is missing", http.StatusBadRequest)
		return
	}

	userData, err := c.DBController.GetUserByEmail(emails[0])
	if err != nil {
		if err.Error() == dbcontrollers.ErrUserNotFound.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(userData)
	if err != nil {
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData(string(b), http.StatusOK)
}

func (c *RESTController) getUsers(w ResponseWriter, r *Request) {
	log.Println("Getting multiple users")
	if err := checkRequestType(GET, w, r); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	idList, err := parseIDList(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	userData, err := c.DBController.GetUsers(idList)
	if err == nil {
		b, err := json.Marshal(userData)
		if err != nil {
			w.writeError(err.Error(), http.StatusInternalServerError)
			return
		}

		w.writeData(string(b), http.StatusOK)
		return
	}

	if err.Error() == dbcontrollers.ErrUserNotFound.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) deleteUser(w ResponseWriter, r *Request) {
	log.Println("Deleting user")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	// Parse data info
	userIDString, ok := data["id"]
	if !ok {
		w.writeError("Missing 'id' element", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(userIDString.(string))
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	nominees := make(map[uuid.UUID]uuid.UUID)
	nomineesMap, ok := data["nominees"]
	if !ok {
		w.writeError("Missing 'nominees' element", http.StatusBadRequest)
		return
	}

	for productIDString, nomineeIDString := range nomineesMap.(map[string]interface{}) {
		productID, err := uuid.Parse(productIDString)
		if err != nil {
			w.writeError(err.Error(), http.StatusBadRequest)
			return
		}
		nomineeID, err := uuid.Parse(nomineeIDString.(string))
		if err != nil {
			w.writeError(err.Error(), http.StatusBadRequest)
			return
		}
		nominees[productID] = nomineeID
	}

	if err = c.DBController.DeleteUser(&id, nominees); err == nil {
		_, err = c.DBController.GetUser(&id)
		if err != nil && err.Error() != dbcontrollers.ErrUserNotFound.Error() {
			w.writeError(err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = c.DBController.GetProduct(&id)
		if err != nil && err.Error() != dbcontrollers.ErrProductNotFound.Error() {
			w.writeError(err.Error(), http.StatusInternalServerError)
			return
		}

		w.writeData(DataOK, http.StatusOK)
		return
	}

	if err.Error() == dbcontrollers.ErrUserNotFound.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) authenticate(w ResponseWriter, r *Request) {
	log.Println("Authenticate")
	if err := checkRequestType(GET, w, r); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		w.writeError("Missing 'email'", http.StatusBadRequest)
		return
	}

	passwords, ok := r.URL.Query()["password"]
	if !ok || len(passwords[0]) < 1 {
		w.writeError("Missing 'password' element", http.StatusBadRequest)
		return
	}

	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		w.writeError("Missing 'id' element", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(ids[0])
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	err = c.DBController.Authenticate(&id, emails[0], passwords[0],
		func(string, pass string, user *models.User) error {
			if diff := pretty.Diff([]byte(pass), user.Password); len(diff) != 0 {
				return errors.New("Invalid password")
			}
			return nil
		})
	if err == nil {
		w.writeData(DataOK, http.StatusOK)
		return
	}

	if err.Error() == "Invalid password" || err.Error() == dbcontrollers.ErrUserNotFound.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}

func (c *RESTController) addProductUser(w ResponseWriter, r *Request) {
	log.Println("Adding product user")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(data["product_id"].(string))
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	for _, users := range data["users"].([]interface{}) {
		userData := users.(map[string]interface{})
		userID, err := uuid.Parse(userData["id"].(string))
		if err != nil {
			w.writeError(err.Error(), http.StatusBadRequest)
			return
		}

		privilege := userData["privilege"].(float64)

		if err := c.DBController.AddProductUser(&productID, &userID, int(privilege)); err == nil {
			w.writeData(DataOK, http.StatusCreated)
			return
		}

		if err.Error() == dbcontrollers.ErrProductNotFound.Error() ||
			err.Error() == dbcontrollers.ErrProductUserNotAssociated.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		w.writeError(err.Error(), http.StatusInternalServerError)
	}
}

func (c *RESTController) deleteProductUser(w ResponseWriter, r *Request) {
	log.Println("Deleting product user")
	data, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(data["product_id"].(string))
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(data["user_id"].(string))
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.DBController.DeleteProductUser(&productID, &userID); err == nil {
		w.writeData(DataOK, http.StatusOK)
		return
	}
	if err.Error() == dbcontrollers.ErrProductUserNotAssociated.Error() {
		w.writeError(err.Error(), http.StatusAccepted)
		return
	}
	w.writeError(err.Error(), http.StatusInternalServerError)
}
