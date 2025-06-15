package models

import (
	"time"

	"github.com/google/uuid"
)

// Server table gorm model
type Server struct {
	// Base Fields
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title      string     `gorm:"index:idx_title_visibility_search,priority:1"`
	OwnerID    uuid.UUID  `gorm:"type:uuid"`
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	Visibility Visibility `gorm:"index:idx_title_visibility_search,priority:2"`

	// Foreign Key for Icon
	IconAssetID *uuid.UUID `gorm:"type:uuid"`

	// Relations
	Owner       User                `gorm:"foreignKey:OwnerID"`     // Relation: A server has one owner
	IconAsset   *MediaAsset         `gorm:"foreignKey:IconAssetID"` // Relation: A server has one icon asset
	Channels    []Channel           `gorm:"foreignKey:ServerID"`    // Relation: A server has many channels
	ServerUsers []ServerUserConnect `gorm:"foreignKey:ServerID"`    // Relation: A server has many connected users
}
