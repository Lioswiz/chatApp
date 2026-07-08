package websocket

import (
	"encoding/json"
	"log"
	"time"

	"chat-platform/models"
	"chat-platform/service"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 10
)

// Client represents one websocket connection.
type Client struct {
	Hub *Hub

	Conn *websocket.Conn

	Send chan []byte

	UserID   int
	Username string

	ChatService *service.ChatService
}

// ReadPump listens for incoming websocket messages.
func (c *Client) ReadPump() {

	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {

		_, data, err := c.Conn.ReadMessage()
		log.Println("Received:", string(data))
		if err != nil {

			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println("websocket:", err)
			}

			break
		}

		var incoming IncomingMessage

		if err := json.Unmarshal(data, &incoming); err != nil {
			log.Println("invalid websocket message:", err)
			continue
		}

		if incoming.Type != MessageTypeChat {
			// We'll handle typing, presence and read receipts later.
			continue
		}

		var payload ChatPayload

		if err := json.Unmarshal(incoming.Data, &payload); err != nil {
			log.Println("invalid chat payload:", err)
			continue
		}

		payload.SenderID = c.UserID
		payload.CreatedAt = time.Now()

		message := &models.Message{
			SenderID:    payload.SenderID,
			ReceiverID:  payload.ReceiverID,
			Message:     payload.Message,
			MessageType: payload.MessageType,
			UploadID:    payload.UploadID,
		}

		savedMessage, err := c.ChatService.SendMessage(message)
		if err != nil {
			log.Println("save message:", err)
			continue
		}
		if err != nil {
			log.Println("SendMessage error:", err)
			continue
		}

		log.Printf("Saved message: %+v\n", savedMessage)

		payload.ID = savedMessage.ID
		payload.CreatedAt = savedMessage.CreatedAt

		if payload.ReceiverID == nil {
			c.Hub.Broadcast <- payload
		} else {
			c.Hub.Private <- payload
		}
	}
}

// WritePump sends outgoing websocket messages.
func (c *Client) WritePump() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {

		select {

		case message, ok := <-c.Send:

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
