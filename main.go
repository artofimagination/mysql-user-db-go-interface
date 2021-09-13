package main

import (
	"fmt"
	"net/http"

	"github.com/artofimagination/mysql-user-db-go-interface/initialization"
	"github.com/artofimagination/mysql-user-db-go-interface/restcontrollers"
)

func main() {
	cfg := &initialization.Config{}
	initialization.InitConfig(cfg)
	_, err := restcontrollers.NewRESTController()
	if err != nil {
		panic(err)
	}

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	port := fmt.Sprintf(":%d", cfg.Port)
	panic(http.ListenAndServe(port, nil))
}
