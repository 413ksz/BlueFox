package command

import (
	valueobject "github.com/413ksz/BlueFox/backEnd/user_menagement/shared/value_object"
)

type UserCreateCommand struct {
	Username    valueobject.Username
	Email       valueobject.Email
	Password    valueobject.Password
	DateOfBirth valueobject.DateOfBirth
}
