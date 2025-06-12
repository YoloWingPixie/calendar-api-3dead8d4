package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// NewTestDB creates a test database connection for testing
func NewTestDB() (*DB, error) {
	// Use test database configuration
	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("TEST_DB_PORT")
	if port == "" {
		port = "5485"
	}

	user := os.Getenv("TEST_DB_USER")
	if user == "" {
		user = "test_user"
	}

	password := os.Getenv("TEST_DB_PASSWORD")
	if password == "" {
		password = "test_pass"
	}

	dbname := os.Getenv("TEST_DB_NAME")
	if dbname == "" {
		dbname = "calendar_test_db"
	}

	sslmode := os.Getenv("TEST_DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open test database: %w", err)
	}

	// Configure connection pool for testing
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(1 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping test database: %w", err)
	}

	return &DB{db}, nil
}

// CleanupTestDB clears test data but keeps tables
func CleanupTestDB(db *DB) error {
	// Delete all events
	if _, err := db.Exec("DELETE FROM events"); err != nil {
		return fmt.Errorf("failed to clean events table: %w", err)
	}

	// Delete all users
	if _, err := db.Exec("DELETE FROM users"); err != nil {
		return fmt.Errorf("failed to clean users table: %w", err)
	}

	return nil
}
