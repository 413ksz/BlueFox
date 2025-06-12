package models

import (
	"time"

	"github.com/google/uuid"
)

type Server struct {
	ID          uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       *string    `json:"title"`
	OwnerID     uuid.UUID  `json:"owner_id" gorm:"type:uuid"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	IconAssetID *uuid.UUID `json:"icon_asset_id" gorm:"type:uuid"`
	Visibility  Visibility `json:"visibility" gorm:"type:varchar(20)"`
	Owner       User       `gorm:"foreignKey:OwnerID;references:ID"`
	IconAsset   MediaAsset `gorm:"foreignKey:IconAssetID;references:ID"`
}
