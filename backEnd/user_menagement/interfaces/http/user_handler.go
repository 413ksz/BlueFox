package http

import (
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/pkg/validation"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/application/query"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/application/service"

	"github.com/gorilla/mux"
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

// GetHandler handles the retrieval of a user data from the database.
// it only retrives public not personal data of the user
// paramaters:
// - w: the HTTP response writer
// - r: the HTTP request
// returns:
// - ApiResponse: the response to be sent to the client
// - CustomError: any error that occurred
func (userHandler *UserHandler) GetHandler(w http.ResponseWriter, r *http.Request) (*models.ApiResponse[any], *models.CustomError) {
	id := mux.Vars(r)["id"]
	apiResponse := models.NewApiResponse[any](nil, 0, "")

	userQuery, querryError := query.NewPublicUserQuery(id)
	if querryError != nil {
		apiResponse.WithError(querryError.Message, querryError, querryError.HttpCode)
		return apiResponse, querryError
	}

	user, serviceError := userHandler.UserService.GetUser(userQuery)
	if serviceError != nil {
		apiResponse.WithError(serviceError.Message, serviceError, serviceError.HttpCode)
		return apiResponse, serviceError
	}
	if user == nil {
		customError := models.NewCustomError(models.ERROR_CODE_NOT_FOUND, "User not found", nil, nil)
		apiResponse.WithError("User not found", customError, customError.HttpCode)
		return apiResponse, customError
	}
	responseData := &models.ResponseData[any]{
		Items: []any{user},
	}
	apiResponse.WithData("User retrieved successfully", responseData, http.StatusOK)
	return apiResponse, nil
}
