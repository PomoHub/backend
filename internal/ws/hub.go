package ws

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

// Message types
const (
	TypeChatMessage    = "chat_message"
	TypePomodoroStatus = "pomodoro_status"
	TypeUserJoined     = "user_joined"
	TypeUserLeft       = "user_left"
)

// WebSocket Message Structure
type WSMessage struct {
	Type    string      `json:"type"`
	SpaceID uuid.UUID   `json:"space_id"`
	Payload interface{} `json:"payload"`
}

// Client represents a connected user
type Client struct {
	ID      uuid.UUID
	Conn    *websocket.Conn
	SpaceID uuid.UUID
	Hub     *Hub
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients map[SpaceID]map[ClientID]*Client
	clients    map[uuid.UUID]map[uuid.UUID]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan WSMessage
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]map[uuid.UUID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan WSMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			if _, ok := h.clients[client.SpaceID]; !ok {
				h.clients[client.SpaceID] = make(map[uuid.UUID]*Client)
			}
			h.clients[client.SpaceID][client.ID] = client
			h.mutex.Unlock()
			log.Printf("Client registered: %s in Space %s", client.ID, client.SpaceID)

		case client := <-h.unregister:
			h.mutex.Lock()
			if space, ok := h.clients[client.SpaceID]; ok {
				if _, ok := space[client.ID]; ok {
					delete(space, client.ID)
					client.Conn.Close()
					if len(space) == 0 {
						delete(h.clients, client.SpaceID)
					}
				}
			}
			h.mutex.Unlock()
			log.Printf("Client unregistered: %s", client.ID)

		case message := <-h.broadcast:
			h.mutex.RLock()
			if clients, ok := h.clients[message.SpaceID]; ok {
				for _, client := range clients {
					if err := client.Conn.WriteJSON(message); err != nil {
						log.Printf("Error sending message: %v", err)
						client.Conn.Close()
						delete(clients, client.ID)
					}
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// Helper to broadcast message from API handlers
func (h *Hub) BroadcastToSpace(spaceID uuid.UUID, msgType string, payload interface{}) {
	h.broadcast <- WSMessage{
		Type:    msgType,
		SpaceID: spaceID,
		Payload: payload,
	}
}
