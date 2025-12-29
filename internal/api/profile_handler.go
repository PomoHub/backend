package api

import (
	"pomodoro-habit-backend/internal/db"
	"pomodoro-habit-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Update Profile
func UpdateProfile(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type Request struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Bio       string `json:"bio"`
		AvatarURL string `json:"avatar_url"`
		BannerURL string `json:"banner_url"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.BannerURL != "" {
		user.BannerURL = req.BannerURL
	}

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update profile"})
	}

	return c.JSON(user)
}

// Get User Profile (Public)
func GetUserProfile(c *fiber.Ctx) error {
	username := c.Params("username")

	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Don't expose sensitive data, although User struct json tags handle some, let's be safe if we add more
	return c.JSON(user)
}

// Create Post
func CreatePost(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type Request struct {
		Content  string `json:"content"`
		ImageURL string `json:"image_url"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Content == "" && req.ImageURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Content or Image is required"})
	}

	post := models.Post{
		UserID:   userID,
		Content:  req.Content,
		ImageURL: req.ImageURL,
	}

	if err := db.DB.Create(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create post"})
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

// Get User Posts
func GetUserPosts(c *fiber.Ctx) error {
	// Can get by userID or username. Let's support username to match profile view.
	username := c.Params("username")

	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	var posts []models.Post
	if err := db.DB.Where("user_id = ?", user.ID).Order("created_at desc").Find(&posts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch posts"})
	}

	return c.JSON(posts)
}

// Delete Post
func DeletePost(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	postIDStr := c.Params("id")
	postID, _ := uuid.Parse(postIDStr)

	if err := db.DB.Where("id = ? AND user_id = ?", postID, userID).Delete(&models.Post{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete post"})
	}

	return c.JSON(fiber.Map{"message": "Post deleted successfully"})
}
