package valueobject

import (
	"strings"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/shared/validation"
)

// Password is a value object that represents a password.
type Password struct {
	value string
}

func NewPassword(value string) (Password, *models.ValidationError) {

	// Trim leading and trailing whitespace from the password.
	// This helps prevent issues with accidental spaces or copy-paste errors.
	// We want to validate the trimmed version of the password.
	// and return the trimmed version of the password so it happens here not the validation function.
	trimmedValue := strings.TrimSpace(value)

	// Check if the password is valid.
	if err := validation.ValidatePasswordForEntropy(trimmedValue); err != nil {
		return Password{}, err
	}

	return Password{value: trimmedValue}, nil
}

// String returns the string representation of the password value object.
func (p Password) String() string {
	return p.value
}
