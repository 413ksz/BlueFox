package models

import (
	"github.com/google/uuid"
)

type Channel struct {
	ID            uuid.UUID   `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ServerID      uuid.UUID   `json:"server_id" gorm:"not null;type:uuid"`
	Name          *string     `json:"name"`
	Icon          *string     `json:"icon"`
	Type          ChannelType `json:"type" gorm:"type:varchar(20)"`
	Topic         *string     `json:"topic"`
	Parent        *uuid.UUID  `json:"parent" gorm:"type:uuid"`
	Server        Server      `gorm:"foreignKey:ServerID;references:ID"`
	ParentChannel *Channel    `gorm:"foreignKey:Parent;references:ID"`
}
