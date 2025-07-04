package model

import (
	"github.com/google/uuid"
)

// MessageAttachment table gorm model
type MessageAttachment struct {
	// Base Fields
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	// Composite Unique Index (Foreign Keys)
	MessageID    uuid.UUID `gorm:"not null;type:uuid;uniqueIndex:idx_message_media_unique,priority:1"`
	MediaAssetID uuid.UUID `gorm:"not null;type:uuid;uniqueIndex:idx_message_media_unique,priority:2"`

	// Relations
	Message    Message    `gorm:"foreignKey:MessageID"`    // Relation: An attachment belongs to one message
	MediaAsset MediaAsset `gorm:"foreignKey:MediaAssetID"` // Relation: An attachment links to one media asset
}
