package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"salestrekker_technical_interview.veljkoilic/internal/data"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config        config
	contactsModel data.ContactsModel
	logger        *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API Server Point")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Create and initialize contactsModel
	// If initialization fails, we log it and exit the app
	contactsModel := data.NewModel()
	err := contactsModel.GetAllContacts()
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		config:        cfg,
		logger:        logger,
		contactsModel: contactsModel,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Starting HTTP server
	logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
