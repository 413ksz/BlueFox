package model

import (
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	valueobject "github.com/413ksz/BlueFox/backEnd/user_menagment/shared/value_object"
	"github.com/google/uuid"
)

type User struct {
	Id                    uuid.UUID
	Username              valueobject.Username
	Email                 valueobject.Email
	PasswordHash          valueobject.PasswordHash
	FirstName             *valueobject.Name
	LastName              *valueobject.Name
	CreatedAt             time.Time
	UpdatedAt             *time.Time
	LastOnline            *time.Time
	Bio                   *string
	DateOfBirth           valueobject.DateOfBirth
	Location              *string
	IsVerified            bool
	ProfilePictureAssetID *uuid.UUID
}

// NewUser creates a new User instance
// it checks for domain invariants and returns an error if the input is invalid
// it is for user cration only
// params:
// - username: The username of the user
// - email: The email address of the user
// - passwordHash: The password hash of the user
// - dateOfBirth: The date of birth of the user
// returns:
// - *User: A pointer to the newly created User instance
// - *models.CustomError: An error if the input is invalid
func NewUser(username string, email string, passwordHash string, dateOfBirth time.Time) (*User, *models.CustomError) {
	var errors []*models.ValidationError

	// ------- Validations --------

	// username
	usernameVO, domainErr := valueobject.NewUsername(username)
	if domainErr != nil {
		errors = append(errors, domainErr)
	}

	// date of birth
	dateOfBirthVO, domainErr := valueobject.NewDateOfBirth(dateOfBirth)
	if domainErr != nil {
		errors = append(errors, domainErr)
	}

	// email
	emailVO, domainErr := valueobject.NewEmail(email)
	if domainErr != nil {
		errors = append(errors, domainErr)
	}

	// passwordHash
	passwordHashVO, domainErr := valueobject.NewPasswordHash(passwordHash)
	if domainErr != nil {
		errors = append(errors, domainErr)
	}

	// check if there are any validation errors
	if len(errors) > 0 {
		return nil, models.NewCustomError(models.ERROR_CODE_UNPROCESSABLE_ENTITY, nil, errors, nil)
	}

	// ------- Create User --------
	return &User{
		Id:                    uuid.New(),
		Username:              usernameVO,
		Email:                 emailVO,
		PasswordHash:          passwordHashVO,
		DateOfBirth:           dateOfBirthVO,
		IsVerified:            false,
		CreatedAt:             time.Now().UTC(),
		UpdatedAt:             nil,
		LastOnline:            nil,
		Bio:                   nil,
		Location:              nil,
		FirstName:             nil,
		LastName:              nil,
		ProfilePictureAssetID: nil,
	}, nil
}
