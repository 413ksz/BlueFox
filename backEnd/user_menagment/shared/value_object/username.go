package valueobject

import (
	"strings"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/shared/validation"
)

type Username struct {
	value string
}

func NewUsername(value string) (Username, *models.ValidationError) {
	trimmedValue := strings.TrimSpace(value)
	if validationError := validation.ValidateUsernameString(value); validationError != nil {
		return Username{}, validationError
	}
	return Username{value: trimmedValue}, nil
}

func (u Username) String() string {
	return u.value
}

func (u Username) Equals(otherUsername Username) bool {
	return u.value == otherUsername.value
}
