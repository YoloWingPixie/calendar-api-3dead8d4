package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

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

// Update modifies an existing event in the database
func (r *EventRepository) Update(id string, event *Event) error {
	event.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE events
		SET title = $2, description = $3, start_time = $4, end_time = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		id,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete removes an event from the database
func (r *EventRepository) Delete(id string) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Ping checks database connectivity
func (r *EventRepository) Ping() error {
	return r.db.Ping()
}
