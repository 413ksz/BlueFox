package valueobject

import (
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/shared/validation"
)

type DateOfBirth struct {
	value time.Time
}

func NewDateOfBirth(value time.Time) (DateOfBirth, *models.ValidationError) {
	if validationError := validation.ValidateDateOfBirthTime(value); validationError != nil {
		return DateOfBirth{}, validationError
	}
	return DateOfBirth{value: value}, nil
}

func (d DateOfBirth) Time() time.Time {
	return d.value
}

func (d DateOfBirth) Equals(otherDateOfBirth DateOfBirth) bool {
	return d.value.Equal(otherDateOfBirth.value)
}
