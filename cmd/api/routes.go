package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	// initialize a new http router instance
	router := httprouter.New()

	// Convert the error helpers to a http.Handler using http.HandlerFunc adapter and
	// set them as custom error handlers
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// register relevant endpoints and their methods
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/contacts/:id", app.showContactHandler)
	router.HandlerFunc(http.MethodPost, "/v1/contacts", app.createContactHandler)
	router.HandlerFunc(http.MethodGet, "/v1/contacts", app.listAllContactsHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/contacts/:id", app.deleteContactHandler)

	// return configured router
	return router
}
