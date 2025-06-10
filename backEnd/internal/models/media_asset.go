package models

import (
	"time"

	"github.com/google/uuid"
)

type MediaAsset struct {
	ID               uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Filename         string     `json:"filename" gorm:"not null"`
	FilePath         string     `json:"file_path" gorm:"not null"`
	FileSize         int        `json:"file_size"`
	MimeType         AssetType  `json:"mime_type" gorm:"type:varchar(50);not null"`
	AssetType        string     `json:"asset_type" gorm:"not null"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UploadedByUserID *uuid.UUID `json:"uploaded_by_user_id" gorm:"type:uuid"`
	UploadedByUser   User       `gorm:"foreignKey:UploadedByUserID;references:ID"`
}
