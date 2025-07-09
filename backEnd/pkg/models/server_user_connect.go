package models

import (
	"os/user"

	"github.com/google/uuid"
)

// ServerUserConnect table gorm model
type ServerUserConnect struct {
	// Composite Primary Keys (Foreign Keys)
	ServerID uuid.UUID `gorm:"not null;type:uuid;primaryKey;autoIncrement:false"`
	UserID   uuid.UUID `gorm:"not null;type:uuid;primaryKey;autoIncrement:false"`

	// Relations
	Server Server    `gorm:"foreignKey:ServerID"` // Relation: Connects to the server
	User   user.User `gorm:"foreignKey:UserID"`   // Relation: Connects to the user
}
