package handlers

import (
	"log"
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/errors"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// UserGetHandler handles HTTP GET requests for fetching a single user by ID.
// It retrieves the user ID (UUID) from the URL path variables (e.g., /users/{id}),
// queries the database, and returns the user data as JSON.
func UserGetHandler(w http.ResponseWriter, r *http.Request) {

	// Define the context and method for the API response.
	const (
		CONTEXT        string = "api/user/"
		METHOD         string = "GET"
		STATUS_DEFAULT int    = http.StatusOK
	)

	apiResponse := &models.ApiResponse[models.User]{}
	apiResponse.Method = METHOD
	apiResponse.Context = CONTEXT
	apiResponse.StatusCode = STATUS_DEFAULT
	// Get the GORM database instance.
	db := database.DB

	// Check if the database connection is initialized.
	if db == nil {
		apiResponse.Error = errors.ERROR_CODE_DATABASE_INITIALIZE.ApiErrorResponse("Database not ready for UserGetHandler", nil)
		log.Printf("ERROR: [%s][%s] Database not initialized. Error: %s", apiResponse.Context, apiResponse.Method, apiResponse.Error.Details)
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Extract the user ID (which is a UUID string) from the URL path variables using mux.Vars.
	// We expect the route to be something like "/users/{id}".
	vars := mux.Vars(r)
	userID := vars["id"] // The key "id" comes from the route definition

	apiResponse.Params = map[string]interface{}{
		"id": userID,
	}

	if userID == "" {
		// This should ideally not happen if the route is defined correctly with {id},
		// but it's a good defensive check.
		apiResponse.Error = errors.ERROR_CODE_INVALID_INPUT.ApiErrorResponse("User ID missing from path", nil)
		log.Printf("WARN: [%s][%s] Invalid input: User ID missing from path. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
		models.SendApiResponse(w, apiResponse)
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
			apiResponse.Error = errors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("User not found", nil)
			models.SendApiResponse(w, apiResponse)
			log.Printf("INFO: [%s][%s] User not found for ID: %s", apiResponse.Context, apiResponse.Method, userID)
			return
		}
		// If there's any other error, log it and return an internal server error.
		apiResponse.Error = errors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("Error fetching user", nil)
		log.Printf("ERROR: [%s][%s] Database error fetching user ID %s: %v", apiResponse.Context, apiResponse.Method, userID, result.Error)
		models.SendApiResponse(w, apiResponse)
		return
	}

	// If the user is found, return the user data as JSON.
	apiResponse.Data = &models.ResponseData[models.User]{
		Items: []models.User{user},
	}
	log.Printf("INFO: [%s][%s] Successfully fetched user ID: %s", apiResponse.Context, apiResponse.Method, userID)
	models.SendApiResponse(w, apiResponse)
}
