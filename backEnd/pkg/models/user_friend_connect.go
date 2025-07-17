package models

import (
	"time"

	"github.com/google/uuid"
)

// UserFriendConnect table gorm model
type UserFriendConnect struct {
	// Composite Primary Keys (Foreign Keys)
	User1ID uuid.UUID `gorm:"not null;type:uuid;primaryKey;autoIncrement:false"`
	User2ID uuid.UUID `gorm:"not null;type:uuid;primaryKey;autoIncrement:false"`

	// Base Fields
	Status      Status
	RequestedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	AcceptedAt  *time.Time

	// Relations
	User1 UserGorm `gorm:"foreignKey:User1ID"` // Relation: Connects to the first user in the friendship
	User2 UserGorm `gorm:"foreignKey:User2ID"` // Relation: Connects to the second user in the friendship
}
