package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(
	ctx context.Context,
	req *logs.LogRequest,
) (*logs.LogResponse, error) {

	// Extract the LogEntry message from the request
	input := req.GetLogEntry()

	// Convert protobuf object to our MongoDB model
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// Save the log to MongoDB
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		// Return failure response
		res := &logs.LogResponse{
			Result: "failed",
		}

		return res, err
	}

	// Return success response
	res := &logs.LogResponse{
		Result: "logged!",
	}

	return res, nil
}

func (app *Config) gRPCListen() {

	// Start listening for incoming gRPC connections on the configured port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	// Create a new gRPC server instance
	s := grpc.NewServer()

	// Register the LogService implementation with the gRPC server
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	// Log that the gRPC server has started successfully
	log.Printf("gRPC Server started on port %s", gRpcPort)

	// Start serving incoming gRPC requests
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
