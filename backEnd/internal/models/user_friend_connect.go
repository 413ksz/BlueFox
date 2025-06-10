package models

import (
	"time"

	"github.com/google/uuid"
)

type UserFriendConnect struct {
	UserID1     uuid.UUID  `json:"user_id_1" gorm:"primaryKey;type:uuid"`
	UserID2     uuid.UUID  `json:"user_id_2" gorm:"primaryKey;type:uuid"`
	Status      Status     `json:"status" gorm:"type:varchar(20)"`
	RequestedAt time.Time  `json:"requested_at" gorm:"autoCreateTime"`
	AcceptedAt  *time.Time `json:"accepted_at"`
	User1       User       `gorm:"foreignKey:UserID1;references:ID"`
	User2       User       `gorm:"foreignKey:UserID2;references:ID"`
}
