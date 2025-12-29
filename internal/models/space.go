package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Space struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	OwnerID   uuid.UUID      `gorm:"type:uuid;not null" json:"owner_id"`
	Owner     User           `gorm:"foreignKey:OwnerID" json:"owner"`
	Members   []SpaceMember  `json:"members"`
	
	// Pomodoro Settings
	PomodoroWorkDuration      int `gorm:"default:25" json:"pomodoro_work_duration"`
	PomodoroShortBreakDuration int `gorm:"default:5" json:"pomodoro_short_break_duration"`
	PomodoroLongBreakDuration  int `gorm:"default:15" json:"pomodoro_long_break_duration"`
	PomodoroRounds            int `gorm:"default:4" json:"pomodoro_rounds"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type SpaceMember struct {
	SpaceID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"space_id"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Role      string    `gorm:"default:'member'" json:"role"` // 'admin' or 'member'
	JoinedAt  time.Time `json:"joined_at"`
	Space     Space     `gorm:"foreignKey:SpaceID" json:"-"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

func (s *Space) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
