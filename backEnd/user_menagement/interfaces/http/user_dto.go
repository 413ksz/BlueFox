package http

import (
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/application/command" // for command struct definitions
	valueobject "github.com/413ksz/BlueFox/backEnd/user_menagement/shared/value_object"
)

// UserCreateRequestDTO represents the request body for creating a new user.
// it contains the necessary fields for creating a user.
type UserCreateRequestDTO struct {
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	DateOfBirth time.Time `json:"dateOfBirth"`
}

// ToCreateUserCommand converts the UserCreateRequestDTO into a
// command.UserCreateCommand. This transformation prepares the data
// for processing by the application layer, stripping away HTTP-specific
// concerns like JSON unmarshalling
func (dto *UserCreateRequestDTO) ToCreateUserCommand() (command.UserCreateCommand, *models.CustomError) {

	// ------------ Validation ------------
	// Validate the input values of the DTO
	// and return an error if any of them are invalid
	validationErrors := models.NewValidationErrors()
	emaiVO, err := valueobject.NewEmail(dto.Email)
	if err != nil {
		validationErrors.AddNewError(*err)
	}
	usernameVO, err := valueobject.NewUsername(dto.Username)
	if err != nil {
		validationErrors.AddNewError(*err)
	}
	passwordVO, err := valueobject.NewPassword(dto.Password)
	if err != nil {
		validationErrors.AddNewError(*err)
	}
	dateOfBirthVO, err := valueobject.NewDateOfBirth(dto.DateOfBirth)
	if err != nil {
		validationErrors.AddNewError(*err)
	}
	if validationErrors.Len() > 0 {
		return command.UserCreateCommand{}, models.NewCustomError(models.ERROR_CODE_UNPROCESSABLE_ENTITY, validationErrors, nil, nil)
	}

	return command.UserCreateCommand{
		Username:    usernameVO,
		Email:       emaiVO,
		Password:    passwordVO,
		DateOfBirth: dateOfBirthVO,
	}, nil

}

// UserUpdateRequestDTO defines the structure for expected incoming user update requests.
// It json tags are used for unmarshalling the request and the validation rules
type UserUpdateRequestDTO struct {
	Username    *string    `json:"username" validate:"min=3,max=20,username"`
	Email       *string    `json:"email" validate:"max=254,email"`
	Password    *string    `json:"password" validate:"min=8,max=72,password"`
	DateOfBirth *time.Time `json:"dateOfBirth" validate:"dateofbirth"`
	FirstName   *string    `json:"firstName" validate:"name"`
	LastName    *string    `json:"lastName" validate:"name"`
	Bio         *string    `json:"bio" validate:"name"`
}

// ToUpdateUserCommand converts the UserUpdateRequestDTO into a command.UserUpdateCommand.
// It strips away the JSON tags and validation rules, leaving only the data for processing
func (dto *UserUpdateRequestDTO) ToUpdateUserCommand() *command.UserUpdateCommand {
	return &command.UserUpdateCommand{
		Username:    dto.Username,
		Email:       dto.Email,
		Password:    dto.Password,
		DateOfBirth: dto.DateOfBirth,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Bio:         dto.Bio,
	}
}
