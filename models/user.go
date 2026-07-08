
package models

import "time"

// User represents a registered user.
type User struct {
	ID           int
	FirstName    string
	LastName     string
	Username     string
	Email        string
	PasswordHash string
	Avatar       string
	CreatedAt    time.Time
}