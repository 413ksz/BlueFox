package models

import (
	"github.com/google/uuid"
)

type ServerUserConnect struct {
	ServerID uuid.UUID `json:"server_id" gorm:"primaryKey;type:uuid"`
	UserID   uuid.UUID `json:"user_id" gorm:"primaryKey;type:uuid"`
	Server   Server    `gorm:"foreignKey:ServerID;references:ID"`
	User     User      `gorm:"foreignKey:UserID;references:ID"`
}
