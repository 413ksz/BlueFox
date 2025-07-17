package repository

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/413ksz/BlueFox/backEnd/user_menagment/domain/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *model.User) *models.CustomError
	Get(id uuid.UUID) (*model.User, *models.CustomError)
	Update(id uuid.UUID, user *model.User) (*model.User, *models.CustomError)
	Delete(id uuid.UUID) *models.CustomError
}
