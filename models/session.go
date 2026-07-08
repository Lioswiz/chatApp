package models

import "time"

// Session represents a logged-in user's session.
type Session struct {
	ID int

	UserID int

	// Random token stored in the user's cookie.
	Token string

	// Session expiration time.
	ExpiresAt time.Time

	CreatedAt time.Time
}