package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// DB wraps the SQL database connection
type DB struct {
	*sql.DB
}

// Event represents a calendar event
type Event struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description,omitempty" db:"description"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateEventRequest represents the request payload for creating an event
type CreateEventRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	StartTime   string  `json:"start_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime     string  `json:"end_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

// UpdateEventRequest represents the request payload for updating an event
type UpdateEventRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	StartTime   string  `json:"start_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime     string  `json:"end_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  string    `json:"database"`
}

// ListEventsResponse represents the response for listing events
type ListEventsResponse struct {
	Events []Event `json:"events"`
	Count  int     `json:"count"`
}

// EventRepository handles database operations for events
type EventRepository struct {
	db *DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *DB) *EventRepository {
	return &EventRepository{db: db}
}

// List retrieves all events from the database
func (r *EventRepository) List() ([]Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time, created_at, updated_at
		FROM events
		ORDER BY start_time ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return events, nil
}

// Get retrieves a single event by ID
func (r *EventRepository) Get(id string) (*Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	var event Event
	err := r.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}

// Create inserts a new event into the database
func (r *EventRepository) Create(event *Event) error {
	event.ID = uuid.New().String()
	event.CreatedAt = time.Now().UTC()
	event.UpdatedAt = event.CreatedAt

	query := `
		INSERT INTO events (id, title, description, start_time, end_time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		query,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.CreatedAt,
		event.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

// Ping checks database connectivity
func (r *EventRepository) Ping() error {
	return r.db.Ping()
}

// RunTestMigrations runs migrations for testing
func (db *DB) RunTestMigrations() error {
	migrationManager := NewMigrationManager(db)
	return migrationManager.RunMigrations()
}

// MigrationManager handles database migrations (copied for tests)
type MigrationManager struct {
	db *DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *DB) *MigrationManager {
	return &MigrationManager{db: db}
}

// RunMigrations executes all pending migrations (simplified for tests)
func (m *MigrationManager) RunMigrations() error {
	// Simple migration for tests - just create the tables
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		description VARCHAR(500) NOT NULL,
		applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		checksum VARCHAR(64) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS events (
		id VARCHAR(36) PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		start_time TIMESTAMP WITH TIME ZONE NOT NULL,
		end_time TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		CHECK (end_time > start_time)
	);

	CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);
	CREATE INDEX IF NOT EXISTS idx_events_end_time ON events(end_time);
	CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at);
	`

	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to run test migrations: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// EventHandler handles HTTP requests for events
type EventHandler struct {
	repo      *EventRepository
	validator *validator.Validate
}

// NewEventHandler creates a new event handler
func NewEventHandler(repo *EventRepository) *EventHandler {
	return &EventHandler{
		repo:      repo,
		validator: validator.New(),
	}
}

// HealthCheck handles GET /health
func (h *EventHandler) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	status := "healthy"
	dbStatus := "connected"

	// Check database connection
	if err := h.repo.Ping(); err != nil {
		status = "unhealthy"
		dbStatus = "disconnected"
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC(),
		Database:  dbStatus,
	}

	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	jsonResponse(w, statusCode, response)
}

// ListEvents handles GET /api/events
func (h *EventHandler) ListEvents(w http.ResponseWriter, _ *http.Request) {
	events, err := h.repo.List()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to retrieve events", err)
		return
	}

	response := ListEventsResponse{
		Events: events,
		Count:  len(events),
	}

	jsonResponse(w, http.StatusOK, response)
}

// Helper functions
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func errorResponse(w http.ResponseWriter, status int, message string, _ error) {
	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}

	jsonResponse(w, status, response)
}
