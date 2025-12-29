package api

import (
	"pomodoro-habit-backend/internal/db"
	"pomodoro-habit-backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Get Habits
func GetHabits(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var habits []models.Habit
	if err := db.DB.Preload("Logs").Where("user_id = ?", userID).Find(&habits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch habits"})
	}

	return c.JSON(habits)
}

// Create Habit
func CreateHabit(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type Request struct {
		Title     string `json:"title"`
		Frequency string `json:"frequency"`
		Color     string `json:"color"`
		Emoji     string `json:"emoji"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	habit := models.Habit{
		UserID:    userID,
		Title:     req.Title,
		Frequency: req.Frequency,
		Color:     req.Color,
		Emoji:     req.Emoji,
	}

	if err := db.DB.Create(&habit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create habit"})
	}

	return c.Status(fiber.StatusCreated).JSON(habit)
}

// Toggle Habit (Log Completion)
func ToggleHabit(c *fiber.Ctx) error {
	// Simple implementation: If logged for today, delete log. Else create log.
	// For robust sync, we might want explicit "log" endpoints.
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	habitIDStr := c.Params("id")
	habitID, _ := uuid.Parse(habitIDStr)

	// Verify ownership
	var habit models.Habit
	if err := db.DB.Where("id = ? AND user_id = ?", habitID, userID).First(&habit).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Habit not found"})
	}

	today := time.Now().Format("2006-01-02")
	var log models.HabitLog
	result := db.DB.Where("habit_id = ? AND date = ?", habitID, today).First(&log)

	if result.RowsAffected > 0 {
		// Already completed today, un-complete it
		db.DB.Delete(&log)
		return c.JSON(fiber.Map{"status": "uncompleted"})
	} else {
		// Complete it
		newLog := models.HabitLog{
			HabitID: habitID,
			Date:    today,
		}
		db.DB.Create(&newLog)
		return c.JSON(fiber.Map{"status": "completed"})
	}
}

// Delete Habit
func DeleteHabit(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	habitIDStr := c.Params("id")
	habitID, _ := uuid.Parse(habitIDStr)

	if err := db.DB.Where("id = ? AND user_id = ?", habitID, userID).Delete(&models.Habit{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete habit"})
	}

	return c.JSON(fiber.Map{"message": "Habit deleted successfully"})
}
