package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"

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

func addUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding user")
	if err := checkRequestType(POST, w, r); err != nil {
		return
	}

	data := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(errors.WithStack(err), "Failed to decode request json")
		fmt.Fprint(w, err.Error())
		return
	}

	name, ok := data["name"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'name'")
		return
	}

	email, ok := data["email"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'email'")
		return
	}

	password, ok := data["password"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'password'")
		return
	}

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

func deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting user")
	if err := checkRequestType(POST, w, r); err != nil {
		return
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(errors.WithStack(err), "Failed to decode request json")
		fmt.Fprint(w, err.Error())
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

func authenticateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Authenticate user")
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]

	passwords, ok := r.URL.Query()["password"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'password' is missing")
		return
	}

	password := []byte(passwords[0])

	err := dbController.Authenticate(email, password, func(string, []byte, *models.User) error { return nil })
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintln(w, "User authenticated")
}

func updateUserSettings(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user settings")

	userData, err := queryUser(r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	settingsList, ok := r.URL.Query()["settings"]
	if !ok || len(settingsList[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'settings' is missing")
		return
	}

	err = dbController.UpdateUserSettings(&userData.Settings)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintln(w, "User settings updated")
}

func updateUserAssets(w http.ResponseWriter, r *http.Request) {
	log.Println("Update user assets")

	userData, err := queryUser(r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	assetList, ok := r.URL.Query()["assets"]
	if !ok || len(assetList[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'assets' is missing")
		return
	}

	err = dbController.UpdateUserAssets(&userData.Assets)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintln(w, "User assets updated")
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding product")
	if err := checkRequestType(POST, w, r); err != nil {
		return
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(errors.WithStack(err), "Failed to decode request json")
		fmt.Fprint(w, err.Error())
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

func updateProductDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("Update product details")

	productData, err := queryProduct(r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	detailsList, ok := r.URL.Query()["details"]
	if !ok || len(detailsList[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'details' is missing")
		return
	}

	err = dbController.UpdateProductDetails(&productData.Details)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintln(w, "Product details updated")
}

func updateProductAssets(w http.ResponseWriter, r *http.Request) {
	log.Println("Update product assets")

	productData, err := queryProduct(r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	assetsQuery, ok := r.URL.Query()["assets"]
	if !ok || len(assetsQuery[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'assets' is missing")
		return
	}

	err = dbController.UpdateProductAssets(&productData.Assets)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintln(w, "Product assets updated")
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting product")
	productData, err := queryProduct(r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	err = dbController.DeleteProduct(&productData.ID)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	fmt.Fprintln(w, "Product delete completed")
}

func addProductUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding product user")
	if err := checkRequestType(POST, w, r); err != nil {
		return
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(errors.WithStack(err), "Failed to decode request json")
		fmt.Fprint(w, err.Error())
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
	if err := checkRequestType(POST, w, r); err != nil {
		return
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(errors.WithStack(err), "Failed to decode request json")
		fmt.Fprint(w, err.Error())
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
	http.HandleFunc("/update-user-settings", updateUserSettings)
	http.HandleFunc("/update-user-assets", updateUserAssets)
	http.HandleFunc("/delete-user", deleteUser)
	http.HandleFunc("/add-product-user", addProductUser)
	http.HandleFunc("/delete-product-user", deleteProductUser)
	http.HandleFunc("/authenticate-user", authenticateUser)
	http.HandleFunc("/add-product", addProduct)
	http.HandleFunc("/get-product", getProduct)
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

	mysqldb.DBConnection = fmt.Sprintf("%s:%s@tcp(user-db:3306)/%s?parseTime=true", os.Getenv("MYSQL_DB_USER"), os.Getenv("MYSQL_DB_PASSWORD"), os.Getenv("MYSQL_DB_NAME"))
	mysqldb.MigrationDirectory = fmt.Sprintf("%s/src/mysql-user-db-go-interface/db/migrations/mysql", os.Getenv("GOPATH"))

	mysqldb.DBConnector = mysqldb.MYSQLConnector{}
	if err := mysqldb.DBConnector.BootstrapSystem(); err != nil {
		log.Fatalf("System bootstrap failed. %s", errors.WithStack(err))
	}

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	panic(http.ListenAndServe(":8080", nil))
}
