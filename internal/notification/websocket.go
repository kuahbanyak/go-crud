package notification

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type NotificationType string

const (
	BookingStatusUpdate NotificationType = "booking_status_update"
	BookingReminder     NotificationType = "booking_reminder"
	ServiceComplete     NotificationType = "service_complete"
	MessageReceived     NotificationType = "message_received"
)

type Notification struct {
	Type    NotificationType `json:"type"`
	UserID  uint             `json:"user_id"`
	Title   string           `json:"title"`
	Message string           `json:"message"`
	Data    interface{}      `json:"data,omitempty"`
}

type Client struct {
	conn   *websocket.Conn
	userID uint
	send   chan Notification
}

type Hub struct {
	clients    map[uint][]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan Notification
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Notification),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.userID] = append(h.clients[client.userID], client)
			h.mutex.Unlock()
			log.Printf("Client connected: UserID %d", client.userID)

		case client := <-h.unregister:
			h.mutex.Lock()
			if clients, exists := h.clients[client.userID]; exists {
				for i, c := range clients {
					if c == client {
						close(c.send)
						h.clients[client.userID] = append(clients[:i], clients[i+1:]...)
						if len(h.clients[client.userID]) == 0 {
							delete(h.clients, client.userID)
						}
						break
					}
				}
			}
			h.mutex.Unlock()
			log.Printf("Client disconnected: UserID %d", client.userID)

		case notification := <-h.broadcast:
			h.mutex.RLock()
			if clients, exists := h.clients[notification.UserID]; exists {
				for _, client := range clients {
					select {
					case client.send <- notification:
					default:
						close(client.send)
					}
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) SendNotification(notification Notification) {
	h.broadcast <- notification
}

func (h *Hub) HandleWebSocket(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	userID := uint(claims["sub"].(float64))

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		conn:   conn,
		userID: userID,
		send:   make(chan Notification, 256),
	}

	h.register <- client

	go h.writePump(client)
	go h.readPump(client)
}

func (h *Hub) writePump(client *Client) {
	defer client.conn.Close()

	for {
		select {
		case notification, ok := <-client.send:
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(notification)
			if err != nil {
				log.Printf("Error marshaling notification: %v", err)
				continue
			}

			if err := client.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}
		}
	}
}

func (h *Hub) readPump(client *Client) {
	defer func() {
		h.unregister <- client
		client.conn.Close()
	}()

	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}
