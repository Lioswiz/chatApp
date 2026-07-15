package service

import (
	"errors"

	"chat-platform/models"
	"chat-platform/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

// Constructor
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(user *models.User, password string) error {

	// Validate the user input.
	if err := ValidateUser(user, password); err != nil {
		return err
	}

	// Check if the email already exists.
	_, err := s.UserRepo.GetUserByEmail(user.Email)
	if err == nil {
		return errors.New("email already exists")
	}

	// Hash the password.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)

	// Save the user.
	return s.UserRepo.CreateUser(user)
}

// Login authenticates a user.
func (s *AuthService) Login(email, password string) (*models.User, error) {

	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// UpdateUser updates the user record.
func (s *AuthService) UpdateUser(user *models.User) error {
	return s.UserRepo.UpdateUser(user)
}