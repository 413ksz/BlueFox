package models

import (
	"time"

	"github.com/google/uuid"
)

// User table gorm model
type User struct {
	// Base Fields
	ID          uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username    string     `json:"username" gorm:"not null;Index"`
	Email       string     `json:"email" gorm:"not null;unique"`
	Password    string     `json:"password_hash" gorm:"not null"`
	FirstName   string     `json:"first_name" gorm:"not null"`
	LastName    string     `json:"last_name" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	LastOnline  *time.Time `json:"last_online"`
	Bio         *string    `json:"bio"`
	DateOfBirth time.Time  `json:"date_of_birth" gorm:"not null"`
	Location    *string    `json:"location"`
	IsVerified  bool       `json:"is_verified" gorm:"default:false"`

	// Foreign Key for Profile Picture
	ProfilePictureAssetID *uuid.UUID `json:"profile_picture_asset_id" gorm:"type:uuid"`

	// Relations (Has One / Has Many)
	ProfilePictureAsset    *MediaAsset         `gorm:"foreignKey:ProfilePictureAssetID"` // Relation: A user has one profile picture
	SentMessages           []Message           `gorm:"foreignKey:AuthorID"`              // Relation: A user sends many messages
	UserFriendConnectsSent []UserFriendConnect `gorm:"foreignKey:User1ID"`               // Relation: A user initiates many friend connections
	UserFriendConnectsRecv []UserFriendConnect `gorm:"foreignKey:User2ID"`               // Relation: A user receives many friend connections
	OwnedServers           []Server            `gorm:"foreignKey:OwnerID"`               // Relation: A user owns many servers
	ServerUserConnects     []ServerUserConnect `gorm:"foreignKey:UserID"`                // Relation: A user is connected to many servers
	UploadedMediaAssets    []MediaAsset        `gorm:"foreignKey:UploadedByUserID"`      // Relation: A user uploads many media assets
}
