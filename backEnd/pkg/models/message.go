package models

import (
	"time"

	"github.com/google/uuid"
)

// Message table gorm model
type Message struct {
	// Base Fields
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AuthorID    uuid.UUID   `gorm:"not null;type:uuid"`
	MessageType MessageType `gorm:"not null"`
	Content     string      `gorm:"not null;index:idx_content_type_search,priority:1"`
	CreatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   *time.Time  `gorm:"autoUpdateTime"`

	// Foreign Key for Reply
	ReplyTo *uuid.UUID `gorm:"type:uuid"` // Can be null if not a reply

	// Relations
	Author         User                `gorm:"foreignKey:AuthorID"`  // Relation: A message has one author
	ReplyToMessage *Message            `gorm:"foreignKey:ReplyTo"`   // Relation: A message can reply to another message
	Replies        []Message           `gorm:"foreignKey:ReplyTo"`   // Relation: A message can have many replies
	Attachments    []MessageAttachment `gorm:"foreignKey:MessageID"` // Relation: A message can have many attachments
}
