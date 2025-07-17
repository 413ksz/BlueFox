package http

import (
	"time"

	"github.com/413ksz/BlueFox/backEnd/user_menagment/application/command" // for command struct definitions
)

// UserCreateRequestDTO defines the structure for incoming HTTP requests
// to create a new user.
// It includes JSON tags for request unmarshaling and validation rules
// for immediate input validation at the API boundary.
type UserCreateRequestDTO struct {
	Username    string    `json:"username" validate:"required,min=3,max=30,username"`
	Email       string    `json:"email" validate:"required,max=254,email"`
	Password    string    `json:"password" validate:"required,min=8,max=72,password"`
	DateOfBirth time.Time `json:"dateOfBirth" validate:"required,dateofbirth"`
}

// ToCreateUserCommand converts the UserCreateRequestDTO into a
// command.UserCreateCommand. This transformation prepares the data
// for processing by the application layer, stripping away HTTP-specific
// concerns like JSON and validation tags.
func (dto *UserCreateRequestDTO) ToCreateUserCommand() command.UserCreateCommand {
	return command.UserCreateCommand{Username: dto.Username, Email: dto.Email, Password: dto.Password, DateOfBirth: dto.DateOfBirth}
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
	Location    *string    `json:"location" validate:"name"`
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
		Location:    dto.Location,
		Bio:         dto.Bio,
	}
}
