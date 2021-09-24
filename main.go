package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artofimagination/mysql-user-db-go-interface/initialization"
	"github.com/artofimagination/mysql-user-db-go-interface/restcontrollers"
)

func main() {
	cfg := &initialization.Config{}
	initialization.InitConfig(cfg)
	r, err := restcontrollers.NewRESTController()
	if err != nil {
		panic(err)
	}

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	port := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Handler:      r,
		Addr:         port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Shutting down")
	os.Exit(0)
}
