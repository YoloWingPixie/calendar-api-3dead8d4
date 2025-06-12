package main

import (
	"time"
)

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

// VersionResponse represents the version check response
type VersionResponse struct {
	Version   string    `json:"version"`
	BuildInfo string    `json:"build_info"`
	Timestamp time.Time `json:"timestamp"`
}

// ListEventsResponse represents the response for listing events
type ListEventsResponse struct {
	Events []Event `json:"events"`
	Count  int     `json:"count"`
}
