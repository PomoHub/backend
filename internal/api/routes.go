package api

import (
	"pomodoro-habit-backend/internal/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Health Check
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "PomoHub API is running",
		})
	})

	// Auth
	auth := v1.Group("/auth")
	auth.Post("/register", Register)
	auth.Post("/login", Login)

	// Protected Routes
	cfg := config.LoadConfig()
	v1.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.JWTSecret)},
	}))

	// Spaces
	spaces := v1.Group("/spaces")
	spaces.Post("/", CreateSpace)
	spaces.Get("/", GetMySpaces)
	spaces.Post("/:spaceId/members", AddMember)
	spaces.Delete("/:spaceId/members/:userId", RemoveMember)
	spaces.Delete("/:spaceId", DeleteSpace)

	// Chat
	spaces.Post("/:spaceId/messages", SendMessage)
	spaces.Get("/:spaceId/messages", GetMessages)

	// Productivity
	// Todos
	todos := v1.Group("/todos")
	todos.Get("/", GetTodos)
	todos.Post("/", CreateTodo)
	todos.Put("/:id/toggle", ToggleTodo)
	todos.Delete("/:id", DeleteTodo)

	// Habits
	habits := v1.Group("/habits")
	habits.Get("/", GetHabits)
	habits.Post("/", CreateHabit)
	habits.Put("/:id/toggle", ToggleHabit)
	habits.Delete("/:id", DeleteHabit)

	// Pomodoro
	pomodoro := v1.Group("/pomodoro")
	pomodoro.Post("/sessions", SavePomodoroSession)
	pomodoro.Get("/stats", GetPomodoroStats)

	// Profile & Posts
	users := v1.Group("/users")
	users.Put("/me", UpdateProfile)
	users.Get("/:username", GetUserProfile)
	users.Get("/:username/posts", GetUserPosts)

	posts := v1.Group("/posts")
	posts.Post("/", CreatePost)
	posts.Delete("/:id", DeletePost)
}
