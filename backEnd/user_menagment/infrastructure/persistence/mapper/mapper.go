package mapper

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	domain "github.com/413ksz/BlueFox/backEnd/user_menagment/domain/model"
	gorm "github.com/413ksz/BlueFox/backEnd/user_menagment/infrastructure/persistence/model"
)

// ToUserGorm maps a domain.User to a gorm.UserGorm
// parameters:
//   - user: the domain.User to map
//
// returns:
//   - gorm.UserGorm: the mapped gorm.UserGorm
func ToUserGorm(user *domain.User) *gorm.UserGorm {
	return &gorm.UserGorm{
		ID:           user.Id,
		Username:     user.Username.String(),
		Email:        user.Email.String(),
		PasswordHash: user.PasswordHash.String(),
		DateOfBirth:  user.DateOfBirth.Time(),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		LastOnline:   user.LastOnline,
		Bio:          user.Bio,
		Location:     user.Location,
		IsVerified:   user.IsVerified.Bool(),
	}
}

func ToUserDomain(user *gorm.UserGorm) (*domain.User, *models.CustomError) {
	domain, domainErr := domain.NewUser(user.Username, user.Email, user.PasswordHash, user.DateOfBirth
	return domain, domainErr
}
