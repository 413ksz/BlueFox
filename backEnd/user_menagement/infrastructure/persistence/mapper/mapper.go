package mapper

import (
	db "github.com/413ksz/BlueFox/backEnd/user_menagement/infrastructure/persistence/database"
	"github.com/413ksz/BlueFox/backEnd/user_menagement/shared/dto"
	"github.com/google/uuid"
)

// ToPublicUserDto maps a database row (db.GetUserRow)
// which may contain nullable fields, to a
// non-nullable PublicUserDTO for public-facing data
// presentation.
// parameters:
// - dbRow: the database row to map
// returns:
// - *dto.PublicUserDTO: the mapped PublicUserDTO
// Example:
// userDto := mapper.ToPublicUserDto(dbRow)
func ToPublicUserDto(dbRow *db.GetUserRow) *dto.PublicUserDTO {
	var bioString *string
	if dbRow.Bio.Valid {
		bioString = &dbRow.Bio.String
	}
	var profilePictureAssetID *uuid.UUID
	if dbRow.ProfilePictureAssetID.Valid {
		bytes := uuid.UUID(dbRow.ProfilePictureAssetID.Bytes)
		profilePictureAssetID = &bytes
	}
	return dto.NewPublicUserDTO(dbRow.ID, dbRow.Username, bioString, profilePictureAssetID)
}
