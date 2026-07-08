package websocket

import (
	"log"
	"net/http"

	"chat-platform/middleware"
	"chat-platform/service"

	"github.com/gorilla/websocket"
)

// Upgrader upgrades an HTTP connection to a WebSocket connection.
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins during development.
		// Restrict this in production.
		return true
	},
}

// WebSocketHandler handles websocket connections.
type WebSocketHandler struct {
	Hub         *Hub
	ChatService *service.ChatService
}

// NewWebSocketHandler creates a websocket handler.
func NewWebSocketHandler(hub *Hub, chatService *service.ChatService) *WebSocketHandler {
	return &WebSocketHandler{
		Hub:         hub,
		ChatService: chatService,
	}
}

// ServeWS upgrades the connection and registers a websocket client.
func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {

	// Get authenticated user from middleware.
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade:", err)
		return
	}

	client := &Client{
		Hub:         h.Hub,
		Conn:        conn,
		Send:        make(chan []byte, 256),
		UserID:      userID,

		// Temporary username.
		// Later we'll fetch it from the authenticated user.
		Username: "User",

		ChatService: h.ChatService,
	}

	h.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}