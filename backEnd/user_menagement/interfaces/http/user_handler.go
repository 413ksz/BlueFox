package http

import (
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/application/service"
)

// UserHandler encapsulates the logic for handling user-related operations.
type UserHandler struct {
	UserService service.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
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
	jsonParseError := validation.ValidateRequestBody(&dto, r, 0)
	if jsonParseError != nil {
		apiResponse.WithError(jsonParseError.Message, jsonParseError, jsonParseError.HttpCode)
		return apiResponse, jsonParseError
	}

	// Convert the DTO to a CreateUserCommand
	command, validationError := dto.ToCreateUserCommand()
	if validationError != nil {
		apiResponse.WithError(validationError.Message, validationError, validationError.HttpCode)
		return apiResponse, validationError
	}

	// Call the service layer to create the user
	createError := userHandler.UserService.CreateUser(command)
	if createError != nil {
		apiResponse.WithError(createError.Message, createError, createError.HttpCode)
		return apiResponse, createError
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
