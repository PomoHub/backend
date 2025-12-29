package api

import (
	"pomodoro-habit-backend/internal/db"
	"pomodoro-habit-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Update Space Settings
func UpdateSpace(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	spaceIDStr := c.Params("spaceId")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid space ID"})
	}

	type Request struct {
		Name                       string `json:"name"`
		PomodoroWorkDuration       int    `json:"pomodoro_work_duration"`
		PomodoroShortBreakDuration int    `json:"pomodoro_short_break_duration"`
		PomodoroLongBreakDuration  int    `json:"pomodoro_long_break_duration"`
		PomodoroRounds             int    `json:"pomodoro_rounds"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var space models.Space
	if err := db.DB.First(&space, spaceID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Space not found"})
	}

	// Only owner or admin can update settings? Let's say Admin.
	var membership models.SpaceMember
	if err := db.DB.Where("space_id = ? AND user_id = ?", spaceID, userID).First(&membership).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not a member of this space"})
	}
	if membership.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only admins can update space settings"})
	}

	// Update fields if provided
	if req.Name != "" {
		space.Name = req.Name
	}
	if req.PomodoroWorkDuration > 0 {
		space.PomodoroWorkDuration = req.PomodoroWorkDuration
	}
	if req.PomodoroShortBreakDuration > 0 {
		space.PomodoroShortBreakDuration = req.PomodoroShortBreakDuration
	}
	if req.PomodoroLongBreakDuration > 0 {
		space.PomodoroLongBreakDuration = req.PomodoroLongBreakDuration
	}
	if req.PomodoroRounds > 0 {
		space.PomodoroRounds = req.PomodoroRounds
	}

	if err := db.DB.Save(&space).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update space"})
	}

	return c.JSON(space)
}
