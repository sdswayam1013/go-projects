package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// =========================
// Common JSON Response Struct
// =========================

// This struct defines a standard response format for your API.
// Why: Keeps all responses consistent across endpoints.
type jsonResponse struct {
	Error   bool        `json:"error"`          // Indicates success or failure
	Message string      `json:"message"`        // Human-readable message
	Data    interface{} `json:"data,omitempty"` // Optional payload (omitted if empty)
}

// =========================
// READ JSON (Request → Go struct)
// =========================

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {

	// Limit request body size to 1MB
	// Why: Prevent malicious users from sending huge payloads and crashing server
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Create JSON decoder to read request body
	// Why: Efficient for streaming data from HTTP request
	dec := json.NewDecoder(r.Body)

	// Decode JSON into provided struct (data)
	// Why: Convert raw JSON → usable Go struct
	err := dec.Decode(data)
	if err != nil {
		return err // return error if JSON is invalid
	}

	// Check if there is more than one JSON object in request
	// Why: Prevent malformed requests like `{...}{...}`
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only one JSON value")
	}

	return nil
}

// =========================
// WRITE JSON (Go struct → Response)
// =========================

func (app *Config) writeJSON(
	w http.ResponseWriter,
	status int,
	data interface{},
	headers ...http.Header, // optional headers
) error {

	// Convert Go struct → JSON bytes
	// Why: HTTP responses must be sent as bytes, not Go structs
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// If optional headers are provided, add them to response
	// Why: Allows flexibility (e.g., custom headers, auth, etc.)
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	// Set content type to JSON
	// Why: Tells client how to interpret the response body
	w.Header().Set("Content-Type", "application/json")

	// Set HTTP status code (e.g., 200, 400, 500)
	// Why: Communicates result of request clearly
	w.WriteHeader(status)

	// Write JSON response body
	// Why: Sends actual data to client
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// =========================
// ERROR JSON (Standardized error response)
// =========================

func (app *Config) errorJSON(
	w http.ResponseWriter,
	err error,
	status ...int, // optional status code
) error {

	// Default status code = 400 (Bad Request)
	// Why: Most errors are client-side issues
	statusCode := http.StatusBadRequest

	// If custom status is provided, override default
	// Why: Allows flexibility (e.g., 500, 401, etc.)
	if len(status) > 0 {
		statusCode = status[0]
	}

	// Create standardized error payload
	// Why: Keeps all error responses consistent
	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	// Reuse writeJSON to send response
	// Why: Avoid duplication and keep response handling centralized
	return app.writeJSON(w, statusCode, payload)
}
