package user

import (
	"encoding/json"
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/apierrors"
	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	passwordHashing "github.com/413ksz/BlueFox/backEnd/pkg/password_hashing"
	jwt_token "github.com/413ksz/BlueFox/backEnd/pkg/token"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	const (
		COMPONENT      string = "user_handler"
		METHOD_NAME    string = "UserLoginHandler"
		CONTEXT        string = "api/user/login"
		METHOD         string = "POST"
		STATUS_DEFAULT int    = http.StatusOK
	)
	db := database.DB

	apiResponse := &models.ApiResponse[models.User]{}
	apiResponse.Method = METHOD
	apiResponse.Context = CONTEXT
	apiResponse.StatusCode = STATUS_DEFAULT

	log.Info().
		Str("component", COMPONENT).
		Str("method_name", METHOD_NAME).
		Str("http_method", METHOD).
		Str("path", CONTEXT).
		Str("event", "http_request_received").
		Msg("Processing user login request.")

	if db == nil {
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_INITIALIZE.ApiErrorResponse("Database not ready for UserLoginHandler", nil)
		log.Error().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "db_not_initialized").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Msg("Database not initialized for loging in the user.")
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Decode the JSON request body into the user struct.
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		apiResponse.Error = apierrors.ERROR_CODE_ENCODE_ERROR.ApiErrorResponse("Invalid request body", err)
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

	// Check if the required fields are present in the request body.
	if user.Email == "" || user.Password == "" {
		apiResponse.Error = apierrors.ERROR_CODE_INVALID_INPUT.ApiErrorResponse("Invalid request body", nil)
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

	// Validate email format
	if !validation.ValidateEmail(user.Email) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid request body", nil)
		log.Warn().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "validation_failed_invalid_email").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Str("email", user.Email).
			Msg("Validation error: invalid email format.")
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Validate email format
	if !validation.ValidatePassword(user.Password) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid request body", nil)
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

	fetchedUser := models.User{}
	// Attempt to find a user with the provided email and password in the database.
	result := db.Where("email = ?", user.Email).Preload("ProfilePictureAsset").First(&fetchedUser).Select("id", "username", "profile_picture_asset_id")
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			apiResponse.Error = apierrors.ERROR_CODE_NOT_FOUND.ApiErrorResponse("User not found", nil)
			log.Error().
				Str("component", COMPONENT).
				Str("method_name", METHOD_NAME).
				Str("event", "user_not_found").
				Str("api_error_code", apiResponse.Error.Code).
				Str("api_error_message", apiResponse.Error.Message).
				Int("api_error_status", apiResponse.Error.HTTPStatusCode).
				Str("id", user.ID.String()).
				Err(result.Error).
				Msg("User not found")
			models.SendApiResponse(w, apiResponse)
			return
		}
		// Other database errors are internal server errors
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("Error fetching user data for login", nil)
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
	if !passwordHashing.VerifyPassword(user.Password, fetchedUser.Password) {
		apiResponse.Error = apierrors.ERROR_CODE_UNAUTHORIZED.ApiErrorResponse("Invalid credentials", nil)
		log.Error().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "invalid_password").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Msg("Invalid password")
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Generate a JWT token for the user
	token, err := jwt_token.GenerateJWTToken(user.Username, user.ID.String(), user.ProfilePictureAsset.UrlPath)

	if err != nil {
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("Error generating JWT token", nil)
		log.Error().
			Str("component", COMPONENT).
			Str("method_name", METHOD_NAME).
			Str("event", "jwt_token_error").
			Str("api_error_code", apiResponse.Error.Code).
			Str("api_error_message", apiResponse.Error.Message).
			Int("api_error_status", apiResponse.Error.HTTPStatusCode).
			Err(err).
			Msg("Error generating JWT token")
		models.SendApiResponse(w, apiResponse)
	}

	log.Info().
		Str("component", COMPONENT).
		Str("method_name", METHOD_NAME).
		Str("event", "user_login_success").
		Msg("User logged in successfully")

	w.Header().Set("Authorization", "Bearer: "+token)
	models.SendApiResponse(w, apiResponse)

}
