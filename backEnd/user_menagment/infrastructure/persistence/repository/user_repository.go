package repository

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/database"
	databaseerrorhelper "github.com/413ksz/BlueFox/backEnd/pkg/database_error_helper"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/domain/model"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/domain/repository"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/infrastructure/persistence/mapper"
	"github.com/google/uuid"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	dB *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{
		dB: db,
	}
}

func (repository UserRepository) Create(user *model.User) *models.CustomError {
	userGorm := mapper.ToUserGorm(user)
	result := repository.dB.DB.Create(&userGorm)

	databaseErr := databaseerrorhelper.GetDatabaseErrorMessage(result)
	if databaseErr != nil {
		return databaseErr
	}

	return nil
}

func (repository UserRepository) Get(id uuid.UUID) (*model.User, *models.CustomError) {
	return nil, nil
}

func (repository UserRepository) Update(id uuid.UUID, user *model.User) (*model.User, *models.CustomError) {
	return nil, nil
}

func (repository UserRepository) Delete(id uuid.UUID) *models.CustomError {
	return nil
}
