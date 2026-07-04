package main

import (
	"net/http"

	"github.com/rs/cors"
)

func (app *Config) routesMux() http.Handler {

	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	mux.HandleFunc("/send", app.SendMail)

	return c.Handler(mux)
}
