package valueobject

import (
	"strings"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/shared/validation"
)

type Name struct {
	value string
}

func NewName(value string) (Name, *models.ValidationError) {
	trimmedValue := strings.TrimSpace(value)
	if err := validation.ValidateNameString(value); err != nil {
		return Name{}, err
	}

	return Name{value: trimmedValue}, nil
}

func (n Name) String() string {
	return n.value
}

func (n Name) Equals(otherName Name) bool {
	return n.value == otherName.value
}
