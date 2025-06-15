package models

import (
	"github.com/google/uuid"
)

type Channel struct {
	// Base Fields
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ServerID uuid.UUID `gorm:"not null;type:uuid"`
	Name     string
	Icon     *string
	Type     ChannelType
	Topic    *string

	// Foreign Key for Parent Channel (for nested channels/categories)
	Parent *uuid.UUID `gorm:"type:uuid"` // Can be null for top-level channels

	// Relations
	Server        Server    `gorm:"foreignKey:ServerID"` // Relation: A channel belongs to one server
	ParentChannel *Channel  `gorm:"foreignKey:Parent"`   // Relation: A channel can have a parent channel
	ChildChannels []Channel `gorm:"foreignKey:Parent"`   // Relation: A channel can have many child channels
}
