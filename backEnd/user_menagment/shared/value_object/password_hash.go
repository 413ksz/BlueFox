package valueobject

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/shared/validation"
)

// PasswordHash is a value object that represents a password hash.
type PasswordHash struct {
	value string
}

// NewPasswordHash creates a new password hash value object.
// It checks for domain invariants and returns an error if the input is invalid
func NewPasswordHash(value string) (PasswordHash, *models.ValidationError) {
	if validationError := validation.ValidatePasswordHash(value); validationError != nil {
		return PasswordHash{}, validationError
	}

	return PasswordHash{value: value}, nil
}

// String returns the string representation of the password hash value object.
func (p PasswordHash) String() string {
	return p.value
}
