package service

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	passwordHashing "github.com/413ksz/BlueFox/backEnd/pkg/password_hashing"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/application/command"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/application/query"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/domain/model"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/domain/repository"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/shared/dto"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/shared/validation"

	"github.com/google/uuid"
)

// UserService is an interface for user-related operations.
type UserService interface {
	CreateUser(command command.UserCreateCommand) *models.CustomError
	GetUser(id query.PublicUserQuery) (*dto.PublicUserDTO, *models.CustomError)
	UpdateUser(id uuid.UUID, command command.UserUpdateCommand) (*model.User, *models.CustomError)
	DeleteUser(id uuid.UUID) *models.CustomError
}

// UserServiceImpl implements the UserService interface.
// and encapsulates the business logic for user operations.
type UserServiceImpl struct {
	userRepo      repository.UserRepository // interface
	pwnedPassword *validation.PwnedPassword
	argon2ID      *passwordHashing.KriptoArgon2ID
}

// NewUserService creates a new instance of UserServiceImpl with the provided dependencies.
func NewUserService(userRepo repository.UserRepository, pwnedPassword *validation.PwnedPassword, argon2ID *passwordHashing.KriptoArgon2ID) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:      userRepo, // interface
		pwnedPassword: pwnedPassword,
		argon2ID:      argon2ID,
	}
}

// CreateUser creates a new user based on the provided command.
//
// Functionality:
// 1. HIBP Check: Verifies if the password has been exposed in a data breach using the HIBP API. If compromised, an error is returned.
// 2. Password Hashing: Securely hashes the password for storage in the database.
// 3. User Creation: Creates the user domain model.
// 4. User Repository: Persists the user data in the repository.
//
// Parameters:
// - command: The command containing the necessary data to create a new user.
//
// Returns:
//   - *models.CustomError: An error if the operation fails, otherwise nil.
//
// Error Conditions:
//   - External dependency error: An error occurs with the Pwned Passwords API (e.g., request timeout).
//   - Internal server error: An error occurs during a critical operation (e.g., failed to hash password).
//
// Example:
//  err := userHandler.UserService.CreateUser(command)
//  if err != nil {
//      return err
//  }
func (s *UserServiceImpl) CreateUser(command command.UserCreateCommand) *models.CustomError {

	if err := s.pwnedPassword.CheckPasswordBreach(command.Password.String()); err != nil {
		return err
	}

	passwordHash, hashingErr := s.argon2ID.GenerateNew(command.Password.String())
	if hashingErr != nil {
		return hashingErr
	}

	user, domainErr := model.NewUser(command.Username.String(), command.Email.String(), passwordHash, command.DateOfBirth.Time())
	if domainErr != nil {
		return domainErr
	}

	if repoError := s.userRepo.Create(user); repoError != nil {
		return repoError
	}

	return nil
}

func (s *UserServiceImpl) GetUser(query query.PublicUserQuery) (*dto.PublicUserDTO, *models.CustomError) {

	user, repoError := s.userRepo.Get(query.Id.UUID())
	if repoError != nil {
		return nil, repoError
	}

	return user, nil
}

func (s *UserServiceImpl) UpdateUser(id uuid.UUID, command command.UserUpdateCommand) (*model.User, *models.CustomError) {
	// implementation for UpdateUser method
	return nil, nil
}

func (s *UserServiceImpl) DeleteUser(id uuid.UUID) *models.CustomError {
	// implementation for DeleteUser method
	return nil
}
