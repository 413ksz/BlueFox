package command

import "time"

type UserCreateCommand struct {
	Username    string
	Email       string
	Password    string
	DateOfBirth time.Time
}
