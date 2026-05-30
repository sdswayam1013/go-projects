package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"3
	"github.com/go-chi/cors"
)

/* chi - A better router than Go’s default ServeMux and offers More features,
cleaner routing, middleware support and more production friendly.
chi middleware - Provides ready-made middlewares (logging, heartbeat, etc.)
cors package go-chi- Handles CORS (Cross-Origin Resource Sharing) */

func routes() http.Handler { /* routes() is a method that returns an HTTP handler */
	mux := chi.NewRouter() /* Create a new router using chi.NewRouter() and mux is an upgraded traffic controller*/

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

	mux.Use(middleware.Heartbeat("/ping")) /*creates a built-in route
	  so when GET /ping -> returns "OK" */
	return mux
}
