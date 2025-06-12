package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	// Create a mock repository
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Clean up any existing data
	if cleanupErr := CleanupTestDB(db); cleanupErr != nil {
		t.Logf("Warning: Failed to cleanup test database: %v", cleanupErr)
	}

	repo := NewEventRepository(db)
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
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if migrationErr := db.RunTestMigrations(); migrationErr != nil {
		t.Fatalf("Failed to run migrations: %v", migrationErr)
	}

	// Clean up any existing data
	if cleanupErr := CleanupTestDB(db); cleanupErr != nil {
		t.Fatalf("Failed to cleanup test database: %v", cleanupErr)
	}

	repo := NewEventRepository(db)
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
