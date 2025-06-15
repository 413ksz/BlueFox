package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/apierrors"
	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	passwordHashing "github.com/413ksz/BlueFox/backEnd/pkg/password_hashing"
	jwt_token "github.com/413ksz/BlueFox/backEnd/pkg/token"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
	"gorm.io/gorm"
)

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	const (
		CONTEXT        string = "api/user/login"
		METHOD         string = "POST"
		STATUS_DEFAULT int    = http.StatusOK
	)
	db := database.DB

	apiResponse := &models.ApiResponse[models.User]{}
	apiResponse.Method = METHOD
	apiResponse.Context = CONTEXT
	apiResponse.StatusCode = STATUS_DEFAULT

	if db == nil {
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_INITIALIZE.ApiErrorResponse("Database not ready for UserUpdateHandler", nil)
		log.Printf("ERROR: [%s][%s] Database not initialized. Error: %s", apiResponse.Context, apiResponse.Method, apiResponse.Error.Details)
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Decode the JSON request body into the user struct.
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		apiResponse.Error = apierrors.ERROR_CODE_ENCODE_ERROR.ApiErrorResponse("Invalid request body", err)
		log.Printf("ERROR: [%s][%s] Invalid request body. Error: %s", apiResponse.Context, apiResponse.Method, apiResponse.Error.Details)
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Check if the required fields are present in the request body.
	if user.Email == "" || user.Password == "" {
		apiResponse.Error = apierrors.ERROR_CODE_INVALID_INPUT.ApiErrorResponse("Invalid request body", nil)
		log.Printf("ERROR: [%s][%s] Invalid request body. Error: %s", apiResponse.Context, apiResponse.Method, apiResponse.Error.Details)
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Validate email format
	if !validation.ValidateEmail(user.Email) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid request body", nil)
		log.Printf("ERROR: [%s][%s] Invalid request body. Error: %s", apiResponse.Context, apiResponse.Method, apiResponse.Error.Details)
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Validate email format
	if !validation.ValidatePassword(user.Password) {
		apiResponse.Error = apierrors.ERROR_CODE_VALIDATION_FAILED.ApiErrorResponse("Invalid request body", nil)
		log.Printf("ERROR: [%s][%s] Invalid request body. Error: %s", apiResponse.Context, apiResponse.Method, apiResponse.Error.Details)
		models.SendApiResponse(w, apiResponse)
		return
	}

	fetchedUser := models.User{}
	// Attempt to find a user with the provided email and password in the database.
	result := db.Where("email = ?", user.Email).Preload("ProfilePictureAsset").First(&fetchedUser).Select("id", "username", "profile_picture_asset_id")
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			apiResponse.Error = apierrors.ERROR_CODE_NOT_FOUND.ApiErrorResponse("User not found", nil)
			log.Printf("INFO: [%s][%s] User not found for ID: %s", apiResponse.Context, apiResponse.Method, user.Email)
			models.SendApiResponse(w, apiResponse)
			return
		}
		// Other database errors are internal server errors
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("Error fetching user for update", nil)
		log.Printf("ERROR: [%s][%s] Database error fetching user ID %s: %v", apiResponse.Context, apiResponse.Method, user.Email, result.Error)
		models.SendApiResponse(w, apiResponse)
		return
	}
	if !passwordHashing.VerifyPassword(user.Password, fetchedUser.Password) {
		apiResponse.Error = apierrors.ERROR_CODE_UNAUTHORIZED.ApiErrorResponse("Invalid credentials", nil)
		log.Printf("ERROR: [%s][%s] Invalid credentials for user ID %s", apiResponse.Context, apiResponse.Method, user.Email)
		models.SendApiResponse(w, apiResponse)
		return
	}

	// Generate a JWT token for the user
	token, err := jwt_token.GenerateJWTToken(user.Username, user.ID.String(), user.ProfilePictureAsset.UrlPath)

	if err != nil {
		apiResponse.Error = apierrors.ERROR_CODE_DATABASE_ERROR.ApiErrorResponse("Error generating JWT token", nil)
		log.Printf("ERROR: [%s][%s] Error generating JWT token for user ID %s: %v", apiResponse.Context, apiResponse.Method, user.Email, result.Error)
		models.SendApiResponse(w, apiResponse)
	}
	w.Header().Set("Authorization", "Bearer: "+token)
	models.SendApiResponse(w, apiResponse)

}
