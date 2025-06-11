package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID   `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthorID       uuid.UUID   `json:"author_id" gorm:"not null;type:uuid"`
	MessageType    MessageType `json:"message_type" gorm:"type:varchar(20);not null"`
	Content        string      `json:"content" gorm:"not null"`
	CreatedAt      time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      *time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	ReplyTo        *uuid.UUID  `json:"reply_to" gorm:"type:uuid"`
	Author         User        `gorm:"foreignKey:AuthorID;references:ID"`
	ReplyToMessage *Message    `gorm:"foreignKey:ReplyTo;references:ID"`
}
