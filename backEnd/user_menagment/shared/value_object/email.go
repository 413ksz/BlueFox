package valueobject

import (
	"strings"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/shared/validation"
)

// Email is a value object that represents an Email address.
type Email struct {
	value string
}

// NewEmail creates a new email value object.
// It checks for domain invariants and returns an error if the input is invalid
func NewEmail(value string) (Email, *models.ValidationError) {
	trimmedValue := strings.TrimSpace(value)
	trimmedValue = strings.ToLower(trimmedValue)
	if validationError := validation.ValidateEmailString(value); validationError != nil {
		return Email{}, validationError
	}

	return Email{value: trimmedValue}, nil
}

// String returns the string representation of the email value object.
func (e Email) String() string {
	return e.value
}

// Equals checks if two email value objects are equal.
func (e Email) Equals(otherEmail Email) bool {
	return e.value == otherEmail.value
}
