package main

import (
	"net/http"
)

/*
Broker is an HTTP handler method attached to Config.
*/
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	// ===== CORS HEADERS (VERY IMPORTANT) =====
	// These allow your frontend (port 80) to talk to backend (port 8080)

	w.Header().Set("Access-Control-Allow-Origin", "*")                   // allow all origins
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // allow headers
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // allowed methods

	// Handle preflight request (browser sends OPTIONS before POST)
	if r.Method == "OPTIONS" {
		return
	}

	// ===== YOUR EXISTING LOGIC =====

	// Create response payload (Go struct)
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

	// // Convert Go struct → JSON
	// out, err := json.MarshalIndent(payload, "", "\t")
	// if err != nil {
	// 	http.Error(w, "Error creating JSON", http.StatusInternalServerError)
	// 	return
	// }

	// // Tell client this is JSON
	// w.Header().Set("Content-Type", "application/json")

	// // Set HTTP status code
	// w.WriteHeader(http.StatusAccepted)

	// // Send JSON response back to client
	//w.Write(out)
}
