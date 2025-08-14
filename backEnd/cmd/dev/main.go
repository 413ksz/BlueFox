// main is the entry point for the BlueFox API server when running locally.
package main

import (
	"net/http"
	"os"

	"github.com/413ksz/BlueFox/backEnd/pkg/app"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

// init is a function that is called when the application is started.
// It sets up the global variables for warm starts.
func init() {
	app.Init()
}

// main is the entry point for the BlueFox API server when running locally.
//
// Functinality:
// 1. Set up local server: configure the router and start the local server
// 2. Set up CORS middleware: add CORS middleware to the router to allow cross-origin requests from the front-end app
// 3. Start local server: start the local server and listen for incoming requests
//
// Error Conditions:
// - if the local server fails to start
func main() {
	// --- Local Server Configuration ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Warn().
			Str("component", "main_local_server").
			Str("event", "port_env_var_missing").
			Str("default_port", port).
			Msg("PORT environment variable not set, defaulting to 8080")
	}
	addr := ":" + port

	log.Info().
		Str("component", "main_local_server").
		Str("event", "local_server_start_attempt").
		Str("listen_address", addr).
		Msg("Starting local HTTP server for testing routes...")

	// --- CORS Middleware ---

	// Add CORS middleware to the router to allow cross-origin requests from the front-end app
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Wrap the router with the CORS middleware
	handlerWithCORS := c.Handler(app.AppRouter)

	// --- Local Server Start ---
	// Start the local HTTP server
	if err := http.ListenAndServe(addr, handlerWithCORS); err != nil {
		log.Fatal().
			Err(err).
			Str("component", "main_local_server").
			Str("event", "local_server_failure").
			Str("listen_address", addr).
			Msg("Local HTTP server failed to start or stopped unexpectedly")
	}
}
