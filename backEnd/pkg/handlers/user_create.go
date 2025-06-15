package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/apierrors"
	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	passwordHashing "github.com/413ksz/BlueFox/backEnd/pkg/password_hashing"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserCreateHandler handles HTTP PUT requests for creating a new user.
// It expects a JSON request body containing user data, saves it to the database,
func UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	const (
		CONTEXT        string = "api/user/"
		METHOD         string = "PUT"
		STATUS_DEFAULT int    = http.StatusCreated
	)

	apiResponse := &models.ApiResponse[models.User]{}
	apiResponse.Method = METHOD
	apiResponse.Context = CONTEXT
	apiResponse.StatusCode = STATUS_DEFAULT

	// Get the GORM database instance from your database package.
	db := database.DB

	// Check if the database connection is initialized.
	if db == nil {
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_INITIALIZE.ApiErrorResponse("Database not ready for UserCreateHandler", nil)
		log.Printf("ERROR: [%s][%s] Database not initialized. Error: %s", apiResponse.Context, apiResponse.Method, apiResponse.Error.Details)
		models.SendApiResponse(w, apiResponse)
		return
	}

	var newUser models.User

	// Decode the JSON request body into the newUser struct.
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		apiResponse.Error = apierrors.ERROR_CODE_ENCODE_ERROR.ApiErrorResponse("Invalid JSON data", nil)
		log.Printf("ERROR: [%s][%s] Error decoding request body: %v", apiResponse.Context, apiResponse.Method, err)
		models.SendApiResponse(w, apiResponse)
		return
	}

	apiResponse.Params = map[string]interface{}{
		"userName":    newUser.Username,
		"email":       newUser.Email,
		"lastName":    newUser.LastName,
		"firstName":   newUser.FirstName,
		"dateOfBirth": newUser.DateOfBirth,
	}

	// Perform validation checks on the user data.
	if newUser.Email == "" || newUser.Username == "" || newUser.Password == "" || newUser.FirstName == "" || newUser.LastName == "" || newUser.DateOfBirth.IsZero() {
		apiResponse.Error = apierrors.ERROR_CODE_INVALID_INPUT.ApiErrorResponse("Missing required fields", nil)
		log.Printf("WARN: [%s][%s] Validation error: required fields are missing. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
		models.SendApiResponse(w, apiResponse)
		return
	}

	if !validation.ValidateEmail(newUser.Email) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid email format", nil)
		log.Printf("WARN: [%s][%s] Validation error: invalid email. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
		models.SendApiResponse(w, apiResponse)
		return
	}

	if !validation.ValidatePassword(newUser.Password) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid password format", nil)
		log.Printf("WARN: [%s][%s] Validation error: invalid password. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
		models.SendApiResponse(w, apiResponse)
		return
	}

	if !validation.ValidateUsername(newUser.Username) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid username format", nil)
		log.Printf("WARN: [%s][%s] Validation error: invalid username. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
		models.SendApiResponse(w, apiResponse)
		return
	}

	if !validation.ValidateName(newUser.FirstName) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid first name format", nil)
		log.Printf("WARN: [%s][%s] Validation error: invalid first name. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
		models.SendApiResponse(w, apiResponse)
		return
	}

	if !validation.ValidateName(newUser.LastName) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid last name format", nil)
		log.Printf("WARN: [%s][%s] Validation error: invalid last name. Params: %+v", apiResponse.Context, apiResponse.Method, apiResponse.Params)
		models.SendApiResponse(w, apiResponse)
		return
	}

	passwordHash, err := passwordHashing.HashPassword(newUser.Password)
	if err != nil {
		apiResponse.Error = apierrors.ERROR_CODE_INTERNAL_SERVER.ApiErrorResponse("Failed to hash password", nil)
		log.Printf("ERROR: [%s][%s] Error hashing password: %v", apiResponse.Context, apiResponse.Method, err)
		models.SendApiResponse(w, apiResponse)
		return
	}
	newUser.Password = passwordHash

	// Attempt to create (insert) the new user record into the database using GORM.
	result := db.Create(&newUser)

	// Check for errors returned by the database operation.
	if result.Error != nil {
		// Attempt to unwrap the error to check for specific database driver errors
		if pgErr, ok := result.Error.(*pgconn.PgError); ok {
			// Check for PostgreSQL unique violation error code (23505)
			if pgErr.Code == "23505" {
				// Differentiate between unique constraints if you have more than one
				// (e.g., on 'email' and 'username')
				if pgErr.ConstraintName == "uni_users_email" { // Make sure "uni_users_email" matches your actual unique constraint name for the email field
					apiResponse.Error = apierrors.ERROR_CODE_UNIQUE_KEY_VIOLATION.ApiErrorResponse("A user with this email address already exists", nil)
					log.Printf("WARN: [%s][%s] Conflict: User with email '%s' already exists.", apiResponse.Context, apiResponse.Method, newUser.Email)
					models.SendApiResponse(w, apiResponse)
					return
				}
				// Fallback for any other unique constraint violation not specifically handled
				apiResponse.Error = apierrors.ERROR_CODE_UNIQUE_KEY_VIOLATION.ApiErrorResponse("A user with similar details already exists", nil)
				log.Printf("WARN: [%s][%s] Conflict: Unhandled unique constraint violation (code 23505) on constraint: %s", apiResponse.Context, apiResponse.Method, pgErr.ConstraintName)
				models.SendApiResponse(w, apiResponse)
				return
			}
		}

		// Fallback for any other database errors (e.g., connection issues, other integrity errors)
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("Error creating user due to a database issue", nil)
		log.Printf("ERROR: [%s][%s] Database error creating user: %v", apiResponse.Context, apiResponse.Method, result.Error)
		models.SendApiResponse(w, apiResponse)
		return
	}

	newUser.Password = "" // Clear the password before sending in the response.

	// If the user is successfully created, return the user data as JSON.
	apiResponse.Message = "User created successfully."

	log.Printf("INFO: [%s][%s] Successfully created user: %s", apiResponse.Context, apiResponse.Method, newUser.Username)
	models.SendApiResponse(w, apiResponse)
}
