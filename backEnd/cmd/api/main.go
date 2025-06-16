// Handler is the primary entry point for the serverless function hosted on Vercel
// Rename this package to main to run run the server locally
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/router"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// AppRouter is the main router for handling API requests. It's initialized
// in the init() function and used by both the Vercel serverless handler and local main() server.
var appRouter *mux.Router

// init initializes the serverless function environment.
// This function is automatically executed by the Go runtime before main() is called
// (in local development) or before the Vercel Handler is invoked.
// It sets up global configurations like logging, database connection, and API routes.
// Any unrecoverable errors during initialization will cause the application to terminate.
func init() {
	// --- Zerolog Configuration ---
	// Set the global log level for Zerolog
	// zerolog.SetGlobalLevel(zerolog.InfoLevel) for production
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Configure Zerolog's global logger to output JSON to os.Stdout
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	// If DEBUG environment variable is set to "true", set the log level to zerolog.DebugLevel
	// and enable console logging for better readability
	if os.Getenv("DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().Timestamp().Caller().Logger()
		log.Debug().
			Str("component", "main_app").
			Str("event", "debug_mode_enabled").
			Msg("Debug mode enabled, logging to console and setting debug level")
	}

	log.Info().
		Str("component", "main_app").
		Str("event", "app_init_start").
		Msg("Serverless function initializing...")

	// --- Database Configuration ---

	// Initialize the global database connection
	if err := database.InitAppDB(); err != nil {
		log.Fatal().
			Err(err).
			Str("component", "main_app").
			Str("event", "app_db_init_failure").
			Msg("Failed to initialize global database")
	}

	log.Info().
		Str("component", "main_app").
		Str("event", "app_db_init_success").
		Msg("Global database connection successfully initialized.")

	// --- API Routes ---
	// Initialize the API router
	appRouter = mux.NewRouter()
	// Register the API routes
	router.RegisterRoutes(appRouter)

	log.Info().
		Str("component", "main_app").
		Str("event", "api_routes_initialized").
		Msg("API routes initialized.")
}

// Handler is the primary entry point for the serverless function hosted on serverless environment.
func Handler(w http.ResponseWriter, r *http.Request) {
	// --- HTTP Request Handler ---
	log.Info().
		Str("component", "main_app_handler").
		Str("event", "http_request_received").
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remote_addr", r.RemoteAddr).
		Msg("Incoming HTTP request")

	// Serve the request to the router
	appRouter.ServeHTTP(w, r)
}

// main is the entry point for the BlueFox API server when running locally.
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
	handlerWithCORS := c.Handler(appRouter)

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
