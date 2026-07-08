package repository

import (
	"database/sql"
	"time"

	"chat-platform/models"
)

type SessionRepository struct {
	DB *sql.DB
}

// NewSessionRepository creates a new session repository.
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		DB: db,
	}
}

// CreateSession stores a new session.
func (r *SessionRepository) CreateSession(session *models.Session) error {

	query := `
	INSERT INTO sessions (
		user_id,
		token,
		expires_at
	)
	VALUES (?, ?, ?)
	`

	result, err := r.DB.Exec(
		query,
		session.UserID,
		session.Token,
		session.ExpiresAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	session.ID = int(id)

	err = r.DB.QueryRow(
		`
		SELECT created_at
		FROM sessions
		WHERE id = ?
		`,
		session.ID,
	).Scan(&session.CreatedAt)

	return err
}

// GetSessionByToken returns a session using its token.
func (r *SessionRepository) GetSessionByToken(token string) (*models.Session, error) {

	session := &models.Session{}

	query := `
	SELECT
		id,
		user_id,
		token,
		expires_at,
		created_at
	FROM sessions
	WHERE token = ?
	`

	err := r.DB.QueryRow(query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return session, nil
}

// DeleteSession removes a session by token.
func (r *SessionRepository) DeleteSession(token string) error {

	_, err := r.DB.Exec(
		`DELETE FROM sessions WHERE token = ?`,
		token,
	)

	return err
}

// DeleteExpiredSessions removes expired sessions.
func (r *SessionRepository) DeleteExpiredSessions() error {

	_, err := r.DB.Exec(
		`
		DELETE FROM sessions
		WHERE expires_at <= ?
		`,
		time.Now(),
	)

	return err
}