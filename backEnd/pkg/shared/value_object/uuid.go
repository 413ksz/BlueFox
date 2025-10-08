package valueobject

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/google/uuid"
)

// UUID is a value object that represents a UUID.
type ID struct {
	value uuid.UUID
}

// NewUUID creates a new UUID value object.
func NewUUID(value string) (ID, *models.ValidationError) {
	if value == "" {
		validationError := models.NewValidationError("uuid", value, "required", "uuid is required")
		return ID{}, validationError
	}

	uuid, parseErr := uuid.Parse(value)
	if parseErr != nil {
		validationError := models.NewValidationError("uuid", value, "invalid", "uuid is in invalid format")
		return ID{}, validationError
	}

	return ID{value: uuid}, nil
}

// String returns the string representation of the UUID.
func (u ID) String() string {
	return u.value.String()
}

// UUID returns the UUID value.
func (u ID) UUID() uuid.UUID {
	return u.value
}
