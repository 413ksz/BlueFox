package model

import (
	"time"

	"github.com/google/uuid"
)

// MediaAsset table gorm model
type MediaAsset struct {
	// Base Fields
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Filename  string    `gorm:"not null"`
	UrlPath   string    `gorm:"not null"` // URL path to the file
	FileSize  int
	MimeType  AssetType `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	// Foreign Key for Uploader (Optional)
	UploadedByUserID *uuid.UUID `gorm:"type:uuid"` // Optional: Track who uploaded it

	// Relations
	UploadedByUser      *User               `gorm:"foreignKey:UploadedByUserID"`      // Relation: A media asset can be uploaded by a user
	MessageAttachments  []MessageAttachment `gorm:"foreignKey:MediaAssetID"`          // Relation: A media asset can be part of many message attachments
	UserProfilePictures []User              `gorm:"foreignKey:ProfilePictureAssetID"` // Relation: A media asset can be a profile picture for multiple users
	ServerIcons         []Server            `gorm:"foreignKey:IconAssetID"`           // Relation: A media asset can be an icon for multiple servers
}
