package main

import (
	"log"
	"net/http"
	"os" // Import the os package to read environment variables

	"github.com/413ksz/BlueFox/backEnd/internal/database"
	"github.com/413ksz/BlueFox/backEnd/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // For loading .env file
	"gorm.io/gorm"
	// Adjust these imports to match your actual module path
)

var db *gorm.DB

func init() {
	var err error
	db, err = database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
}

func main() {
	// Load environment variables from .env file (for local development)
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("No .env file found or error loading: %v. Assuming environment variables are set directly.", err)
	}

	// Get the API base route from environment variable
	apiBaseRoute := os.Getenv("API_BASE_ROUTE")
	if apiBaseRoute == "" {
		apiBaseRoute = "/api"
		log.Printf("API_BASE_ROUTE environment variable not set. Using default: %s", apiBaseRoute)
	} else {
		log.Printf("API base route: %s", apiBaseRoute)
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting underlying SQL DB: %v", err)
		} else {
			sqlDB.Close()
		}
	}()

	// Setup router
	r := mux.NewRouter()

	// Create a subrouter for the base route
	apiRouter := r.PathPrefix(apiBaseRoute).Subrouter()

	// Define API routes on the subrouter
	apiRouter.HandleFunc("/test", handlers.TestHandler).Methods("GET")

	// Start server
	port := os.Getenv("PORT") // Get port from environment variable, default to 8080
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
