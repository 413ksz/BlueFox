package query

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	valueobject "github.com/413ksz/BlueFox/backEnd/pkg/shared/value_object"
)

// PublicUserQuerry is a query to get a public user by id
type PublicUserQuery struct {
	Id valueobject.ID `json:"id"`
}

func NewPublicUserQuery(id string) (PublicUserQuery, *models.CustomError) {
	uuid, validationError := valueobject.NewUUID(id)
	if validationError != nil {
		return PublicUserQuery{}, models.NewCustomError(models.ERROR_CODE_BAD_REQUEST, validationError, nil, nil)

	}
	return PublicUserQuery{Id: uuid}, nil
}
