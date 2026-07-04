package main

import (
	"authentication/data" // Our custom package containing models

	"database/sql" // Provides generic SQL database functionality
	"fmt"          // Used for string formatting
	"log"          // Used for logging messages
	"net/http"     // Used to create and run HTTP servers
	"os"           // Used to read environment variables
	"time"         // Used for retry delays

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver (imported for side effects)
)

// Port on which the authentication service will run
const webPort = "80"

// Counts how many times we've tried connecting to PostgreSQL
var counts int64

// Config acts as a dependency container.
// Instead of passing DB and Models everywhere,
// we store them in one struct.
type Config struct {
	DB     *sql.DB     // Database connection
	Models data.Models // Database models
}

func main() {

	// Log startup message
	log.Println("Starting authentication service")

	// Try connecting to PostgreSQL
	conn := connectToDB()

	// If database connection failed after retries,
	// stop the application.
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// Create application configuration
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// Create HTTP server configuration
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),

		// routes() returns our router (mux)
		Handler: app.routes(),
	}

	// Start listening for HTTP requests
	err := srv.ListenAndServe()

	// Crash application if server fails
	if err != nil {
		log.Panic(err)
	}
}

// openDB performs ONE attempt to connect to PostgreSQL.
func openDB(dsn string) (*sql.DB, error) {

	// Open a database connection
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	// Verify that database is actually reachable
	err = db.Ping()
	//sql.Open() creates a database handle but doen't actually connect to the database.
	//  db.Ping() tries to connect to the database and returns an error if it fails.
	if err != nil {
		return nil, err
	}

	// Success
	return db, nil
}

// connectToDB keeps retrying until PostgreSQL is ready.
func connectToDB() *sql.DB {

	// Read DSN from environment variable. DSN (Data Source Name) is a string that contains the information needed to connect to the database, such as username, password, host, port, and database name.
	dsn := os.Getenv("DSN")

	for {

		// Try opening database
		connection, err := openDB(dsn)

		if err != nil {

			log.Println("Postgres not yet ready...")

			// Increment retry counter
			counts++

		} else {

			log.Println("Connected to Postgres!")

			// Success
			return connection
		}

		// Stop after too many attempts
		if counts > 10 {

			log.Println(err)

			return nil
		}

		// Wait before retrying
		log.Println("Backing off for two seconds...")

		time.Sleep(2 * time.Second)

		// Start next iteration
		continue
	}
}
