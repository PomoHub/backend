package db

import (
	"log"
	"pomodoro-habit-backend/internal/config"
	"pomodoro-habit-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to Database successfully")

	// Auto Migrate
	log.Println("Running Migrations...")
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Space{},
		&models.SpaceMember{},
		&models.Message{},
		&models.Todo{},
		&models.Habit{},
		&models.HabitLog{},
		&models.PomodoroSession{},
		&models.Post{},
	); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")
}
