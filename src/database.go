package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// DB wraps the SQL database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection using the provided configuration
func NewDB(config *Config) (*DB, error) {
	log.Printf("🔌 Connecting to database: %s@%s:%s/%s", config.DBUser, config.DBHost, config.DBPort, config.DBName)

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// Note: Migration functionality has been moved to migrations.go
// Use MigrationManager for proper migration tracking and logging

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
