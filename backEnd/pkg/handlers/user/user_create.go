package user

import (
	"encoding/json"
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/apierrors"
	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	passwordHashing "github.com/413ksz/BlueFox/backEnd/pkg/password_hashing"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

// UserCreateHandler handles HTTP PUT requests for creating a new user.
// It expects a JSON request body containing user data, saves it to the database,
func UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	const (
		COMPONENT      string = "user_handler"
		METHOD_NAME    string = "UserCreateHandler"
		CONTEXT        string = "api/user"
		METHOD         string = "PUT"
		STATUS_DEFAULT int    = http.StatusCreated
	)

	apiResponse := &models.ApiResponse[models.User]{}
	apiResponse.Method = METHOD
	apiResponse.Context = CONTEXT
	apiResponse.StatusCode = STATUS_DEFAULT

	// Get the GORM database instance from your database package.
	db := database.DB

	log.Info().
		Str("component", COMPONENT).
		Str("method_name", METHOD_NAME).
		Str("http_method", r.Method).
		Str("path", r.URL.Path).
		Str("event", "http_request_received").
		Msg("Processing new user creation request.")

	// Check if the database connection is initialized.
	if db == nil {
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_INITIALIZE.ApiErrorResponse("Database not ready for UserCreateHandler", nil)
		log.Error().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "db_not_initialized").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Msg("Database not initialized for user creation.")
		models.SendApiResponse(w, apiResponse)
		return
	}

	var newUser models.User

	// Decode the JSON request body into the newUser struct.
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		apiResponse.Error = apierrors.ERROR_CODE_ENCODE_ERROR.ApiErrorResponse("Invalid JSON data", nil)
		log.Warn().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "request_body_decode_failed").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Err(err).
			Msg("Error decoding request body.")
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
		log.Warn().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "validation_failed_missing_fields").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Fields(apiResponse.Params).
			Msg("Validation error: required fields are missing.")
		models.SendApiResponse(w, apiResponse)
		return
	}

	if !validation.ValidateEmail(newUser.Email) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid email format", nil)
		log.Warn().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "validation_failed_invalid_email").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Str("email", newUser.Email).
			Msg("Validation error: invalid email format.")
		models.SendApiResponse(w, apiResponse)
		return
	}

	if !validation.ValidatePassword(newUser.Password) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid password format", nil)
		log.Warn().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "validation_failed_invalid_password").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Msg("Validation error: invalid password format.")
		models.SendApiResponse(w, apiResponse)
		return
	}

	if !validation.ValidateUsername(newUser.Username) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid username format", nil)
		log.Warn().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "validation_failed_invalid_username").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Str("username", newUser.Username).
			Msg("Validation error: invalid username format.")
		models.SendApiResponse(w, apiResponse)
		return
	}

	if !validation.ValidateName(newUser.FirstName) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid first name format", nil)
		log.Warn().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "validation_failed_invalid_lastname").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Str("lastName", newUser.LastName).
			Msg("Validation error: invalid last name format.")
		models.SendApiResponse(w, apiResponse)
		return
	}

	if !validation.ValidateName(newUser.LastName) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid last name format", nil)
		log.Warn().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "validation_failed_invalid_firstname").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Str("firstName", newUser.FirstName).
			Msg("Validation error: invalid last name format.")
		models.SendApiResponse(w, apiResponse)
		return
	}

	passwordHash, err := passwordHashing.HashPassword(newUser.Password)
	if err != nil {
		apiResponse.Error = apierrors.ERROR_CODE_INTERNAL_SERVER.ApiErrorResponse("Failed to hash password", nil)
		log.Error().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "password_hashing_failed").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Err(err). // Logs the underlying hashing error
			Msg("Error hashing password.")
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
					log.Warn().
						Str("component", COMPONENT).
						Str("method_name", METHOD_NAME).
						Str("event", "unique_constraint_violation").
						Str("api_error_code", apiResponse.Error.Code).
						Str("api_error_message", apiResponse.Error.Message).
						Str("constraint_name", pgErr.ConstraintName).
						Str("duplicate_field", "email").
						Str("email", newUser.Email).
						Err(pgErr).
						Msg("Conflict: User with this email already exists.")
					models.SendApiResponse(w, apiResponse)
					return
				}
				// Fallback for any other unique constraint violation not specifically handled
				apiResponse.Error = apierrors.ERROR_CODE_UNIQUE_KEY_VIOLATION.ApiErrorResponse("A user with similar details already exists", nil)
				log.Warn().
					Str("component", COMPONENT).
					Str("method_name", METHOD_NAME).
					Str("event", "unhandled_unique_constraint_violation").
					Str("api_error_code", apiResponse.Error.Code).
					Str("api_error_message", apiResponse.Error.Message).
					Str("constraint_code", pgErr.Code).
					Str("constraint_name", pgErr.ConstraintName).
					Err(pgErr).
					Msg("Conflict: Unhandled unique constraint violation.")
				models.SendApiResponse(w, apiResponse)
				return
			}
		}

		// Fallback for any other database errors (e.g., connection issues, other integrity errors)
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("Error creating user due to a database issue", nil)
		log.Error().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "database_error_creating_user").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Err(result.Error).
			Msg("Database error creating user.")
		models.SendApiResponse(w, apiResponse)
		return
	}

	newUser.Password = "" // Clear the password before sending in the response.

	// If the user is successfully created, return the user data as JSON.
	apiResponse.Message = "User created successfully."

	log.Info().
		Str("component", COMPONENT).
		Str("method_name", METHOD_NAME).
		Str("event", "user_created_success").
		Str("username", newUser.Username).
		Str("email", newUser.Email).
		Msg("Successfully created user.")

	models.SendApiResponse(w, apiResponse)
}
