package main

import (
	"broker/event"
	"bytes"         // converts JSON bytes into request body
	"encoding/json" // marshal/unmarshal JSON
	"errors"
	"net/http"
	"net/rpc"
)

// ====================================================
// STRUCTS USED FOR REQUESTS
// ====================================================

// RequestPayload represents the JSON received
// by Broker Service from Frontend.
//
// Example:
//
//	{
//	   "action":"auth",
//	   "auth":{
//	      "email":"admin@example.com",
//	      "password":"secret"
//	   }
//	}
type RequestPayload struct {

	// tells broker what operation to perform
	Action string `json:"action"`

	// nested auth payload
	Auth AuthPayload `json:"auth,omitempty"`

	Log LogPayload `json:"log,omitempty"`

	Mail MailPayload `json:"mail,omitempty"`
}

// Represents login credentials
//
//	{
//	   "email":"admin@example.com",
//	   "password":"secret"
//	}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type RPCPayload struct {
	Name string
	Data string
}

// ====================================================
// BROKER ENDPOINT
// ====================================================

func (app *Config) Broker(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ------------------------------------------------
	// This endpoint exists mainly for testing.
	// Frontend can hit broker and verify it works.
	// ------------------------------------------------

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(
		w,
		http.StatusOK,
		payload,
	)
}

// ====================================================
// HANDLE SUBMISSION
// ====================================================
//
// Receives request from frontend.
//
// Example:
//
// {
//   "action":"auth",
//   "auth":{
//      "email":"admin@example.com",
//      "password":"secret"
//   }
// }
//
// ====================================================

func (app *Config) HandleSubmission(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ------------------------------------------------
	// Create empty structure to receive JSON
	// ------------------------------------------------

	var requestPayload RequestPayload

	// ------------------------------------------------
	// Convert request JSON into Go struct
	// ------------------------------------------------

	err := app.readJSON(
		w,
		r,
		&requestPayload,
	)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// ------------------------------------------------
	// Decide what action to perform
	// ------------------------------------------------

	switch requestPayload.Action {

	// Login request
	case "auth":

		app.authenticate(
			w,
			requestPayload.Auth,
		)
	case "log":
		app.logItemViaRPC(w, requestPayload.Log)

	case "mail":

		app.sendMail(w, requestPayload.Mail)

	// Unknown action
	default:

		app.errorJSON(
			w,
			errors.New("unknown action"),
		)
	}
}

// ====================================================
// AUTHENTICATE
// ====================================================
//
// Sends request to Authentication Service
//
// Broker does NOT validate user itself.
//
// It forwards credentials to auth service.
//
// ====================================================

// ====================================================
// LOG ITEM
// ====================================================
//
// Sends log information to Logger Service.
//
// Frontend
//    |
//    ▼
// Broker Service
//    |
//    ▼
// Logger Service
//    |
//    ▼
// MongoDB
//
// ====================================================

func (app *Config) logItem(
	w http.ResponseWriter,
	entry LogPayload,
) {

	// ----------------------------------------
	// Convert LogPayload struct into JSON
	//
	// Go Struct:
	// {
	//   Name: "movie-service",
	//   Data: "something happened"
	// }
	//
	// becomes:
	//
	// {
	//   "name":"movie-service",
	//   "data":"something happened"
	// }
	// ----------------------------------------

	jsonData, _ := json.MarshalIndent(
		entry,
		"",
		"\t",
	)

	// ----------------------------------------
	// URL of Logger Service
	//
	// Docker automatically resolves:
	//
	// logger-service
	//
	// to Logger Service container
	// ----------------------------------------

	logServiceURL := "http://logger-service/log"

	// ----------------------------------------
	// Create HTTP POST request
	//
	// Broker ----> Logger Service
	// ----------------------------------------

	request, err := http.NewRequest(
		"POST",
		logServiceURL,
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Tell Logger Service we're sending JSON
	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	// ----------------------------------------
	// HTTP Client
	//
	// Responsible for sending requests
	// ----------------------------------------

	client := &http.Client{}

	// ----------------------------------------
	// Send request to Logger Service
	// ----------------------------------------

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	// ----------------------------------------
	// Verify Logger Service accepted request
	// ----------------------------------------

	if response.StatusCode != http.StatusAccepted {

		app.errorJSON(
			w,
			errors.New("error calling log service"),
		)

		return
	}

	// ----------------------------------------
	// Success Response back to Frontend
	// ----------------------------------------

	var payload jsonResponse

	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(
		w,
		http.StatusAccepted,
		payload,
	)
}

func (app *Config) authenticate(
	w http.ResponseWriter,
	a AuthPayload,
) {

	// ------------------------------------------------
	// Convert AuthPayload struct
	// into JSON bytes
	//
	// {
	//   "email":"admin@example.com",
	//   "password":"secret"
	// }
	// ------------------------------------------------

	jsonData, _ := json.MarshalIndent(
		a,
		"",
		"\t",
	)

	// ------------------------------------------------
	// Create HTTP request to Auth Service
	//
	// Docker service name:
	// authentication-service
	//
	// Docker DNS automatically resolves it.
	// ------------------------------------------------

	request, err := http.NewRequest(
		"POST",
		"http://authentication-service/authenticate",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// ------------------------------------------------
	// HTTP client
	// ------------------------------------------------

	client := &http.Client{}

	// ------------------------------------------------
	// Call Authentication Service
	// ------------------------------------------------

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// close body when function exits
	defer response.Body.Close()

	// ------------------------------------------------
	// Validate response status
	// ------------------------------------------------

	if response.StatusCode ==
		http.StatusUnauthorized {

		app.errorJSON(
			w,
			errors.New("invalid credentials"),
		)

		return

	} else if response.StatusCode !=
		http.StatusAccepted {

		app.errorJSON(
			w,
			errors.New("error calling auth service"),
		)

		return
	}

	// ------------------------------------------------
	// Read JSON returned by auth service
	// ------------------------------------------------

	var jsonFromService jsonResponse

	err = json.NewDecoder(
		response.Body,
	).Decode(&jsonFromService)

	if err != nil {

		app.errorJSON(w, err)

		return
	}

	// ------------------------------------------------
	// Build final response
	// ------------------------------------------------

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	// ------------------------------------------------
	// Send response back to frontend
	// ------------------------------------------------

	app.writeJSON(
		w,
		http.StatusAccepted,
		payload,
	)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	//call the mail service
	mailServiceURL := "http://mailer-service/send"

	//post to mail service

	request, err := http.NewRequest(
		"POST",
		mailServiceURL,
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		app.errorJSON(w, err)
		return

	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {

	// Connect to the Logger Service's RPC server
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Prepare the data to send to the RPC server
	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	// Variable to store the response from the RPC server
	var result string

	// Call Logger Service's LogInfo() method remotely
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Prepare HTTP response for the frontend
	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	// Send response back to the frontend
	app.writeJSON(w, http.StatusAccepted, payload)
}
