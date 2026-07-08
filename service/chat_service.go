package service

import (
	"errors"
	"strings"

	"chat-platform/models"
	"chat-platform/repository"
)

type ChatService struct {
	MessageRepo *repository.MessageRepository
}

// NewChatService creates a new ChatService.
func NewChatService(messageRepo *repository.MessageRepository) *ChatService {
	return &ChatService{
		MessageRepo: messageRepo,
	}
}

// SendMessage validates and stores a message.
func (s *ChatService) SendMessage(message *models.Message) (*models.Message, error) {

	if message == nil {
		return nil, errors.New("message cannot be nil")
	}

	// Sender must exist.
	if message.SenderID <= 0 {
		return nil, errors.New("invalid sender")
	}

	// Default message type.
	if message.MessageType == "" {
		message.MessageType = models.MessageTypeText
	}

	// Trim text/caption.
	message.Message = strings.TrimSpace(message.Message)

	switch message.MessageType {

	case models.MessageTypeText:

		if message.Message == "" {
			return nil, errors.New("text message cannot be empty")
		}

	case models.MessageTypeImage,
		models.MessageTypeVideo,
		models.MessageTypeDocument:

		// Media messages require an uploaded file.
		if message.UploadID == nil {
			return nil, errors.New("media message requires an upload")
		}

	default:
		return nil, errors.New("invalid message type")
	}

	// Save and return the stored message.
	return s.MessageRepo.SaveMessage(message)
}

// GetPublicMessages returns all public chat messages.
func (s *ChatService) GetPublicMessages() ([]models.Message, error) {
	return s.MessageRepo.GetPublicMessages()
}

// GetPrivateMessages returns all messages exchanged between two users.
func (s *ChatService) GetPrivateMessages(user1, user2 int) ([]models.Message, error) {

	if user1 <= 0 || user2 <= 0 {
		return nil, errors.New("invalid user id")
	}

	if user1 == user2 {
		return nil, errors.New("cannot chat with yourself")
	}

	return s.MessageRepo.GetPrivateMessages(user1, user2)
}

// GetMessageByID returns a single message.
func (s *ChatService) GetMessageByID(id int) (*models.Message, error) {

	if id <= 0 {
		return nil, errors.New("invalid message id")
	}

	return s.MessageRepo.GetMessageByID(id)
}

// DeleteMessage removes a message.
func (s *ChatService) DeleteMessage(id int) error {

	if id <= 0 {
		return errors.New("invalid message id")
	}

	return s.MessageRepo.DeleteMessage(id)
}