package main

import (
	"log"
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/router"
	"github.com/gorilla/mux"
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
