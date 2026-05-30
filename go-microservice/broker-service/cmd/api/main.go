package main

import (
	"fmt"      // used for string formatting (e.g., building ":80")
	"log"      // used for logging messages to console
	"net/http" // core HTTP server package in Go
)

// webPort defines the port your server will run on
// Note: port 80 may require sudo/root privileges on Linux/macOS
const webPort = "80"

// Config is your application struct
// You attach methods (routes, handlers) to this struct
// This helps in organizing code and injecting dependencies later (DB, config, etc.)
type Config struct{} //conceptually it represents my app

// main is the entry point of your application
func main() {

	// I am creating my application object(instance of Config struct))
	app := Config{} //it carries my methods like routes() and btroker()

	// Log that your service is starting
	log.Printf("Starting broker service on port %s\n", webPort)

	// Create a new HTTP server i.e build a machine that listens for HTTP requests
	// http.Server is a struct from net/http that represents your server
	srv := &http.Server{

		// Addr defines which port the server listens on
		// ":80" means listen on port 80 on all interfaces
		Addr: fmt.Sprintf(":%s", webPort),

		// Handler is VERY IMPORTANT
		// when a request comes in, use my routing logic to decide what to do
		// app.routes() will return a router (like ServeMux)
		Handler: app.routes(), //routes() decides which function to call. It connects main.go to handlers.go
	}

	// Start the server
	// ListenAndServe:
	// - opens the port
	// - starts listening for incoming HTTP requests
	// - blocks the program (keeps it running)
	err := srv.ListenAndServe()

	// If server fails (port busy, permission issue, etc.), log and crash
	if err != nil {
		log.Panic(err)
	}
}

/*
routes() defines how URLs map to handler functions.

It returns something that implements http.Handler.
In this case, we use http.ServeMux (default router in Go).
Defines a method called routes,
*/
func (app *Config) routes() http.Handler { //

	// Creates a router called mux(multiplexer) that matches incoming requests to their handlers)
	mux := http.NewServeMux()

	// Define route:
	// When someone hits "/broker", call app.Broker
	mux.HandleFunc("/broker", app.Broker)

	// Return the router to the server
	return mux
}

/*
jsonResponse defines the structure of the JSON response.

The struct tags (`json:"..."`) control how fields appear in JSON.
*/

/*
Broker is an HTTP handler function.

It matches the required signature:
func(w http.ResponseWriter, r *http.Request)
*/
// Broker handler is defined in handlers.go
