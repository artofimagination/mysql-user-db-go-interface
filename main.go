package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/artofimagination/mysql-user-db-go-interface/controllers"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi! I am Server!")
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Inserting user")
	names, ok := r.URL.Query()["name"]
	if !ok || len(names[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'name' is missing")
		return
	}

	name := names[0]
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

	password := passwords[0]
	models.Interface = models.RepoInterface{}
	mysqldb.FunctionInterface = mysqldb.MYSQLFunctionInterface{}
	user, err := controllers.CreateUser(name, email, password,
		func(*uuid.UUID) string {
			return "testPath"
		}, func(string) ([]byte, error) {
			return []byte{}, nil
		})
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintln(w, user)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Getting user")
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]
	mysqldb.FunctionInterface = mysqldb.MYSQLFunctionInterface{}
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	result, err := mysqldb.FunctionInterface.GetUserByEmail(email, tx)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, result)
	}
}

func getSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Getting user settings")
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]
	mysqldb.FunctionInterface = mysqldb.MYSQLFunctionInterface{}
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	resultUser, err := mysqldb.FunctionInterface.GetUserByEmail(email, tx)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, resultUser)
	}

	resultSettings, err := mysqldb.GetSettings(&resultUser.SettingsID)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, *resultSettings)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Deleting user")
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]
	mysqldb.FunctionInterface = mysqldb.MYSQLFunctionInterface{}
	err := mysqldb.DeleteUser(email)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
}

func deleteSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Deleting user settings")
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]
	mysqldb.FunctionInterface = mysqldb.MYSQLFunctionInterface{}
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	resultUser, err := mysqldb.FunctionInterface.GetUserByEmail(email, tx)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, resultUser)
	}

	err = mysqldb.DeleteSettings(&resultUser.SettingsID)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
}

func checkUserPass(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Check user pass")
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]
	mysqldb.FunctionInterface = mysqldb.MYSQLFunctionInterface{}
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	_, err = mysqldb.FunctionInterface.GetUserByEmail(email, tx)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, "User pass matched!")
	}
}

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/insert", insertUser)
	http.HandleFunc("/get", getUser)
	http.HandleFunc("/delete", deleteUser)
	http.HandleFunc("/check", checkUserPass)
	http.HandleFunc("/get-settings", getSettings)
	http.HandleFunc("/delete-settings", deleteSettings)

	mysqldb.DBConnection = "root:123secure@tcp(user-db:3306)/user_database?parseTime=true"
	mysqldb.MigrationDirectory = fmt.Sprintf("%s/src/mysql-user-db-go-interface/db/migrations/mysql", os.Getenv("GOPATH"))

	mysqldb.DBConnector = mysqldb.MYSQLConnector{}
	if err := mysqldb.DBConnector.BootstrapSystem(); err != nil {
		log.Fatalf("System bootstrap failed. %s", errors.WithStack(err))
	}

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	panic(http.ListenAndServe(":8080", nil))
}
