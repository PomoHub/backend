package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Habit struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Title     string         `gorm:"not null" json:"title"`
	Frequency string         `gorm:"default:'daily'" json:"frequency"` // daily, weekly
	Color     string         `json:"color"`
	Emoji     string         `json:"emoji"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Logs      []HabitLog     `json:"logs"`
}

type HabitLog struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	HabitID   uuid.UUID `gorm:"type:uuid;not null;index" json:"habit_id"`
	Date      string    `gorm:"type:date;not null" json:"date"` // YYYY-MM-DD
	Completed bool      `gorm:"default:true" json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Habit) BeforeCreate(tx *gorm.DB) (err error) {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return
}

func (hl *HabitLog) BeforeCreate(tx *gorm.DB) (err error) {
	if hl.ID == uuid.Nil {
		hl.ID = uuid.New()
	}
	return
}
