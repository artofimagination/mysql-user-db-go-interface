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
	mysqldb.Functions = mysqldb.MYSQLFunctions{}
	user, err := controllers.CreateUser(name, email, []byte(password),
		func(*uuid.UUID) string {
			return "testPath"
		}, func([]byte) ([]byte, error) {
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
	mysqldb.Functions = mysqldb.MYSQLFunctions{}
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	result, err := mysqldb.Functions.GetUser(mysqldb.ByEmail, email, tx)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, result)
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
	mysqldb.Functions = mysqldb.MYSQLFunctions{}
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	user, err := mysqldb.Functions.GetUser(mysqldb.ByEmail, email, tx)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	err = mysqldb.Functions.DeleteUser(&user.ID, tx)
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
	mysqldb.Functions = mysqldb.MYSQLFunctions{}
	tx, err := mysqldb.DBConnector.ConnectSystem()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	_, err = mysqldb.Functions.GetUser(mysqldb.ByEmail, email, tx)
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

	mysqldb.DBConnection = "root:123secure@tcp(user-db:3306)/user_database?parseTime=true"
	mysqldb.MigrationDirectory = fmt.Sprintf("%s/src/mysql-user-db-go-interface/db/migrations/mysql", os.Getenv("GOPATH"))

	mysqldb.DBConnector = mysqldb.MYSQLConnector{}
	if err := mysqldb.DBConnector.BootstrapSystem(); err != nil {
		log.Fatalf("System bootstrap failed. %s", errors.WithStack(err))
	}

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	panic(http.ListenAndServe(":8080", nil))
}
