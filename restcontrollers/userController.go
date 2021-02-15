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

func (c *RESTController) addUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding user")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	name, ok := data["name"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'name'")
		return
	}

	email, ok := data["email"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'email'")
		return
	}

	password, ok := data["password"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'password'")
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
			err = errors.Wrap(errors.WithStack(err), "Failed to encode response")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, string(b))
		return
	}

	if err.Error() == dbcontrollers.ErrDuplicateEmailEntry.Error() ||
		err.Error() == dbcontrollers.ErrDuplicateNameEntry.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	err = errors.Wrap(errors.WithStack(err), "Failed to create user")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting user")
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

	userData, err := c.DBController.GetUser(&id)
	if err == nil {
		b, err := json.Marshal(userData)
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

	if err.Error() == dbcontrollers.ErrUserNotFound.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	err = errors.Wrap(errors.WithStack(err), "Failed to get user")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) getUserByEmail(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting user by email")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}

	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("Url Param 'email' is missing"))
		return
	}

	userData, err := c.DBController.GetUserByEmail(emails[0])
	if err != nil {
		if err.Error() == dbcontrollers.ErrUserNotFound.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		err = errors.Wrap(errors.WithStack(err), "Failed to get user")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	b, err := json.Marshal(userData)
	if err != nil {
		err = errors.Wrap(errors.WithStack(err), "Failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

func (c *RESTController) getUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting multiple users")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}

	idList, err := parseIDList(w, r)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	userData, err := c.DBController.GetUsers(idList)
	if err == nil {
		b, err := json.Marshal(userData)
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

	if err.Error() == dbcontrollers.ErrUserNotFound.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	err = errors.Wrap(errors.WithStack(err), "Failed to get users")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting user")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	// Parse data info
	userIDString, ok := data["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'id'")
		return
	}

	id, err := uuid.Parse(userIDString.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid 'id'")
		return
	}

	nominees := make(map[uuid.UUID]uuid.UUID)
	nomineesMap, ok := data["nominees"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'nominees'")
		return
	}

	for productIDString, nomineeIDString := range nomineesMap.(map[string]interface{}) {
		productID, err := uuid.Parse(productIDString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid 'product id'")
			return
		}
		nomineeID, err := uuid.Parse(nomineeIDString.(string))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid 'nominee id'")
			return
		}
		nominees[productID] = nomineeID
	}

	if err = c.DBController.DeleteUser(&id, nominees); err == nil {
		_, err = c.DBController.GetUser(&id)
		if err != nil && err.Error() != dbcontrollers.ErrUserNotFound.Error() {
			err = errors.Wrap(errors.WithStack(err), "Failed to get user")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		_, err = c.DBController.GetProduct(&id)
		if err != nil && err.Error() != dbcontrollers.ErrProductNotFound.Error() {
			err = errors.Wrap(errors.WithStack(err), "Failed to get product")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Delete completed")
		return
	}

	if err.Error() == dbcontrollers.ErrUserNotFound.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("Authenticate")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}

	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("Missing 'email'"))
		return
	}

	passwords, ok := r.URL.Query()["password"]
	if !ok || len(passwords[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("Missing 'password'"))
		return
	}

	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("Missing 'id'"))
		return
	}

	id, err := uuid.Parse(ids[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
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
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Authentication successful")
		return
	}

	if err.Error() == "Invalid password" || err.Error() == dbcontrollers.ErrUserNotFound.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	err = errors.Wrap(errors.WithStack(err), "Failed to get user")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (c *RESTController) addProductUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding product user")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	productID, err := uuid.Parse(data["product_id"].(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid 'productID'")
		return
	}

	for _, users := range data["users"].([]interface{}) {
		userData := users.(map[string]interface{})
		userID, err := uuid.Parse(userData["id"].(string))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid 'productID'")
			return
		}

		privilege := userData["privilege"].(float64)

		if err := c.DBController.AddProductUser(&productID, &userID, int(privilege)); err == nil {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, "Add product user completed")
			return
		}

		if err.Error() == dbcontrollers.ErrProductNotFound.Error() ||
			err.Error() == dbcontrollers.ErrProductUserNotAssociated.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		err = errors.Wrap(errors.WithStack(err), "Failed to add product user")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	}
}

func (c *RESTController) deleteProductUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting product user")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	productID, err := uuid.Parse(data["product_id"].(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid 'productID'")
		return
	}

	userID, err := uuid.Parse(data["user_id"].(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid 'userID'")
		return
	}

	if err := c.DBController.DeleteProductUser(&productID, &userID); err == nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Delete product user completed")
		return
	}
	if err.Error() == dbcontrollers.ErrProductUserNotAssociated.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	err = errors.Wrap(errors.WithStack(err), "Failed to delete product user")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}
