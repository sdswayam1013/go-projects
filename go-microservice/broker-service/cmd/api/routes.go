package main

import (
	"net/http"

	"github.com/rs/cors"
)

func (app *Config) routesMux() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/broker", app.Broker)

	mux.HandleFunc("/handle", app.HandleSubmission)

	mux.HandleFunc("/log-grpc", app.LogViaGRPC)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	return c.Handler(mux)
}
