package websocket

import (
	"encoding/json"
	"time"
)

const (
	MessageTypeChat        = "chat"
	MessageTypeTyping      = "typing"
	MessageTypePresence    = "presence"
	MessageTypeReadReceipt = "read_receipt"
)

//
// Every websocket message is wrapped in this structure.
//
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

//
// Chat message payload.
//
type ChatPayload struct {
	ID int `json:"id"`

	SenderID int `json:"sender_id"`

	// nil means public chat
	ReceiverID *int `json:"receiver_id,omitempty"`

	Message string `json:"message"`

	// text, image, video, document
	MessageType string `json:"message_type"`

	// Optional uploaded file
	UploadID *int `json:"upload_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

//
// Typing indicator.
//
type TypingPayload struct {
	SenderID int `json:"sender_id"`

	ReceiverID *int `json:"receiver_id,omitempty"`

	IsTyping bool `json:"is_typing"`
}

//
// User online/offline.
//
type PresencePayload struct {
	UserID int `json:"user_id"`

	Username string `json:"username"`

	Online bool `json:"online"`
}

//
// Read receipt.
//
type ReadReceiptPayload struct {
	MessageID int `json:"message_id"`

	UserID int `json:"user_id"`
}

//
// Optional helper for decoding incoming websocket messages.
//
type IncomingMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}