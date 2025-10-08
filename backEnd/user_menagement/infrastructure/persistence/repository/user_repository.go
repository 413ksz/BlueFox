package repository

import (
	"context"

	databaseerrorhelper "github.com/413ksz/BlueFox/backEnd/pkg/database_error_helper"
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/domain/model"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/domain/repository"
	db "github.com/413ksz/BlueFox/backEnd/user_menagement/infrastructure/persistence/database"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/infrastructure/persistence/mapper"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/shared/dto"
	"github.com/google/uuid"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	Querier db.Querier
}

func NewUserRepository(querier db.Querier) *UserRepository {
	return &UserRepository{
		Querier: querier,
	}
}

func (repository UserRepository) Create(user *model.User) *models.CustomError {

	userDbModel := db.CreateUserParams{
		ID:           user.Id,
		Email:        user.Email.String(),
		Username:     user.Username.String(),
		PasswordHash: user.PasswordHash.String(),
		DateOfBirth:  user.DateOfBirth.Time(),
	}

	ctx := context.Background()
	err := repository.Querier.CreateUser(ctx, userDbModel)

	if err != nil {
		databaseErr := databaseerrorhelper.GetDatabaseErrorMessage(err)
		return databaseErr
	}

	return nil
}

func (repository UserRepository) Get(id uuid.UUID) (*dto.PublicUserDTO, *models.CustomError) {
	ctx := context.Background()
	requestData, databaseError := repository.Querier.GetUser(ctx, id)

	if databaseError != nil {
		repoError := databaseerrorhelper.GetDatabaseErrorMessage(databaseError)
		return nil, repoError
	}

	user := mapper.ToPublicUserDto(&requestData)

	return user, nil
}

func (repository UserRepository) Update(id uuid.UUID, user *model.User) (*model.User, *models.CustomError) {
	return nil, nil
}

func (repository UserRepository) Delete(id uuid.UUID) *models.CustomError {
	return nil
}
