package api

import (
	"pomodoro-habit-backend/internal/db"
	"pomodoro-habit-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

// Save Session
func SavePomodoroSession(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type Request struct {
		Duration  int  `json:"duration"`
		Completed bool `json:"completed"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	session := models.PomodoroSession{
		UserID:    userID,
		Duration:  req.Duration,
		Completed: req.Completed,
	}

	if err := db.DB.Create(&session).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save session"})
	}

	return c.Status(fiber.StatusCreated).JSON(session)
}

// Get Stats (Simple aggregation)
func GetPomodoroStats(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Total Minutes
	var totalMinutes int64
	db.DB.Model(&models.PomodoroSession{}).
		Where("user_id = ? AND completed = ?", userID, true).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&totalMinutes)

	return c.JSON(fiber.Map{
		"total_minutes": totalMinutes,
	})
}
