package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(
	w http.ResponseWriter,
	r *http.Request,
) {

	// Create a temporary structure
	// to hold incoming JSON data
	//
	// Incoming JSON:
	//
	// {
	//   "email":"admin@example.com",
	//   "password":"secret"
	// }
	// --------------------------------------------------
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// --------------------------------------------------
	// Read JSON from request body
	//
	// r.Body contains raw JSON
	//
	// Convert:
	//
	// {
	//   "email":"admin@example.com"
	// }
	//
	// into:
	//
	// requestPayload.Email
	// --------------------------------------------------
	err := app.readJSON(
		w,
		r,
		&requestPayload,
	)

	// If JSON is invalid
	if err != nil {

		// Send error response
		app.errorJSON(
			w,
			err,
			http.StatusBadRequest,
		)

		return
	}

	// --------------------------------------------------
	// Find user in PostgreSQL
	//
	// SELECT * FROM users
	// WHERE email = ?
	// --------------------------------------------------
	user, err := app.Models.User.GetByEmail(
		requestPayload.Email,
	)

	// User not found
	if err != nil {

		app.errorJSON(
			w,
			errors.New("invalid credentials"),
			http.StatusBadRequest,
		)

		return
	}

	// --------------------------------------------------
	// Compare entered password
	// with BCrypt hash stored in DB
	// --------------------------------------------------
	valid, err := user.PasswordMatches(
		requestPayload.Password,
	)

	// Password mismatch
	if err != nil || !valid {

		app.errorJSON(
			w,
			errors.New("invalid credentials"),
			http.StatusBadRequest,
		)

		return
	}
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// --------------------------------------------------
	// Build success response
	// --------------------------------------------------
	payload := jsonResponse{

		// Login successful
		Error: false,

		// Success message
		Message: fmt.Sprintf(
			"Logged in user %s",
			user.Email,
		),

		// User information
		Data: user,
	}

	// --------------------------------------------------
	// Send JSON response
	// --------------------------------------------------
	app.writeJSON(
		w,
		http.StatusAccepted,
		payload,
	)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
