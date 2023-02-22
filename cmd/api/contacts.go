package main

import (
	"errors"
	"fmt"
	"net/http"
	"salestrekker_technical_interview.veljkoilic/internal/data"
	"salestrekker_technical_interview.veljkoilic/internal/validator"
)

func (app *application) createContactHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Telephone string `json:"telephone"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	contact := &data.Contact{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Telephone: input.Telephone,
	}

	// check if any validation errors have been found
	v := validator.New()
	if data.ValidateContact(v, contact); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.contactsModel.InsertContact(contact)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Include location header, so user can access the created contact
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/contacts/%d", contact.ID))

	// Write a JSON response with a 201 Created status code
	err = app.writeJSON(w, http.StatusCreated, envelope{"contact": contact}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Envelope all contacts and send them to the user
func (app *application) listAllContactsHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"contacts": app.contactsModel.Contacts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showContactHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	contact, err := app.contactsModel.GetContact(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"contact": contact}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteContactHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.contactsModel.DeleteContact(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "contact successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
