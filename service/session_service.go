package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"chat-platform/models"
	"chat-platform/repository"
)

const (
	SessionDuration = 24 * time.Hour
)

type SessionService struct {
	SessionRepo *repository.SessionRepository
	UserRepo    *repository.UserRepository
}

// Constructor
func NewSessionService(
	sessionRepo *repository.SessionRepository,
	userRepo *repository.UserRepository,
) *SessionService {
	return &SessionService{
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
	}
}

// CreateSession creates a new login session for a user.
func (s *SessionService) CreateSession(userID int) (*models.Session, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(SessionDuration),
	}

	err = s.SessionRepo.CreateSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// ValidateSession checks whether a session token is valid.
func (s *SessionService) ValidateSession(token string) (*models.User, error) {

	if token == "" {
		return nil, errors.New("missing session token")
	}

	session, err := s.SessionRepo.GetSessionByToken(token)
	if err != nil {
		return nil, errors.New("invalid session")
	}

	if time.Now().After(session.ExpiresAt) {

		_ = s.SessionRepo.DeleteSession(token)

		return nil, errors.New("session expired")
	}

	user, err := s.UserRepo.GetUserByID(session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Logout removes a session.
func (s *SessionService) Logout(token string) error {

	if token == "" {
		return nil
	}

	return s.SessionRepo.DeleteSession(token)
}

// generateToken creates a secure random session token.
func generateToken() (string, error) {

	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}