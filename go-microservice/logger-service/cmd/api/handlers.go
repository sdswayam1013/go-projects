package main

import (
	"log-service/data"
	"net/http"
)

/////////////////////////////////////////////////
// REQUEST STRUCTURE
/////////////////////////////////////////////////

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

/////////////////////////////////////////////////
// WRITE LOG HANDLER
/////////////////////////////////////////////////

func (app *Config) WriteLog(
	w http.ResponseWriter,
	r *http.Request,
) {

	// Holds incoming JSON request
	var requestPayload JSONPayload

	// Read JSON from request body
	err := app.readJSON(
		w,
		r,
		&requestPayload,
	)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Create log object
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	// Save into MongoDB
	err = app.Models.LogEntry.Insert(event)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Success response
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(
		w,
		http.StatusAccepted,
		resp,
	)
}
