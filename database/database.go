package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Connect opens the SQLite database, creates the tables if they
// don't already exist, and returns the database connection.
func Connect(databasePath string) (*sql.DB, error) {
	// Open the database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create all tables
	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// createTables reads schema.sql and executes it.
func createTables(db *sql.DB) error {
	schema, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		return fmt.Errorf("unable to read schema.sql: %w", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}
