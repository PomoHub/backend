package api

import (
	"pomodoro-habit-backend/internal/db"
	"pomodoro-habit-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Get Todos
func GetTodos(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var todos []models.Todo
	if err := db.DB.Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch todos"})
	}

	return c.JSON(todos)
}

// Create Todo
func CreateTodo(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type Request struct {
		Title string `json:"title"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required"})
	}

	todo := models.Todo{
		UserID: userID,
		Title:  req.Title,
	}

	if err := db.DB.Create(&todo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create todo"})
	}

	return c.Status(fiber.StatusCreated).JSON(todo)
}

// Toggle Todo
func ToggleTodo(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	todoIDStr := c.Params("id")
	todoID, _ := uuid.Parse(todoIDStr)

	var todo models.Todo
	if err := db.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	}

	todo.Completed = !todo.Completed
	if err := db.DB.Save(&todo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update todo"})
	}

	return c.JSON(todo)
}

// Delete Todo
func DeleteTodo(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	todoIDStr := c.Params("id")
	todoID, _ := uuid.Parse(todoIDStr)

	if err := db.DB.Where("id = ? AND user_id = ?", todoID, userID).Delete(&models.Todo{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete todo"})
	}

	return c.JSON(fiber.Map{"message": "Todo deleted successfully"})
}
