package models

import (
	"github.com/google/uuid"
)

type MessageAttachment struct {
	ID           uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	MessageID    uuid.UUID  `json:"message_id" gorm:"uniqueIndex:idx_message_media;type:uuid"`
	MediaAssetID uuid.UUID  `json:"media_asset_id" gorm:"uniqueIndex:idx_message_media;type:uuid"`
	Message      Message    `gorm:"foreignKey:MessageID;references:ID"`
	MediaAsset   MediaAsset `gorm:"foreignKey:MediaAssetID;references:ID"`
}
