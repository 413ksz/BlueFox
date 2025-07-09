package users

import (
	"time"
)

// UserCreateRequestDTO defines the structure for expected incoming user creation requests.
type UserCreateRequestDTO struct {
	Username    string    `json:"username" validate:"required,min=3,max=20,username"`
	Email       string    `json:"email" validate:"required,max=254,email"`
	Password    string    `json:"password" validate:"required,min=8,max=72,password"`
	DateOfBirth time.Time `json:"dateOfBirth" validate:"required,dateofbirth"`
}

// UserUpdateRequestDTO defines the structure for expected incoming user update requests.
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
