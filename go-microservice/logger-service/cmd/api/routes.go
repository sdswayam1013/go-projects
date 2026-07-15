package main

import (
	"net/http"
)

func (app *Config) routesMux() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/log", app.WriteLog)
	// simple CORS middleware (replaces external github.com/rs/cors)
	c := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			h.ServeHTTP(w, r)
		})
	}

	return c(mux)
}
