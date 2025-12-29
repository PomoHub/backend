package api

import (
	"pomodoro-habit-backend/internal/db"
	"pomodoro-habit-backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// SendFriendRequest sends a friend request to a user
func SendFriendRequest(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	targetId := c.Params("userId")

	if userId == targetId {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot send friend request to yourself"})
	}

	var existing models.Friend
	result := db.DB.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userId, targetId, targetId, userId).First(&existing)
	
	if result.RowsAffected > 0 {
		if existing.Status == models.FriendStatusBlocked {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Cannot send request"})
		}
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Friend request already exists or you are already friends"})
	}

	friendRequest := models.Friend{
		UserID:    uuid.MustParse(userId),
		FriendID:  uuid.MustParse(targetId),
		Status:    models.FriendStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.DB.Create(&friendRequest).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send friend request"})
	}

	return c.Status(fiber.StatusCreated).JSON(friendRequest)
}

// AcceptFriendRequest accepts a received friend request
func AcceptFriendRequest(c *fiber.Ctx) error {
	requestId := c.Params("requestId")
	userId := c.Locals("user_id").(string)

	var request models.Friend
	if err := db.DB.First(&request, "id = ?", requestId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Request not found"})
	}

	// Verify the request was sent TO this user
	if request.FriendID.String() != userId {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not authorized to accept this request"})
	}

	request.Status = models.FriendStatusAccepted
	request.UpdatedAt = time.Now()

	if err := db.DB.Save(&request).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to accept request"})
	}

	return c.JSON(request)
}

// RejectFriendRequest rejects/deletes a friend request
func RejectFriendRequest(c *fiber.Ctx) error {
	requestId := c.Params("requestId")
	userId := c.Locals("user_id").(string)

	var request models.Friend
	if err := db.DB.First(&request, "id = ?", requestId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Request not found"})
	}

	// Verify ownership (either sender can cancel or receiver can reject)
	if request.UserID.String() != userId && request.FriendID.String() != userId {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not authorized"})
	}

	if err := db.DB.Delete(&request).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove request"})
	}

	return c.JSON(fiber.Map{"message": "Friend request removed"})
}

// RemoveFriend removes a friend connection
func RemoveFriend(c *fiber.Ctx) error {
	friendId := c.Params("friendId")
	userId := c.Locals("user_id").(string)

	// Delete record where (user=me AND friend=them) OR (user=them AND friend=me)
	result := db.DB.Where(
		"((user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)) AND status = ?", 
		userId, friendId, friendId, userId, models.FriendStatusAccepted,
	).Delete(&models.Friend{})

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove friend"})
	}

	return c.JSON(fiber.Map{"message": "Friend removed"})
}

// BlockUser blocks a user
func BlockUser(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	targetId := c.Params("userId")

	// Check if any relationship exists and update/overwrite it
	var relationship models.Friend
	result := db.DB.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userId, targetId, targetId, userId).First(&relationship)

	if result.RowsAffected > 0 {
		// Update existing
		relationship.Status = models.FriendStatusBlocked
		relationship.UserID = uuid.MustParse(userId) // Ensure I am the one blocking
		relationship.FriendID = uuid.MustParse(targetId)
		relationship.UpdatedAt = time.Now()
		db.DB.Save(&relationship)
	} else {
		// Create new block
		block := models.Friend{
			UserID:    uuid.MustParse(userId),
			FriendID:  uuid.MustParse(targetId),
			Status:    models.FriendStatusBlocked,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		db.DB.Create(&block)
	}

	return c.JSON(fiber.Map{"message": "User blocked"})
}

// GetFriends returns list of accepted friends
func GetFriends(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)

	var friends []models.Friend
	db.DB.Preload("User").Preload("Friend").Where(
		"((user_id = ? OR friend_id = ?) AND status = ?)", 
		userId, userId, models.FriendStatusAccepted,
	).Find(&friends)

	// Format response to just return the "other" user
	var friendUsers []models.User
	for _, f := range friends {
		if f.UserID.String() == userId {
			friendUsers = append(friendUsers, f.Friend)
		} else {
			friendUsers = append(friendUsers, f.User)
		}
	}

	return c.JSON(friendUsers)
}

// GetFriendRequests returns pending requests received by user
func GetFriendRequests(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)

	var requests []models.Friend
	db.DB.Preload("User").Where("friend_id = ? AND status = ?", userId, models.FriendStatusPending).Find(&requests)

	return c.JSON(requests)
}
