package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

// Global MongoDB client
var client *mongo.Client

// Application configuration
type Config struct {
	Models data.Models
}

func main() {

	// Connect to MongoDB
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	// Save mongo connection globally
	client = mongoClient

	// Context used when disconnecting
	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	// Disconnect Mongo when application stops
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	/////////////////////////////////////////////////
	// NEW PART
	/////////////////////////////////////////////////

	app := Config{
		Models: data.New(client),
	}

	log.Println("Starting logger service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routesMux(),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

/////////////////////////////////////////////////
// HTTP SERVER
/////////////////////////////////////////////////

// func (app *Config) serve() {

// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%s", webPort),
// 		Handler: app.routesMux(),
// 	}

// 	err := srv.ListenAndServe()

// 	if err != nil {
// 		log.Panic(err)
// 	}
// }

/////////////////////////////////////////////////
// MONGODB CONNECTION
/////////////////////////////////////////////////

func connectToMongo() (*mongo.Client, error) {

	clientOptions := options.Client().
		ApplyURI(mongoURL)

	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(
		context.TODO(),
		clientOptions,
	)

	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return c, nil
}
