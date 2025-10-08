package command

import (
	"time"

	"github.com/google/uuid"
)

type UserUpdateCommand struct {
	Id          uuid.UUID
	Username    *string
	Email       *string
	Password    *string
	FirstName   *string
	LastName    *string
	Bio         *string
	DateOfBirth *time.Time
	Location    *string
}
