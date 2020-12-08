package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var dbController *dbcontrollers.MYSQLController

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi! I am an example server!")
}

const (
	POST = "POST"
	GET  = "GET"
)

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

func addUser(w http.ResponseWriter, r *http.Request) {
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
	user, err := dbController.CreateUser(name, email, []byte(password),
		func(*uuid.UUID) string {
			return "testPath"
		}, func(password []byte) ([]byte, error) {
			return password, nil
		})
	if err != nil {
		if err.Error() == dbcontrollers.ErrDuplicateEmailEntry.Error() ||
			err.Error() == dbcontrollers.ErrDuplicateNameEntry.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		err = errors.Wrap(errors.WithStack(err), "Failed to create user")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		err = errors.Wrap(errors.WithStack(err), "Failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(b))
}

func queryUser(r *http.Request) (*models.UserData, error) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		return nil, errors.New("Url Param 'id' is missing")
	}

	id, err := uuid.Parse(ids[0])
	if err != nil {
		return nil, err
	}

	userData, err := dbController.GetUser(&id)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting user")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}

	userData, err := queryUser(r)
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

func getUserByEmail(w http.ResponseWriter, r *http.Request) {
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

	userData, err := dbController.GetUserByEmail(emails[0])
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

func getUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting multiple users")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}

	idList, err := parseIDList(w, r)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	userData, err := dbController.GetUsers(idList)
	if err != nil {
		if err.Error() == dbcontrollers.ErrUserNotFound.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		err = errors.Wrap(errors.WithStack(err), "Failed to get users")
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

func deleteUser(w http.ResponseWriter, r *http.Request) {
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

	err = dbController.DeleteUser(&id, nominees)
	if err != nil {
		if err.Error() == dbcontrollers.ErrUserNotFound.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	_, err = dbController.GetUser(&id)
	if err != nil && err.Error() != dbcontrollers.ErrUserNotFound.Error() {
		err = errors.Wrap(errors.WithStack(err), "Failed to get user")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Delete completed")
}

func authenticate(w http.ResponseWriter, r *http.Request) {
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

	err = dbController.Authenticate(&id, emails[0], passwords[0],
		func(string, pass string, user *models.User) error {
			if !cmp.Equal([]byte(pass), user.Password) {
				return errors.New("Invalid password")
			}
			return nil
		})
	if err != nil {
		if err.Error() == "Invalid password" || err.Error() == dbcontrollers.ErrUserNotFound.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		err = errors.Wrap(errors.WithStack(err), "Failed to get user")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Authentication successful")
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

func validateUser(expected *models.UserData) (int, error) {
	user, err := dbController.GetUser(&expected.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !cmp.Equal(user, expected) {
		return http.StatusAccepted, errors.New("Failed to update user details")
	}
	return http.StatusOK, nil
}

func updateUserSettings(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user settings")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	userData, err := parseUserData(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to get suer")
	}

	err = dbController.UpdateUserSettings(userData)
	if err != nil {
		if err.Error() == dbcontrollers.ErrMissingUserSettings.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	statusCode, err := validateUser(userData)
	if err != nil {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(statusCode)
	fmt.Fprint(w, "User settings updated")
}

func updateUserAssets(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user assets")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	userData, err := parseUserData(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to get suer")
	}

	err = dbController.UpdateUserAssets(userData)
	if err != nil {
		if err.Error() == dbcontrollers.ErrMissingUserAssets.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	statusCode, err := validateUser(userData)
	if err != nil {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(statusCode)
	fmt.Fprint(w, "User assets updated")
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding product")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	// Parse product info
	productJSON, ok := data["product"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'product'")
		return
	}

	name, ok := productJSON.(map[string]interface{})["name"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'name'")
		return
	}

	public, ok := productJSON.(map[string]interface{})["public"].(bool)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'public'")
		return
	}

	// Get user ID
	userIDString, ok := data["user"]
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

	product, err := dbController.CreateProduct(name, public, &userID,
		func(*uuid.UUID) string {
			return "testPath"
		})
	if err != nil {
		duplicateProduct := fmt.Errorf(dbcontrollers.ErrProductExistsString, name)
		if err.Error() == duplicateProduct.Error() || err.Error() == dbcontrollers.ErrEmptyUsersList.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}

		err = errors.Wrap(errors.WithStack(err), "Failed to create product")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	b, err := json.Marshal(product)
	if err != nil {
		err = errors.Wrap(errors.WithStack(err), "Failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(b))
}

func queryProduct(r *http.Request) (*models.ProductData, error) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		return nil, errors.New("Url Param 'id' is missing")
	}

	id, err := uuid.Parse(ids[0])
	if err != nil {
		return nil, err
	}

	productData, err := dbController.GetProduct(&id)
	if err != nil {
		return nil, err
	}
	return productData, nil
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting product")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}

	productData, err := queryProduct(r)
	if err != nil {
		if err.Error() == dbcontrollers.ErrProductNotFound.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		err = errors.Wrap(errors.WithStack(err), "Failed to get user")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	b, err := json.Marshal(productData)
	if err != nil {
		err = errors.Wrap(errors.WithStack(err), "Failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(b))
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting multiple products")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}

	idList, err := parseIDList(w, r)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	productData, err := dbController.GetProducts(idList)
	if err != nil {
		if err.Error() == dbcontrollers.ErrProductNotFound.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		err = errors.Wrap(errors.WithStack(err), "Failed to get products")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	b, err := json.Marshal(productData)
	if err != nil {
		err = errors.Wrap(errors.WithStack(err), "Failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

func parseProductData(data map[string]interface{}) (*models.ProductData, error) {
	productData := models.ProductData{}
	productDataMap, ok := data["product"]
	if !ok {
		return nil, errors.New("Missing 'product'")
	}

	productDataByte, err := json.Marshal(productDataMap)
	if err != nil {
		return nil, errors.New("Invalid 'product json'")
	}

	if err := json.Unmarshal(productDataByte, &productData); err != nil {
		return nil, errors.New("Invalid 'product'")
	}
	return &productData, nil
}

func validateProduct(expected *models.ProductData) (int, error) {
	product, err := dbController.GetProduct(&expected.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !cmp.Equal(product, expected) {
		return http.StatusAccepted, errors.New("Failed to update product details")
	}
	return http.StatusOK, nil
}

func updateProductDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("Update product details")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	productData, err := parseProductData(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to get product")
	}

	err = dbController.UpdateProductDetails(productData)
	if err != nil {
		if err.Error() == dbcontrollers.ErrMissingProductDetail.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	statusCode, err := validateProduct(productData)
	if err != nil {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(statusCode)
	fmt.Fprint(w, "Product details updated")
}

func updateProductAssets(w http.ResponseWriter, r *http.Request) {
	log.Println("Update product assets")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	productData, err := parseProductData(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to get product")
	}

	err = dbController.UpdateProductAssets(productData)
	if err != nil {
		if err.Error() == dbcontrollers.ErrMissingProductAsset.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	statusCode, err := validateProduct(productData)
	if err != nil {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(statusCode)
	fmt.Fprint(w, "Product assets updated")
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting product")
	data, err := decodePostData(w, r)
	if err != nil {
		return
	}

	productIDString, ok := data["product_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'product_id'")
		return
	}

	productID, err := uuid.Parse(productIDString.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid 'productID'")
		return
	}

	err = dbController.DeleteProduct(&productID)
	if err != nil {
		if err.Error() == dbcontrollers.ErrProductNotFound.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	_, err = dbController.GetProduct(&productID)
	if err != nil && err.Error() != dbcontrollers.ErrProductNotFound.Error() {
		err = errors.Wrap(errors.WithStack(err), "Failed to get product")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Delete completed")
}

func addProductUser(w http.ResponseWriter, r *http.Request) {
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
		if err := dbController.AddProductUser(&productID, &userID, int(privilege)); err != nil {
			if err.Error() == dbcontrollers.ErrProductNotFound.Error() ||
				err.Error() == dbcontrollers.ErrProductUserNotAssociated.Error() {
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprint(w, err.Error())
				return
			}
			err = errors.Wrap(errors.WithStack(err), "Failed to add product user")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Add product user completed")
}

func deleteProductUser(w http.ResponseWriter, r *http.Request) {
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

	if err := dbController.DeleteProductUser(&productID, &userID); err != nil {
		if err.Error() == dbcontrollers.ErrProductUserNotAssociated.Error() {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, err.Error())
			return
		}
		err = errors.Wrap(errors.WithStack(err), "Failed to delete product user")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Delete product user completed")
}

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/add-user", addUser)
	http.HandleFunc("/get-user", getUser)
	http.HandleFunc("/get-user-by-email", getUserByEmail)
	http.HandleFunc("/get-users", getUsers)
	http.HandleFunc("/update-user-settings", updateUserSettings)
	http.HandleFunc("/update-user-assets", updateUserAssets)
	http.HandleFunc("/delete-user", deleteUser)
	http.HandleFunc("/add-product-user", addProductUser)
	http.HandleFunc("/delete-product-user", deleteProductUser)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/add-product", addProduct)
	http.HandleFunc("/get-product", getProduct)
	http.HandleFunc("/get-products", getProducts)
	http.HandleFunc("/update-product-details", updateProductDetails)
	http.HandleFunc("/update-product-assets", updateProductAssets)
	http.HandleFunc("/delete-product", deleteProduct)

	controller, err := dbcontrollers.NewDBController()
	if err != nil {
		panic(err)
	}
	project := dbcontrollers.ProjectDBDummy{}
	dbcontrollers.SetProjectDB(project)
	dbController = controller

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	panic(http.ListenAndServe(":8080", nil))
}
