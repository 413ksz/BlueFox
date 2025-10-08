package dto

import "github.com/google/uuid"

// PublicUserDTO defines the structure for a public facing user data.
type PublicUserDTO struct {
	Id                    uuid.UUID  `json:"id"`
	Username              string     `json:"username"`
	Bio                   *string    `json:"bio,omitempty"`
	ProfilePictureAssetID *uuid.UUID `json:"profilePictureAssetId,omitempty"`
}

// NewPublicUserDTO creates a new PublicUserDTO instance
func NewPublicUserDTO(id uuid.UUID, username string, bio *string, profilePictureAssetID *uuid.UUID) *PublicUserDTO {
	return &PublicUserDTO{
		Id:                    id,
		Username:              username,
		Bio:                   bio,
		ProfilePictureAssetID: profilePictureAssetID,
	}
}
