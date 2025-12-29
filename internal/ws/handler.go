package ws

import (
	"log"
	"pomodoro-habit-backend/internal/config"
	"pomodoro-habit-backend/internal/db"
	"pomodoro-habit-backend/internal/models"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Initialize Global Hub
var GlobalHub = NewHub()

func SetupWebSockets(app *fiber.App) {
	// Start Hub
	go GlobalHub.Run()

	// WebSocket Middleware for Auth
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return c.Status(fiber.StatusUpgradeRequired).SendString("Upgrade Required")
	})

	app.Get("/ws/:spaceId", websocket.New(func(c *websocket.Conn) {
		// Extract Token from Query Param (WS doesn't support headers well in standard JS API)
		tokenStr := c.Query("token")
		if tokenStr == "" {
			log.Println("WS: No token provided")
			c.Close()
			return
		}

		// Validate Token
		cfg := config.LoadConfig()
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			log.Println("WS: Invalid token")
			c.Close()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userIDStr := claims["user_id"].(string)
		userID, _ := uuid.Parse(userIDStr)

		spaceIDStr := c.Params("spaceId")
		spaceID, err := uuid.Parse(spaceIDStr)
		if err != nil {
			log.Println("WS: Invalid Space ID")
			c.Close()
			return
		}

		// Check Space Membership
		var membership models.SpaceMember
		if err := db.DB.Where("space_id = ? AND user_id = ?", spaceID, userID).First(&membership).Error; err != nil {
			log.Println("WS: Not a member of space")
			c.Close()
			return
		}

		client := &Client{
			ID:      userID,
			Conn:    c,
			SpaceID: spaceID,
			Hub:     GlobalHub,
		}

		client.Hub.register <- client
		defer func() {
			client.Hub.unregister <- client
		}()

		// Notify others that user joined
		client.Hub.BroadcastToSpace(spaceID, TypeUserJoined, map[string]string{
			"user_id": userID.String(),
		})

		for {
			var msg WSMessage
			if err := c.ReadJSON(&msg); err != nil {
				log.Println("WS: Read error:", err)
				break
			}

			// Handle incoming messages if needed (e.g., live typing status, or pomodoro updates sent via WS)
			// For simplicity, we can just broadcast it back to the room if it's a valid type
			msg.SpaceID = spaceID // Ensure space ID is correct
			
			// If it's a pomodoro status update
			if msg.Type == TypePomodoroStatus {
				// Broadcast to others
				client.Hub.broadcast <- msg
			}
		}
	}))
}
