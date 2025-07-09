package users

import (
	"time"

	"github.com/google/uuid"
)

// User is the domain entity for a user.
type User struct {
	Id           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	FirstName    *string
	LastName     *string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	LastOnline   *time.Time
	Bio          *string
	DateOfBirth  time.Time
	Location     *string
	IsVerified   bool
}

// NewUser creates a new User instance.
// params:
// - username: The username of the user.
// - email: The email address of the user.
// - passwordhash: The hashed password of the user.
// - dateofbirth: The date of birth of the user.
// returns:
// - *User: A pointer to the newly created User instance.
func NewUser(username string, email string, passwordhash string, dateofbirth time.Time) *User {
	return &User{
		Id:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordhash,
		DateOfBirth:  dateofbirth,
		CreatedAt:    time.Now(),
		IsVerified:   false,
	}
}

func (user *User) ToUserGorm() *UserGorm {
	return &UserGorm{
		ID:           user.Id,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		LastOnline:   user.LastOnline,
		Bio:          user.Bio,
		DateOfBirth:  user.DateOfBirth,
		Location:     user.Location,
		IsVerified:   user.IsVerified,
	}
}
