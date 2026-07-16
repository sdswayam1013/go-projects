package main

import (
	"context"          // Used for MongoDB operations
	"log"              // For logging errors/messages
	"log-service/data" // Contains the LogEntry struct
	"time"             // To store CreatedAt timestamps
	// MongoDB driver
)

// RPCServer is an empty struct.
//
// Why do we need this?
// The net/rpc package requires methods to belong to a type.
// So we create a receiver type called RPCServer.
type RPCServer struct{}

// RPCPayload represents the data coming from the Broker Service.
type RPCPayload struct {
	Name string
	Data string
}

// LogInfo is the RPC method that the Broker will call remotely.
//
// Signature required by net/rpc:
//
// func (t *T) MethodName(args T1, reply *T2) error
//
// Here:
//
// receiver      -> *RPCServer
// arguments     -> RPCPayload
// response       -> *string
// return value   -> error
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {

	//
	// Get MongoDB collection.
	//
	collection := client.Database("logs").Collection("logs")

	//
	// Insert log into MongoDB.
	//
	_, err := collection.InsertOne(
		context.TODO(),
		data.LogEntry{
			Name:      payload.Name,
			Data:      payload.Data,
			CreatedAt: time.Now(),
		},
	)

	//
	// If Mongo insertion fails,
	// log the error and return it.
	//
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	//
	// Send response back to Broker.
	//
	*resp = "Processed payload via RPC: " + payload.Name

	return nil
}
