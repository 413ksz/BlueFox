package service

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/apierrors"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	passwordHashing "github.com/413ksz/BlueFox/backEnd/pkg/password_hashing"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/application/command"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/domain/model"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/domain/repository"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(command command.UserCreateCommand) *models.CustomError
	GetUser(id uuid.UUID) (*model.User, *models.CustomError)
	UpdateUser(id uuid.UUID, command command.UserUpdateCommand) (*model.User, *models.CustomError)
	DeleteUser(id uuid.UUID) *models.CustomError
}

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (s *UserServiceImpl) CreateUser(command command.UserCreateCommand) *models.CustomError {

	// Hash the password and check for errors
	passwordHash, err := passwordHashing.HashPassword(command.Password)
	if err != nil {
		customError := apierrors.ERROR_CODE_INTERNAL_SERVER.NewApiError("error hashing password", err)
		return customError
	}

	// Create a new User(domain) instance from the command and hash the password
	// pointer receiver
	user, domainErr := model.NewUser(command.Username, command.Email, passwordHash, command.DateOfBirth)
	if domainErr != nil {
		return domainErr
	}

	// Call the repository layer to create the user and check for errors
	repoError := s.userRepo.Create(user)
	if repoError != nil {
		return repoError
	}

	return nil
}

func (s *UserServiceImpl) GetUser(id uuid.UUID) (*model.User, *models.CustomError) {
	// implementation for GetUser method
	return nil, nil
}

func (s *UserServiceImpl) UpdateUser(id uuid.UUID, command command.UserUpdateCommand) (*model.User, *models.CustomError) {
	// implementation for UpdateUser method
	return nil, nil
}

func (s *UserServiceImpl) DeleteUser(id uuid.UUID) *models.CustomError {
	// implementation for DeleteUser method
	return nil
}
