package user

import (
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/internal/apierrors"
	"github.com/413ksz/BlueFox/backEnd/internal/database"
	"github.com/413ksz/BlueFox/backEnd/internal/model"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// UserGetHandler handles HTTP GET requests for fetching a single user by ID.
// It retrieves the user ID (UUID) from the URL path variables (e.g., /users/{id}),
// queries the database, and returns the user data as JSON.
func UserGetHandler(w http.ResponseWriter, r *http.Request) {

	// Define the context and method for the API response.
	const (
		COMPONENT      string = "user_handler"
		METHOD_NAME    string = "UserGetHandler"
		CONTEXT        string = "api/user/"
		METHOD         string = "GET"
		STATUS_DEFAULT int    = http.StatusOK
	)

	apiResponse := &models.ApiResponse[model.User]{}
	apiResponse.Method = METHOD
	apiResponse.Context = CONTEXT
	apiResponse.StatusCode = STATUS_DEFAULT
	// Get the GORM database instance.
	db := database.DB

	log.Info().
		Str("component", COMPONENT).
		Str("method_name", METHOD_NAME).
		Str("http_method", METHOD).
		Str("path", CONTEXT).
		Str("event", "http_request_received").
		Msg("Processing user retrieval request.")

	// Check if the database connection is initialized.
	if db == nil {
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_INITIALIZE.ApiErrorResponse("Database not ready for UserGetHandler", nil)
		log.Error().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "db_not_initialized").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Msg("Database not initialized for getting user data")
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
		apiResponse.Error = apierrors.ERROR_CODE_INVALID_INPUT.ApiErrorResponse("User ID missing from path", nil)
		log.Error().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "invalid_id").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Msg("User ID missing from path")
		models.SendApiResponse(w, apiResponse)
		return
	}

	var user model.User // Declare a variable to hold the fetched user data

	// Query the database for the user with the specified ID (UUID).
	// GORM will populate the 'user' variable if a record is found.
	result := db.First(&user, "id = ?", userID)

	// Check for errors returned by the database query.
	if result.Error != nil {
		// If gorm.ErrRecordNotFound is returned, it means no user with that ID was found.
		if result.Error == gorm.ErrRecordNotFound {
			apiResponse.Error = apierrors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("User not found", nil)
			log.Error().
				Str("component", COMPONENT).
				Str("method_name", METHOD_NAME).
				Str("event", "user_not_found").
				Str("api_error_code", apiResponse.Error.Code).
				Str("api_error_message", apiResponse.Error.Message).
				Int("api_error_status", apiResponse.Error.HTTPStatusCode).
				Str("id", userID).
				Err(result.Error).
				Msg("User not found")
			models.SendApiResponse(w, apiResponse)
			return
		}
		// If there's any other error, log it and return an internal server error.
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("Error fetching user", nil)
		log.Error().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "database_error").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Err(result.Error).
			Msg("Error fetching user")
		models.SendApiResponse(w, apiResponse)
		return
	}

	// If the user is found, return the user data as JSON.
	apiResponse.Data = &models.ResponseData[model.User]{
		Items: []model.User{user},
	}

	log.Info().
		Str("component", COMPONENT).
		Str("method_name", METHOD_NAME).
		Str("event", "user_fetched").
		Str("user_id", userID).
		Msg("User fetched successfully")

	models.SendApiResponse(w, apiResponse)
}
