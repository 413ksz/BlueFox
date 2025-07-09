package users

import (
	"time"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/google/uuid"
)

// User table gorm model
type UserGorm struct {
	// Base Fields
	ID           uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid,default:gen_random_uuid()"`
	Username     string     `json:"username" gorm:"not null;Index"`
	Email        string     `json:"email" gorm:"size:254;not null;unique"`
	PasswordHash string     `json:"password_hash" gorm:"not null"`
	FirstName    *string    `json:"first_name,omitempty"`
	LastName     *string    `json:"last_name,omitempty"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	LastOnline   *time.Time `json:"last_online,omitempty"`
	Bio          *string    `json:"bio,omitempty"`
	DateOfBirth  time.Time  `json:"date_of_birth" gorm:"not null"`
	Location     *string    `json:"location,omitempty"`
	IsVerified   bool       `json:"is_verified" gorm:"default:false"`

	// Foreign Key for Profile Picture
	ProfilePictureAssetID *uuid.UUID `json:"profile_picture_asset_id" gorm:"type:uuid"`

	// Relations (Has One / Has Many)
	ProfilePictureAsset    *models.MediaAsset         `gorm:"foreignKey:ProfilePictureAssetID"` // Relation: A user has one profile picture
	SentMessages           []models.Message           `gorm:"foreignKey:AuthorID"`              // Relation: A user sends many messages
	UserFriendConnectsSent []models.UserFriendConnect `gorm:"foreignKey:User1ID"`               // Relation: A user initiates many friend connections
	UserFriendConnectsRecv []models.UserFriendConnect `gorm:"foreignKey:User2ID"`               // Relation: A user receives many friend connections
	OwnedServers           []models.Server            `gorm:"foreignKey:OwnerID"`               // Relation: A user owns many servers
	ServerUserConnects     []models.ServerUserConnect `gorm:"foreignKey:UserID"`                // Relation: A user is connected to many servers
	UploadedMediaAssets    []models.MediaAsset        `gorm:"foreignKey:UploadedByUserID"`      // Relation: A user uploads many media assets
}
