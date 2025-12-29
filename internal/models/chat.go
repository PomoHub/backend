package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SpaceID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"space_id"`
	SenderID  uuid.UUID      `gorm:"type:uuid;not null" json:"sender_id"`
	Content   string         `gorm:"not null" json:"content"`
	Sender    User           `gorm:"foreignKey:SenderID" json:"sender"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}
