package api

import (
	"pomodoro-habit-backend/internal/db"
	"pomodoro-habit-backend/internal/models"
	"pomodoro-habit-backend/internal/ws"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Send Message
func SendMessage(c *fiber.Ctx) error {
	senderID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	spaceIDStr := c.Params("spaceId")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid space ID"})
	}

	type Request struct {
		Content string `json:"content"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Content cannot be empty"})
	}

	// Check if user is member of space
	var membership models.SpaceMember
	if err := db.DB.Where("space_id = ? AND user_id = ?", spaceID, senderID).First(&membership).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not a member of this space"})
	}

	message := models.Message{
		SpaceID:   spaceID,
		SenderID:  senderID,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.DB.Create(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not send message"})
	}

	// Broadcast via WebSocket
	ws.GlobalHub.BroadcastToSpace(spaceID, ws.TypeChatMessage, message)

	return c.Status(fiber.StatusCreated).JSON(message)
}

// Get Messages
func GetMessages(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	spaceIDStr := c.Params("spaceId")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid space ID"})
	}

	// Check membership
	var membership models.SpaceMember
	if err := db.DB.Where("space_id = ? AND user_id = ?", spaceID, userID).First(&membership).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	var messages []models.Message
	// Pagination could be added here
	err = db.DB.Where("space_id = ?", spaceID).
		Order("created_at desc").
		Limit(50).
		Preload("Sender"). // Load sender details
		Find(&messages).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch messages"})
	}

	return c.JSON(messages)
}
