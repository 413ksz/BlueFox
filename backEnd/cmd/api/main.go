// Package handler is the entry point for the serverless function hosted on serverless invironment on Vercel.
// if you want to run locally, you can run go run main.go and change the package to main
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/router"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var appRouter *mux.Router

// init initializes the serverless function environment.
// This function is automatically executed by the Go runtime before main() is called.
// It sets up the global database connection and registers all API routes using gorilla/mux.
// Any failure during initialization is considered fatal and will cause the application to exit.
func init() {
	log.Println("Serverless function initializing...")

	if err := database.InitAppDB(); err != nil {
		log.Fatalf("Failed to initialize global database: %v", err)
	}
	log.Println("Global database connection successfully initialized.")

	appRouter = mux.NewRouter()

	router.RegisterRoutes(appRouter)
	log.Println("API routes initialized.")
}

// Handler is the primary entry point for the serverless hosted on serverless invironment.
// It handles incoming HTTP requests and dispatches them to the appropriate API handlers registered within 'appRouter'.
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: Method=%s, Path=%s", r.Method, r.URL.Path)
	appRouter.ServeHTTP(w, r)
}

// main is the entry point for the BlueFox API server when running locally.
// In the serverless environment, this function is not executed it is for local development testing.
// It starts a local HTTP server for testing purposes by listening on a port
// specified by the PORT environment variable (defaulting to 8080 if not set).
// A fatal error is logged if the server fails to start.
func main() {
	// Get the port from environment variable, default to 8080 if not set.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	log.Printf("Starting local HTTP server on %s for testing routes...", addr)

	// Configure CORS options for local development.
	// You MUST replace "http://localhost:3000" with the actual origin of your frontend development server.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}, // Allow all common HTTP methods
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},          // Allow common headers
		AllowCredentials: true,                                                         // Allow cookies, auth tokens, etc.
		Debug:            true,
	})

	handlerWithCORS := c.Handler(appRouter)

	// Use http.ListenAndServe to start the server.
	// The appRouter will handle all incoming requests.
	// Log a fatal error if the server fails to start.
	log.Fatal(http.ListenAndServe(addr, handlerWithCORS))
}
