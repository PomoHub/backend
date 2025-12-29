package main

import (
	"log"
	"pomodoro-habit-backend/internal/api"
	"pomodoro-habit-backend/internal/config"
	"pomodoro-habit-backend/internal/db"
	"pomodoro-habit-backend/internal/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Connect to Database
	db.ConnectDB(cfg)

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: "PomoHub API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	api.SetupRoutes(app)
	ws.SetupWebSockets(app)

	// Start Server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
