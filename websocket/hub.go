package websocket

import (
	"encoding/json"
	"log"
)

type Hub struct {
	// Connected clients
	Clients map[int]*Client

	// Client management
	Register   chan *Client
	Unregister chan *Client

	// Public messages
	Broadcast chan ChatPayload

	// Private messages
	Private chan ChatPayload

	// Presence updates
	Presence chan PresencePayload
}

// NewHub creates a new websocket hub.
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan ChatPayload),
		Private:    make(chan ChatPayload),
		Presence:   make(chan PresencePayload),
	}
}

// Run starts the hub.
func (h *Hub) Run() {
	for {
		select {

		//------------------------------------
		// Register Client
		//------------------------------------
		case client := <-h.Register:

			h.Clients[client.UserID] = client

			log.Printf("%s connected", client.Username)

			h.broadcast(
				"presence",
				PresencePayload{
					UserID:   client.UserID,
					Username: client.Username,
					Online:   true,
				},
			)

		//------------------------------------
		// Unregister Client
		//------------------------------------
		case client := <-h.Unregister:

			if _, ok := h.Clients[client.UserID]; ok {

				delete(h.Clients, client.UserID)

				close(client.Send)

				log.Printf("%s disconnected", client.Username)

				h.broadcast(
					"presence",
					PresencePayload{
						UserID:   client.UserID,
						Username: client.Username,
						Online:   false,
					},
				)
			}

		//------------------------------------
		// Public Chat
		//------------------------------------
		case message := <-h.Broadcast:

			h.broadcast("chat", message)

		//------------------------------------
		// Private Chat
		//------------------------------------
		case message := <-h.Private:

			if message.ReceiverID != nil {

				if receiver, ok := h.Clients[*message.ReceiverID]; ok {
					h.send(receiver, "chat", message)
				}
			}

			// Echo back to sender
			if sender, ok := h.Clients[message.SenderID]; ok {
				h.send(sender, "chat", message)
			}

		//------------------------------------
		// Presence
		//------------------------------------
		case presence := <-h.Presence:

			h.broadcast("presence", presence)
		}
	}
}

// send sends one websocket message to a client.
func (h *Hub) send(client *Client, messageType string, payload interface{}) {

	data, err := json.Marshal(WSMessage{
		Type: messageType,
		Data: payload,
	})

	if err != nil {
		log.Println(err)
		return
	}

	select {

	case client.Send <- data:

	default:

		close(client.Send)

		delete(h.Clients, client.UserID)
	}
}

// broadcast sends a websocket message to every connected client.
func (h *Hub) broadcast(messageType string, payload interface{}) {

	for _, client := range h.Clients {

		h.send(client, messageType, payload)
	}
}