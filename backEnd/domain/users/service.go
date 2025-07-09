package users

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/apierrors"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	passwordHashing "github.com/413ksz/BlueFox/backEnd/pkg/password_hashing"
	"github.com/google/uuid"
)

type UserServiceLayer interface {
	CreateUser(dto UserCreateRequestDTO) *models.CustomError
	GetUser(id uuid.UUID) (*User, *models.CustomError)
	UpdateUser(id uuid.UUID, dto UserUpdateRequestDTO) (*User, *models.CustomError)
	DeleteUser(id uuid.UUID) *models.CustomError
}

type UserService struct {
	userRepo UserRepositoryLayer
}

func NewUserService(userRepo UserRepositoryLayer) UserServiceLayer {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(dto UserCreateRequestDTO) *models.CustomError {

	// Hash the password and check for errors
	passwordHash, err := passwordHashing.HashPassword(dto.Password)
	if err != nil {
		customError := apierrors.ERROR_CODE_INTERNAL_SERVER.NewApiError("error hashing password", err)
		return customError
	}

	// Create a new User(domain) instance from the DTO and hash the password
	// pointer receiver
	user := NewUser(dto.Username, dto.Email, passwordHash, dto.DateOfBirth)

	// Call the repository layer to create the user and check for errors
	customError := s.userRepo.Create(user)
	if customError != nil {
		return customError
	}

	return nil
}

func (s *UserService) GetUser(id uuid.UUID) (*User, *models.CustomError) {
	// implementation for GetUser method
	return nil, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, dto UserUpdateRequestDTO) (*User, *models.CustomError) {
	// implementation for UpdateUser method
	return nil, nil
}

func (s *UserService) DeleteUser(id uuid.UUID) *models.CustomError {
	// implementation for DeleteUser method
	return nil
}
