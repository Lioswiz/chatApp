package models

import "time"

const (
	MessageTypeText     = "text"
	MessageTypeImage    = "image"
	MessageTypeVideo    = "video"
	MessageTypeDocument = "document"
)

type Message struct {
	ID int `json:"id"`

	SenderID int `json:"sender_id"`

	// nil = public chat
	ReceiverID *int `json:"receiver_id,omitempty"`

	Message string `json:"message"`

	MessageType string `json:"message_type"`

	UploadID *int `json:"upload_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}