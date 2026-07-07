package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
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

	//register the rpc server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

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

func (app *Config) rpcListen() error {

	// Start the RPC server
	log.Println("Starting RPC server on port", rpcPort)

	// Listen for incoming TCP connections on the RPC port
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}

	// Close the listener when the application shuts down
	defer listen.Close()

	// Keep accepting new client connections forever
	for {

		// Wait for a client (Broker) to connect
		rpcConn, err := listen.Accept()
		if err != nil {
			continue // Ignore failed connections and keep listening
		}

		// Handle each RPC connection in a separate goroutine
		go rpc.ServeConn(rpcConn)
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
