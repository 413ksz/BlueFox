package users

import (
	databaseerrorhelper "github.com/413ksz/BlueFox/backEnd/pkg/database_error_helper"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryLayer interface {
	Create(user *User) *models.CustomError
	Get(id uuid.UUID) (*User, *models.CustomError)
	Update(id uuid.UUID, dto UserUpdateRequestDTO) (*User, *models.CustomError)
	Delate(id uuid.UUID) *models.CustomError
}

type UserRepository struct {
	dB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		dB: db,
	}
}

func (repository UserRepository) Create(user *User) *models.CustomError {

	result := repository.dB.Create(user)

	err := databaseerrorhelper.GetDatabaseErrorMessage(result)
	if err != nil {
		return err
	}

	return nil
}

func (repository UserRepository) Get(id uuid.UUID) (*User, *models.CustomError) {
	return nil, nil
}

func (repository UserRepository) Update(id uuid.UUID, dto UserUpdateRequestDTO) (*User, *models.CustomError) {
	return nil, nil
}

func (repository UserRepository) Delate(id uuid.UUID) *models.CustomError {
	return nil
}
