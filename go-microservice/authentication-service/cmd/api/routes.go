package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter() //creates a router using chi
	//sp[ecify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{ /*Adds a middleware to your router, Middleware = code that runs before your handler*/
		AllowedOrigins:   []string{"https://*", "http://*"},                                   /* Allws request from any domain */
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                 /*which http methods are allowed */
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, /*which headers can be sent by the client */
		ExposedHeaders:   []string{"Link"},                                                    /*headers visible to frontend */
		AllowCredentials: true,                                                                /*allow cookies, authorization headers, etc. */
		MaxAge:           300,                                                                 // Browser caches this CORS config for 300 seconds
		/*let browsers safely talk to my backend API */
	}))
	return mux
}
