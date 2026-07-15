package repository

import (
	"database/sql"

	"chat-platform/models"
)

type UserRepository struct {
	DB *sql.DB
}

// Constructor
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// CreateUser inserts a new user into the database.
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
	INSERT INTO users
	(first_name, last_name, username, email, password_hash, avatar)
	VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.DB.Exec(
		query,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Avatar,
	)

	return err
}

// GetUserByEmail returns a user by email.
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
	SELECT
		id,
		first_name,
		last_name,
		username,
		email,
		password_hash,
		avatar,
		created_at
	FROM users
	WHERE email = ?`

	user := &models.User{}

	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Avatar,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID returns a user using their ID.
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `
	SELECT
		id,
		first_name,
		last_name,
		username,
		email,
		password_hash,
		avatar,
		created_at
	FROM users
	WHERE id = ?`

	user := &models.User{}

	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Avatar,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates user details.
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
	UPDATE users
	SET first_name = ?, last_name = ?, avatar = ?
	WHERE id = ?`

	_, err := r.DB.Exec(
		query,
		user.FirstName,
		user.LastName,
		user.Avatar,
		user.ID,
	)

	return err
}