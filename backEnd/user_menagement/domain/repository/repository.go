package repository

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/domain/model"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/shared/dto"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *model.User) *models.CustomError
	Get(id uuid.UUID) (*dto.PublicUserDTO, *models.CustomError)
	Update(id uuid.UUID, user *model.User) (*model.User, *models.CustomError)
	Delete(id uuid.UUID) *models.CustomError
}
