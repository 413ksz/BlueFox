package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// UserGetHandler handles HTTP GET requests for fetching a single user by ID.
// It retrieves the user ID (UUID) from the URL path variables (e.g., /users/{id}),
// queries the database, and returns the user data as JSON.
func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	// Get the GORM database instance.
	db := database.DB

	// Check if the database connection is initialized.
	if db == nil {
		log.Println("Database connection is not initialized in UserGetHandler.")
		http.Error(w, "Internal server error: Database not ready", http.StatusInternalServerError)
		return
	}

	// Extract the user ID (which is a UUID string) from the URL path variables using mux.Vars.
	// We expect the route to be something like "/users/{id}".
	vars := mux.Vars(r)
	userID := vars["id"] // The key "id" comes from the route definition

	if userID == "" {
		// This should ideally not happen if the route is defined correctly with {id},
		// but it's a good defensive check.
		http.Error(w, "Bad request: User ID missing from path", http.StatusBadRequest)
		return
	}

	var user models.User // Declare a variable to hold the fetched user data

	// Query the database for the user with the specified ID (UUID).
	// GORM will populate the 'user' variable if a record is found.
	result := db.First(&user, "id = ?", userID)

	// Check for errors returned by the database query.
	if result.Error != nil {
		// If gorm.ErrRecordNotFound is returned, it means no user with that ID was found.
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("User with ID %s not found", userID)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		log.Printf("Database error fetching user with ID %s: %v", userID, result.Error)
		http.Error(w, "Internal server error: Database query failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(user)

	if err != nil {
		// If there's an error during JSON encoding, log it and return an internal server error.
		log.Printf("Error encoding JSON response for user ID %s: %v", userID, err)
		http.Error(w, "Internal server error during response encoding", http.StatusInternalServerError)
		return
	}
}
