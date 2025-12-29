package api

import (
	"pomodoro-habit-backend/internal/config"
	"pomodoro-habit-backend/internal/db"
	"pomodoro-habit-backend/internal/models"
	"pomodoro-habit-backend/internal/utils"
	"regexp"
	"time"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func validatePassword(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasSpecial = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasSpecial
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate (basic)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// Username Validation (English letters only)
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(req.Username) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username must contain only English letters, numbers, and underscores"})
	}

	// Password Validation
	if !validatePassword(req.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one special character"})
	}

	// Check if user exists
	var existingUser models.User
	if result := db.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser); result.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already exists"})
	}

	// Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	// Parse BirthDate
	var birthDate time.Time
	if req.BirthDate != "" {
		parsed, err := time.Parse("2006-01-02", req.BirthDate)
		if err == nil {
			birthDate = parsed
		}
	}

	// Create User
	newUser := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BirthDate:    birthDate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if result := db.DB.Create(&newUser); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	// Generate Token
	cfg := config.LoadConfig()
	token, err := utils.GenerateJWT(newUser.ID, cfg.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		Token: token,
		User:  newUser,
	})
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var user models.User
	if result := db.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate Token
	cfg := config.LoadConfig()
	token, err := utils.GenerateJWT(user.ID, cfg.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(AuthResponse{
		Token: token,
		User:  user,
	})
}
