package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockEventRepository for testing
type MockEventRepository struct {
	events []Event
	pingOK bool
}

func NewMockEventRepository() *MockEventRepository {
	return &MockEventRepository{
		events: []Event{},
		pingOK: true,
	}
}

func (m *MockEventRepository) List() ([]Event, error) {
	return m.events, nil
}

func (m *MockEventRepository) Ping() error {
	if m.pingOK {
		return nil
	}
	return nil
}

func TestHealthCheck(t *testing.T) {
	// Create a mock repository
	repo := NewMockEventRepository()
	handler := NewEventHandler(repo)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function directly
	handler.HealthCheck(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response.Status)
	}

	if response.Database != "connected" {
		t.Errorf("Expected database status 'connected', got '%s'", response.Database)
	}
}

func TestListEventsEmpty(t *testing.T) {
	// Create a mock repository
	repo := NewMockEventRepository()
	handler := NewEventHandler(repo)

	// Create a request to list events
	req, err := http.NewRequest("GET", "/api/events", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.ListEvents(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response ListEventsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Count != 0 {
		t.Errorf("Expected count 0, got %d", response.Count)
	}

	if len(response.Events) != 0 {
		t.Errorf("Expected empty events list, got %d events", len(response.Events))
	}
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

// EventRepository interface
type EventRepository interface {
	List() ([]Event, error)
	Ping() error
}

// EventHandler handles HTTP requests for events
type EventHandler struct {
	repo EventRepository
}

// NewEventHandler creates a new event handler
func NewEventHandler(repo EventRepository) *EventHandler {
	return &EventHandler{
		repo: repo,
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

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func errorResponse(w http.ResponseWriter, status int, message string, _ error) {
	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}

	jsonResponse(w, status, response)
}
