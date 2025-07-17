package http

import (
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/application/service"
)

// UserHandler encapsulates the logic for handling user-related operations.
type UserHandler struct {
	Validator   *validation.Validator
	UserService service.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userService service.UserService, validator *validation.Validator) *UserHandler {
	return &UserHandler{
		Validator:   validator,
		UserService: userService,
	}
}

// UserCreateHandler handles the creation of a new user.
// it handles request body validation, DTO conversion, and service layer call.
// paramaters:
// - w: the HTTP response writer
// - r: the HTTP request
// returns:
// - ApiResponse: the response to be sent to the client
// - CustomError: any error that occurred
func (userHandler *UserHandler) UserCreateHandler(w http.ResponseWriter, r *http.Request) (*models.ApiResponse[any], *models.CustomError) {

	// Create a new ApiResponse and a UserCreateRequestDTO
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

	// Convert the DTO to a CreateUserCommand
	command := dto.ToCreateUserCommand()

	// Call the service layer to create the user
	customError = userHandler.UserService.CreateUser(command)
	if customError != nil {
		apiResponse.WithError(customError.Message, customError, customError.HttpCode)
		return apiResponse, customError
	}

	// Return a success response
	apiResponse.WithData("User created successfully", nil, http.StatusCreated)
	return apiResponse, nil
}

// UserGetHandler handles HTTP GET requests for fetching a single user by ID.
// It retrieves the user ID (UUID) from the URL path variables (e.g., /users/{id}),
// queries the database, and returns the user data as JSON.
func UserGetHandler(w http.ResponseWriter, r *http.Request) (*models.ApiResponse[any], *models.CustomError) {

	apiReasponse := models.NewApiResponse[any](nil, http.StatusBadRequest, "Bad request")

	return apiReasponse, nil
}
