package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                    uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username              string     `json:"username" gorm:"not null;uniqueIndex"`
	Email                 string     `json:"email" gorm:"not null;unique"`
	PasswordHash          string     `json:"password_hash" gorm:"not null"`
	ProfilePictureAssetID *uuid.UUID `json:"profile_picture_asset_id" gorm:"type:uuid"`
	CreatedAt             time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	LastOnline            *time.Time `json:"last_online"`
	Bio                   *string    `json:"bio"`
	DateOfBirth           time.Time  `json:"date_of_birth" gorm:"not null"`
	Location              *string    `json:"location"`
	IsVerified            bool       `json:"is_verified" gorm:"default:false"`
}
