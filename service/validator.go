package service

import (
	"errors"
	"strings"

	"chat-platform/models"
)

func ValidateUser(user *models.User, password string) error {

	if strings.TrimSpace(user.FirstName) == "" {
		return errors.New("first name is required")
	}

	if strings.TrimSpace(user.LastName) == "" {
		return errors.New("last name is required")
	}

	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username is required")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}