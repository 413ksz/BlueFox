package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	// Required for time.Time in DTO
	// Required for uuid.UUID in DTO

	"github.com/413ksz/BlueFox/backEnd/pkg/apierrors"
	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	passwordHashing "github.com/413ksz/BlueFox/backEnd/pkg/password_hashing"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn" // For PostgreSQL specific errors
	"gorm.io/gorm"
)

// UserUpdateHandler handles HTTP PATCH requests for updating an existing user.
// It expects a JSON request body containing fields to update and the user ID
// in the URL path. It performs validation, hashes passwords, and updates
// the user in the database, returning the full updated user object.
func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	const (
		CONTEXT        string = "api/user/"
		METHOD         string = "PATCH"
		STATUS_DEFAULT int    = http.StatusOK
	)

	// Initialize the API response structure
	apiResponse := &models.ApiResponse[models.User]{}
	apiResponse.Method = METHOD
	apiResponse.Context = CONTEXT
	apiResponse.StatusCode = STATUS_DEFAULT

	// Get the GORM database instance
	db := database.DB

	// Check if the database connection is initialized.
	if db == nil {
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_INITIALIZE.ApiErrorResponse("Database not ready for UserUpdateHandler", nil)
		log.Printf("ERROR: [%s][%s] Database not initialized. Error: %s", apiResponse.Context, apiResponse.Method, apiResponse.Error.Details)
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Extract the user ID (UUID) from the URL path variables.
	vars := mux.Vars(r)
	userID := vars["id"]

	// Add the user ID to response parameters for logging/context
	apiResponse.Params = map[string]interface{}{
		"id": userID,
	}

	// Validate that the user ID is present in the path.
	if userID == "" {
		apiResponse.Error = apierrors.ERROR_CODE_INVALID_INPUT.ApiErrorResponse("User ID missing from path", nil)
		log.Printf("WARN: [%s][%s] Invalid input: User ID missing from path. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
		models.SendApiResponse(w, apiResponse)
		return
	}

	var existingUser models.User
	// Fetch the existing user from the database to ensure it exists and for GORM's Model context.
	result := db.First(&existingUser, "id = ?", userID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// User not found is a client-side error (404 Not Found)
			apiResponse.Error = apierrors.ERROR_CODE_NOT_FOUND.ApiErrorResponse("User not found", nil)
			log.Printf("INFO: [%s][%s] User not found for ID: %s", apiResponse.Context, apiResponse.Method, userID)
			models.SendApiResponse(w, apiResponse)
			return
		}
		// Other database errors are internal server errors
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("Error fetching user for update", nil)
		log.Printf("ERROR: [%s][%s] Database error fetching user ID %s: %v", apiResponse.Context, apiResponse.Method, userID, result.Error)
		models.SendApiResponse(w, apiResponse)
		return
	}

	var updates models.User
	// Decode the JSON request body into the models.User struct.
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		apiResponse.Error = apierrors.ERROR_CODE_ENCODE_ERROR.ApiErrorResponse("Invalid JSON data for update", nil)
		log.Printf("ERROR: [%s][%s] Error decoding request body for user update: %v", apiResponse.Context, apiResponse.Method, err)
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Prepare a map of fields to update for GORM.
	// This allows selective updates based on what was provided in the request.
	updateParams := make(map[string]interface{})

	// Conditionally add non-pointer fields to updateParams if they are not empty/zero.
	// This means if a client sends an empty string for username, it won't be updated.
	if updates.Username != "" {
		updateParams["username"] = updates.Username
	}
	if updates.Email != "" {
		updateParams["email"] = updates.Email
	}
	if updates.FirstName != "" {
		updateParams["first_name"] = updates.FirstName
	}
	if updates.LastName != "" {
		updateParams["last_name"] = updates.LastName
	}
	if !updates.DateOfBirth.IsZero() { // Check if date of birth is not its zero value
		updateParams["date_of_birth"] = updates.DateOfBirth
	}

	// Conditionally add pointer fields to updateParams if they are not nil.
	// Dereference the pointer to get the actual value for the map.
	// If the pointer is nil, the field was not provided in the request body.
	if updates.Bio != nil {
		updateParams["bio"] = *updates.Bio
	}
	if updates.Location != nil {
		updateParams["location"] = *updates.Location
	}
	if updates.ProfilePictureAssetID != nil {
		updateParams["profile_picture_asset_id"] = *updates.ProfilePictureAssetID
	}
	// The `IsVerified` field is intentionally not handled here,
	// preventing its update via this endpoint.

	// Prepare loggable parameters by making a copy and redacting sensitive info.
	loggableParams := make(map[string]interface{})
	for k, v := range updateParams {
		if k != "password" { // Ensure password (if added later) is never logged
			loggableParams[k] = v
		}
	}
	loggableParams["id"] = userID       // Add user ID to loggable parameters
	apiResponse.Params = loggableParams // Update apiResponse.Params for logging

	// --- VALIDATION SECTION ---
	// Perform validation only on the fields that were provided for update.

	// Validate email format if email was provided in the update.
	if _, ok := updateParams["email"]; ok {
		if !validation.ValidateEmail(updateParams["email"].(string)) {
			apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid email format", nil)
			log.Printf("WARN: [%s][%s] Validation error: invalid email for update. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
			models.SendApiResponse(w, apiResponse)
			return
		}
	}

	// Handle password update: validate and hash if provided.
	if updates.Password != "" { // Check from the DTO, which holds the raw password input
		println("whyimhere")
		if !validation.ValidatePassword(updates.Password) { // Validate the raw password
			apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid password format", nil)
			log.Printf("WARN: [%s][%s] Validation error: invalid password for update. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
			models.SendApiResponse(w, apiResponse)
			return
		}
		// Hash the new password before adding it to updateParams.
		hashedPassword, err := passwordHashing.HashPassword(updates.Password)
		if err != nil {
			apiResponse.Error = apierrors.ERROR_CODE_INTERNAL_SERVER.ApiErrorResponse("Failed to hash new password", nil)
			log.Printf("ERROR: [%s][%s] Error hashing new password for user ID %s: %v", apiResponse.Context, apiResponse.Method, userID, err)
			models.SendApiResponse(w, apiResponse)
			return
		}
		updateParams["password"] = hashedPassword // Add the hashed password to the map
	}

	// Validate username format if username was provided.
	if _, ok := updateParams["username"]; ok {
		if !validation.ValidateUsername(updateParams["username"].(string)) {
			apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid username format", nil)
			log.Printf("WARN: [%s][%s] Validation error: invalid username for update. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
			models.SendApiResponse(w, apiResponse)
			return
		}
	}

	// Validate first name format if first name was provided.
	if _, ok := updateParams["first_name"]; ok {
		if !validation.ValidateName(updateParams["first_name"].(string)) {
			apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid first name format", nil)
			log.Printf("WARN: [%s][%s] Validation error: invalid first name for update. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
			models.SendApiResponse(w, apiResponse)
			return
		}
	}

	// Validate last name format if last name was provided.
	if _, ok := updateParams["last_name"]; ok {
		if !validation.ValidateName(updateParams["last_name"].(string)) {
			apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid last name format", nil)
			log.Printf("WARN: [%s][%s] Validation error: invalid last name for update. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
			models.SendApiResponse(w, apiResponse)
			return
		}
	}
	// --- END VALIDATION SECTION ---

	// Perform the database update using GORM's Updates method with the map.
	// This method updates only the columns specified in the map.
	result = db.Model(&existingUser).Updates(updateParams)

	if result.Error != nil {
		// Check for specific PostgreSQL unique constraint violation errors.
		if pgErr, ok := result.Error.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // Unique violation error code
				// Handle email unique constraint specifically
				if pgErr.ConstraintName == "uni_users_email" {
					apiResponse.Error = apierrors.ERROR_CODE_UNIQUE_KEY_VIOLATION.ApiErrorResponse("A user with this email address already exists", nil)
					log.Printf("WARN: [%s][%s] Conflict: User with email '%s' already exists during update for ID: %s.", apiResponse.Context, apiResponse.Method, updates.Email, userID)
					models.SendApiResponse(w, apiResponse)
					return
				}
				apiResponse.Error = apierrors.ERROR_CODE_UNIQUE_KEY_VIOLATION.ApiErrorResponse("A user with similar details already exists", nil)
				log.Printf("WARN: [%s][%s] Conflict: Unhandled unique constraint violation (code 23505) on constraint %s for user ID %s", apiResponse.Context, apiResponse.Method, pgErr.ConstraintName, userID)
				models.SendApiResponse(w, apiResponse)
				return
			}
		}
		// Handle generic database errors
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("Error updating user due to a database issue", nil)
		log.Printf("ERROR: [%s][%s] Database error updating user ID %s: %v", apiResponse.Context, apiResponse.Method, userID, result.Error)
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Re-fetch the user to get the latest state from the database,
	// including any fields updated by GORM (like UpdatedAt timestamp).
	result = db.First(&existingUser, "id = ?", userID)
	if result.Error != nil {
		// This scenario should be rare if the update just succeeded, but handles potential issues.
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("Successfully updated user but failed to re-fetch", nil)
		log.Printf("ERROR: [%s][%s] Successfully updated user ID %s but failed to re-fetch it: %v", apiResponse.Context, apiResponse.Method, userID, result.Error)
		models.SendApiResponse(w, apiResponse)
		return
	}

	// For security, clear the password hash from the user object before sending it in the response.
	existingUser.Password = ""

	// Prepare and send the successful API response.
	apiResponse.Message = "User updated successfully."
	apiResponse.Data = &models.ResponseData[models.User]{
		Items: []models.User{existingUser}, // Return the full, updated user object
	}

	log.Printf("INFO: [%s][%s] Successfully updated user ID: %s", apiResponse.Context, apiResponse.Method, userID)
	models.SendApiResponse(w, apiResponse)
}
