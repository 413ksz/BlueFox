package users

import (
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
)

// UserHandler encapsulates the logic for handling user-related operations.
type UserHandler struct {
	Validator   *validation.Validator
	UserService UserServiceLayer
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userService UserServiceLayer, validator *validation.Validator) *UserHandler {
	return &UserHandler{
		Validator:   validator,
		UserService: userService,
	}
}

func (userHandler *UserHandler) UserCreateHandler(w http.ResponseWriter, r *http.Request) (*models.ApiResponse[any], *models.CustomError) {

	apiResponse := models.NewApiResponse[any](nil, 0, "")
	var dto UserCreateRequestDTO

	// Validate the request body and unmarshal it into the DTO
	customError := userHandler.Validator.ValidateRequestBody(&dto, r)
	if customError != nil {
		apiResponse.WithError(customError.Message, customError, customError.HttpCode)
		return apiResponse, customError
	}

	// Validate the input values of the DTO
	customError = userHandler.Validator.ValidateDto(&dto)
	if customError != nil {
		apiResponse.WithError(customError.Message, customError, customError.HttpCode)
		return apiResponse, customError
	}

	// Call the service layer to create the user
	customError = userHandler.UserService.CreateUser(dto)
	if customError != nil {
		apiResponse.WithError(customError.Message, customError, customError.HttpCode)
		return apiResponse, customError
	}

	apiResponse.WithData("User created successfully", nil, http.StatusCreated)
	// return success
	return apiResponse, nil
}

// UserGetHandler handles HTTP GET requests for fetching a single user by ID.
// It retrieves the user ID (UUID) from the URL path variables (e.g., /users/{id}),
// queries the database, and returns the user data as JSON.
func UserGetHandler(w http.ResponseWriter, r *http.Request) (*models.ApiResponse[User], *models.CustomError) {

	apiReasponse := models.NewApiResponse[User](nil, http.StatusBadRequest, "Bad request")

	return apiReasponse, nil
}
