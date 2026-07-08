package repository

import (
	"database/sql"

	"chat-platform/models"
)

type MessageRepository struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{
		DB: db,
	}
}

// SaveMessage stores a new message and returns it.
func (r *MessageRepository) SaveMessage(message *models.Message) (*models.Message, error) {

	query := `
	INSERT INTO messages
	(sender_id, receiver_id, message, message_type, upload_id)
	VALUES (?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(
		query,
		message.SenderID,
		message.ReceiverID,
		message.Message,
		message.MessageType,
		message.UploadID,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	message.ID = int(id)

	err = r.DB.QueryRow(
		`SELECT created_at FROM messages WHERE id = ?`,
		message.ID,
	).Scan(&message.CreatedAt)

	if err != nil {
		return nil, err
	}

	return message, nil
}

// GetPublicMessages returns every public message.
func (r *MessageRepository) GetPublicMessages() ([]models.Message, error) {

	query := `
	SELECT
		id,
		sender_id,
		receiver_id,
		message,
		message_type,
		upload_id,
		created_at
	FROM messages
	WHERE receiver_id IS NULL
	ORDER BY created_at ASC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message

	for rows.Next() {

		var message models.Message

		err := rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.ReceiverID,
			&message.Message,
			&message.MessageType,
			&message.UploadID,
			&message.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, rows.Err()
}

// GetPrivateMessages returns the conversation between two users.
func (r *MessageRepository) GetPrivateMessages(user1, user2 int) ([]models.Message, error) {

	query := `
	SELECT
		id,
		sender_id,
		receiver_id,
		message,
		message_type,
		upload_id,
		created_at
	FROM messages
	WHERE
	(sender_id = ? AND receiver_id = ?)
	OR
	(sender_id = ? AND receiver_id = ?)
	ORDER BY created_at ASC`

	rows, err := r.DB.Query(
		query,
		user1,
		user2,
		user2,
		user1,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message

	for rows.Next() {

		var message models.Message

		err := rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.ReceiverID,
			&message.Message,
			&message.MessageType,
			&message.UploadID,
			&message.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, rows.Err()
}

// GetMessageByID returns a single message.
func (r *MessageRepository) GetMessageByID(id int) (*models.Message, error) {

	query := `
	SELECT
		id,
		sender_id,
		receiver_id,
		message,
		message_type,
		upload_id,
		created_at
	FROM messages
	WHERE id = ?`

	var message models.Message

	err := r.DB.QueryRow(query, id).Scan(
		&message.ID,
		&message.SenderID,
		&message.ReceiverID,
		&message.Message,
		&message.MessageType,
		&message.UploadID,
		&message.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &message, nil
}

// DeleteMessage removes a message.
func (r *MessageRepository) DeleteMessage(id int) error {

	_, err := r.DB.Exec(
		`DELETE FROM messages WHERE id = ?`,
		id,
	)

	return err
}