package models

import (
	"time"

	"github.com/google/uuid"
)

type FriendStatus string

const (
	FriendStatusPending  FriendStatus = "pending"
	FriendStatusAccepted FriendStatus = "accepted"
	FriendStatusBlocked  FriendStatus = "blocked"
)

type Friend struct {
	ID        uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
	FriendID  uuid.UUID    `gorm:"type:uuid;not null;index" json:"friend_id"`
	Status    FriendStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	User   User `gorm:"foreignKey:UserID" json:"user"`
	Friend User `gorm:"foreignKey:FriendID" json:"friend"`
}

// Composite unique index to prevent duplicate requests
func (Friend) TableName() string {
	return "friends"
}
