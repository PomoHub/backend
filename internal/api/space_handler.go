package api

import (
	"pomodoro-habit-backend/internal/db"
	"pomodoro-habit-backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Middleware to extract user_id from JWT
func getUserID(c *fiber.Ctx) (uuid.UUID, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	return uuid.Parse(userIDStr)
}

// Create Space
func CreateSpace(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type Request struct {
		Name string `json:"name"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Space name is required"})
	}

	// Check Space Limit (Max 3 owned spaces)
	var count int64
	db.DB.Model(&models.Space{}).Where("owner_id = ?", userID).Count(&count)
	if count >= 3 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You can only create up to 3 spaces"})
	}

	space := models.Space{
		Name:    req.Name,
		OwnerID: userID,
	}

	tx := db.DB.Begin()

	if err := tx.Create(&space).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create space"})
	}

	// Add owner as a member (admin)
	member := models.SpaceMember{
		SpaceID:  space.ID,
		UserID:   userID,
		Role:     "admin",
		JoinedAt: time.Now(),
	}

	if err := tx.Create(&member).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not add owner to space"})
	}

	tx.Commit()

	return c.Status(fiber.StatusCreated).JSON(space)
}

// Get My Spaces
func GetMySpaces(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var spaces []models.Space
	// Find spaces where user is a member
	err = db.DB.Joins("JOIN space_members ON space_members.space_id = spaces.id").
		Where("space_members.user_id = ?", userID).
		Preload("Members").
		Preload("Members.User"). // Load user details for members
		Find(&spaces).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch spaces"})
	}

	return c.JSON(spaces)
}

// Get Space Details
func GetSpaceDetails(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	spaceIDStr := c.Params("spaceId")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid space ID"})
	}

	var space models.Space
	// Check if user is a member
	var membership models.SpaceMember
	if err := db.DB.Where("space_id = ? AND user_id = ?", spaceID, userID).First(&membership).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not a member of this space"})
	}

	if err := db.DB.Preload("Members").Preload("Members.User").First(&space, spaceID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Space not found"})
	}

	return c.JSON(space)
}

// Add Member to Space (by User ID)
func AddMember(c *fiber.Ctx) error {
	currentUserID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	spaceIDStr := c.Params("spaceId")
	spaceID, err := uuid.Parse(spaceIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid space ID"})
	}

	type Request struct {
		UserID uuid.UUID `json:"user_id"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check if current user is admin/owner of the space
	var membership models.SpaceMember
	if err := db.DB.Where("space_id = ? AND user_id = ?", spaceID, currentUserID).First(&membership).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not a member of this space"})
	}
	if membership.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only admins can add members"})
	}

	// Check if target user exists
	var targetUser models.User
	if err := db.DB.First(&targetUser, req.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if already a member
	var existingMember models.SpaceMember
	if result := db.DB.Where("space_id = ? AND user_id = ?", spaceID, req.UserID).First(&existingMember); result.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User is already a member"})
	}

	// Add Member
	newMember := models.SpaceMember{
		SpaceID:  spaceID,
		UserID:   req.UserID,
		Role:     "member",
		JoinedAt: time.Now(),
	}

	if err := db.DB.Create(&newMember).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not add member"})
	}

	return c.JSON(fiber.Map{"message": "Member added successfully", "member": newMember})
}

// Leave Space / Remove Member
func RemoveMember(c *fiber.Ctx) error {
	currentUserID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	spaceIDStr := c.Params("spaceId")
	userIDStr := c.Params("userId")

	spaceID, _ := uuid.Parse(spaceIDStr)
	targetUserID, _ := uuid.Parse(userIDStr)

	// Check permissions
	var currentUserMembership models.SpaceMember
	if err := db.DB.Where("space_id = ? AND user_id = ?", spaceID, currentUserID).First(&currentUserMembership).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	// Allow if removing self OR if admin removing someone else
	if currentUserID != targetUserID && currentUserMembership.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only admins can remove other members"})
	}

	if err := db.DB.Delete(&models.SpaceMember{}, "space_id = ? AND user_id = ?", spaceID, targetUserID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not remove member"})
	}

	return c.JSON(fiber.Map{"message": "Member removed successfully"})
}

// Delete Space
func DeleteSpace(c *fiber.Ctx) error {
	currentUserID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	spaceIDStr := c.Params("spaceId")
	spaceID, _ := uuid.Parse(spaceIDStr)

	var space models.Space
	if err := db.DB.First(&space, spaceID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Space not found"})
	}

	if space.OwnerID != currentUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only the owner can delete the space"})
	}

	// Delete space (Cascade delete should handle members/messages if configured in DB, but GORM soft delete might need manual handling or constraint setup. For now, simple delete)
	// Ideally we should delete members first or rely on DB FK constraints ON DELETE CASCADE

	// Transaction
	tx := db.DB.Begin()
	if err := tx.Delete(&models.SpaceMember{}, "space_id = ?", spaceID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete members"})
	}
	if err := tx.Delete(&space).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete space"})
	}
	tx.Commit()

	return c.JSON(fiber.Map{"message": "Space deleted successfully"})
}
