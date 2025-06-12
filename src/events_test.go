package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

// MockEventRepository for testing
type MockEventRepository struct {
	events  map[string]*Event
	counter int
}

func NewMockEventRepository() *MockEventRepository {
	return &MockEventRepository{
		events:  make(map[string]*Event),
		counter: 0,
	}
}

func (m *MockEventRepository) Create(event *Event) error {
	m.counter++
	event.ID = fmt.Sprintf("test-id-%d", m.counter)
	event.CreatedAt = time.Now()
	event.UpdatedAt = event.CreatedAt
	m.events[event.ID] = event
	return nil
}

func (m *MockEventRepository) Get(id string) (*Event, error) {
	event, exists := m.events[id]
	if !exists {
		return nil, nil
	}
	return event, nil
}

func (m *MockEventRepository) Update(id string, event *Event) error {
	event.ID = id
	event.UpdatedAt = time.Now()
	m.events[id] = event
	return nil
}

func (m *MockEventRepository) Delete(id string) error {
	if _, exists := m.events[id]; !exists {
		return sql.ErrNoRows
	}
	delete(m.events, id)
	return nil
}

func (m *MockEventRepository) List() ([]Event, error) {
	events := make([]Event, 0, len(m.events))
	for _, event := range m.events {
		events = append(events, *event)
	}
	// Sort by start time to match the real repository behavior
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Before(events[j].StartTime)
	})
	return events, nil
}

func (m *MockEventRepository) Ping() error {
	return nil
}

// MockAuthMiddleware for testing
type MockAuthMiddleware struct {
	config *Config
}

func NewMockAuthMiddleware(config *Config) *MockAuthMiddleware {
	return &MockAuthMiddleware{
		config: config,
	}
}

func (m *MockAuthMiddleware) RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(m.config.APIKeyHeader)
		if apiKey == "" {
			errorResponse(w, http.StatusUnauthorized, "API key required", nil)
			return
		}

		// For testing, accept the bootstrap key
		if apiKey != m.config.BootstrapAdminKey {
			errorResponse(w, http.StatusUnauthorized, "Invalid API key", nil)
			return
		}

		// Create a test user context
		testUser := &User{
			ID:       "test-user",
			Username: "testuser",
			APIKey:   apiKey,
		}
		r = r.WithContext(WithUser(r.Context(), testUser))
		next.ServeHTTP(w, r)
	})
}

// setupEventTest creates a test setup with mocks
func setupEventTest(_ *testing.T) (*EventHandler, *MockAuthMiddleware) {
	// Test configuration
	config := &Config{
		BootstrapAdminKey: "test-admin-key-123",
		APIKeyHeader:      "X-API-Key",
		Environment:       "test",
	}

	// Setup mock repositories
	eventRepo := NewMockEventRepository()

	// Setup auth middleware and handlers
	authMiddleware := NewMockAuthMiddleware(config)
	eventHandler := NewEventHandler(eventRepo)

	return eventHandler, authMiddleware
}

func TestEventsCreateEvent(t *testing.T) {
	handler, auth := setupEventTest(t)

	tests := []struct {
		name           string
		payload        interface{}
		apiKey         string
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Valid event creation",
			payload: CreateEventRequest{
				Title:       "Team Meeting",
				Description: stringPtr("Weekly team sync"),
				StartTime:   time.Now().Add(time.Hour).Format(time.RFC3339),
				EndTime:     time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			},
			apiKey:         "test-admin-key-123",
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "No authentication",
			payload: CreateEventRequest{
				Title:     "Test Event",
				StartTime: time.Now().Add(time.Hour).Format(time.RFC3339),
				EndTime:   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			},
			apiKey:         "",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name: "Invalid API key",
			payload: CreateEventRequest{
				Title:     "Test Event",
				StartTime: time.Now().Add(time.Hour).Format(time.RFC3339),
				EndTime:   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			},
			apiKey:         "invalid-key",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name: "Missing title",
			payload: CreateEventRequest{
				StartTime: time.Now().Add(time.Hour).Format(time.RFC3339),
				EndTime:   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			},
			apiKey:         "test-admin-key-123",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "End time before start time",
			payload: CreateEventRequest{
				Title:     "Test Event",
				StartTime: time.Now().Add(2 * time.Hour).Format(time.RFC3339),
				EndTime:   time.Now().Add(time.Hour).Format(time.RFC3339),
			},
			apiKey:         "test-admin-key-123",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("Failed to marshal payload: %v", err)
			}

			req := httptest.NewRequest("POST", "/api/events", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			w := httptest.NewRecorder()

			// Apply auth middleware for protected endpoints
			if tt.expectedStatus == http.StatusUnauthorized || tt.apiKey != "" {
				handler := auth.RequireAPIKey(http.HandlerFunc(handler.CreateEvent))
				handler.ServeHTTP(w, req)
			} else {
				handler.CreateEvent(w, req)
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if !tt.expectError && w.Code == http.StatusCreated {
				var event Event
				if err := json.Unmarshal(w.Body.Bytes(), &event); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if event.ID == "" {
					t.Error("Expected event ID to be generated")
				}
				if event.CreatedAt.IsZero() {
					t.Error("Expected CreatedAt to be set")
				}
				if event.Title != tt.payload.(CreateEventRequest).Title {
					t.Errorf("Expected title %s, got %s", tt.payload.(CreateEventRequest).Title, event.Title)
				}
			}
		})
	}
}

func TestEventsGetEvent(t *testing.T) {
	handler, auth := setupEventTest(t)

	// Create a test event first
	createPayload := CreateEventRequest{
		Title:       "Test Event",
		Description: stringPtr("Test Description"),
		StartTime:   time.Now().Add(time.Hour).Format(time.RFC3339),
		EndTime:     time.Now().Add(2 * time.Hour).Format(time.RFC3339),
	}

	body, _ := json.Marshal(createPayload)
	req := httptest.NewRequest("POST", "/api/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "test-admin-key-123")
	w := httptest.NewRecorder()

	createHandler := auth.RequireAPIKey(http.HandlerFunc(handler.CreateEvent))
	createHandler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create test event: %d %s", w.Code, w.Body.String())
	}

	var createdEvent Event
	json.Unmarshal(w.Body.Bytes(), &createdEvent)

	tests := []struct {
		name           string
		eventID        string
		apiKey         string
		expectedStatus int
	}{
		{
			name:           "Get existing event",
			eventID:        createdEvent.ID,
			apiKey:         "test-admin-key-123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get non-existent event",
			eventID:        "non-existent-id",
			apiKey:         "test-admin-key-123",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "No authentication",
			eventID:        createdEvent.ID,
			apiKey:         "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/events/"+tt.eventID, nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			// Set up router with path variable
			w := httptest.NewRecorder()
			router := mux.NewRouter()

			if tt.expectedStatus == http.StatusUnauthorized || tt.apiKey != "" {
				finalHandler := auth.RequireAPIKey(http.HandlerFunc(handler.GetEvent))
				router.Handle("/api/events/{id}", finalHandler).Methods("GET")
			} else {
				router.HandleFunc("/api/events/{id}", handler.GetEvent).Methods("GET")
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.expectedStatus == http.StatusOK {
				var event Event
				if err := json.Unmarshal(w.Body.Bytes(), &event); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if event.ID != createdEvent.ID {
					t.Errorf("Expected event ID %s, got %s", createdEvent.ID, event.ID)
				}
				if event.Title != "Test Event" {
					t.Errorf("Expected title 'Test Event', got '%s'", event.Title)
				}
			}
		})
	}
}

func TestEventsUpdateEvent(t *testing.T) {
	handler, auth := setupEventTest(t)

	// Create a test event first
	createPayload := CreateEventRequest{
		Title:       "Original Event",
		Description: stringPtr("Original Description"),
		StartTime:   time.Now().Add(time.Hour).Format(time.RFC3339),
		EndTime:     time.Now().Add(2 * time.Hour).Format(time.RFC3339),
	}

	body, _ := json.Marshal(createPayload)
	req := httptest.NewRequest("POST", "/api/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "test-admin-key-123")
	w := httptest.NewRecorder()

	createHandler := auth.RequireAPIKey(http.HandlerFunc(handler.CreateEvent))
	createHandler.ServeHTTP(w, req)

	var createdEvent Event
	json.Unmarshal(w.Body.Bytes(), &createdEvent)

	tests := []struct {
		name           string
		eventID        string
		payload        UpdateEventRequest
		apiKey         string
		expectedStatus int
	}{
		{
			name:    "Update existing event",
			eventID: createdEvent.ID,
			payload: UpdateEventRequest{
				Title:       "Updated Event",
				Description: stringPtr("Updated Description"),
				StartTime:   time.Now().Add(time.Hour).Format(time.RFC3339),
				EndTime:     time.Now().Add(3 * time.Hour).Format(time.RFC3339),
			},
			apiKey:         "test-admin-key-123",
			expectedStatus: http.StatusOK,
		},
		{
			name:    "Update non-existent event",
			eventID: "non-existent-id",
			payload: UpdateEventRequest{
				Title:     "Updated Event",
				StartTime: time.Now().Add(time.Hour).Format(time.RFC3339),
				EndTime:   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			},
			apiKey:         "test-admin-key-123",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:    "No authentication",
			eventID: createdEvent.ID,
			payload: UpdateEventRequest{
				Title:     "Updated Event",
				StartTime: time.Now().Add(time.Hour).Format(time.RFC3339),
				EndTime:   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			},
			apiKey:         "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("PUT", "/api/events/"+tt.eventID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			w := httptest.NewRecorder()
			router := mux.NewRouter()

			if tt.expectedStatus == http.StatusUnauthorized || tt.apiKey != "" {
				finalHandler := auth.RequireAPIKey(http.HandlerFunc(handler.UpdateEvent))
				router.Handle("/api/events/{id}", finalHandler).Methods("PUT")
			} else {
				router.HandleFunc("/api/events/{id}", handler.UpdateEvent).Methods("PUT")
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.expectedStatus == http.StatusOK {
				var event Event
				if err := json.Unmarshal(w.Body.Bytes(), &event); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if event.Title != tt.payload.Title {
					t.Errorf("Expected title '%s', got '%s'", tt.payload.Title, event.Title)
				}
			}
		})
	}
}

func TestEventsDeleteEvent(t *testing.T) {
	handler, auth := setupEventTest(t)

	// Create a test event first
	createPayload := CreateEventRequest{
		Title:     "Event to Delete",
		StartTime: time.Now().Add(time.Hour).Format(time.RFC3339),
		EndTime:   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
	}

	body, _ := json.Marshal(createPayload)
	req := httptest.NewRequest("POST", "/api/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "test-admin-key-123")
	w := httptest.NewRecorder()

	createHandler := auth.RequireAPIKey(http.HandlerFunc(handler.CreateEvent))
	createHandler.ServeHTTP(w, req)

	var createdEvent Event
	json.Unmarshal(w.Body.Bytes(), &createdEvent)

	tests := []struct {
		name           string
		eventID        string
		apiKey         string
		expectedStatus int
	}{
		{
			name:           "Delete existing event",
			eventID:        createdEvent.ID,
			apiKey:         "test-admin-key-123",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Delete non-existent event",
			eventID:        "non-existent-id",
			apiKey:         "test-admin-key-123",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "No authentication",
			eventID:        createdEvent.ID,
			apiKey:         "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/events/"+tt.eventID, nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			w := httptest.NewRecorder()
			router := mux.NewRouter()

			if tt.expectedStatus == http.StatusUnauthorized || tt.apiKey != "" {
				finalHandler := auth.RequireAPIKey(http.HandlerFunc(handler.DeleteEvent))
				router.Handle("/api/events/{id}", finalHandler).Methods("DELETE")
			} else {
				router.HandleFunc("/api/events/{id}", handler.DeleteEvent).Methods("DELETE")
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			// For successful deletion, verify event is gone
			if tt.expectedStatus == http.StatusNoContent {
				// Try to get the deleted event
				getReq := httptest.NewRequest("GET", "/api/events/"+tt.eventID, nil)
				getReq.Header.Set("X-API-Key", "test-admin-key-123")
				getW := httptest.NewRecorder()

				getRouter := mux.NewRouter()
				getHandler := auth.RequireAPIKey(http.HandlerFunc(handler.GetEvent))
				getRouter.Handle("/api/events/{id}", getHandler).Methods("GET")
				getRouter.ServeHTTP(getW, getReq)

				if getW.Code != http.StatusNotFound {
					t.Errorf("Expected deleted event to return 404, got %d", getW.Code)
				}
			}
		})
	}
}

func TestEventsListEvents(t *testing.T) {
	handler, auth := setupEventTest(t)

	// Test empty list first
	t.Run("Empty list", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/events", nil)
		req.Header.Set("X-API-Key", "test-admin-key-123")
		w := httptest.NewRecorder()

		listHandler := auth.RequireAPIKey(http.HandlerFunc(handler.ListEvents))
		listHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response ListEventsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Count != 0 {
			t.Errorf("Expected count 0, got %d", response.Count)
		}
	})

	// Create some test events
	events := []CreateEventRequest{
		{
			Title:     "Event 1",
			StartTime: time.Now().Add(time.Hour).Format(time.RFC3339),
			EndTime:   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
		},
		{
			Title:     "Event 2",
			StartTime: time.Now().Add(3 * time.Hour).Format(time.RFC3339),
			EndTime:   time.Now().Add(4 * time.Hour).Format(time.RFC3339),
		},
	}

	// Create events
	for _, eventReq := range events {
		body, _ := json.Marshal(eventReq)
		req := httptest.NewRequest("POST", "/api/events", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "test-admin-key-123")
		w := httptest.NewRecorder()

		createHandler := auth.RequireAPIKey(http.HandlerFunc(handler.CreateEvent))
		createHandler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Failed to create test event: %d", w.Code)
		}
	}

	// Test list with events
	t.Run("List with events", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/events", nil)
		req.Header.Set("X-API-Key", "test-admin-key-123")
		w := httptest.NewRecorder()

		listHandler := auth.RequireAPIKey(http.HandlerFunc(handler.ListEvents))
		listHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response ListEventsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Count != 2 {
			t.Errorf("Expected count 2, got %d", response.Count)
		}

		if len(response.Events) != 2 {
			t.Errorf("Expected 2 events, got %d", len(response.Events))
		}

		// Verify events are ordered by start time
		if len(response.Events) >= 2 {
			if !response.Events[0].StartTime.Before(response.Events[1].StartTime) {
				t.Error("Events should be ordered by start time")
			}
		}
	})

	// Test unauthorized access
	t.Run("Unauthorized access", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/events", nil)
		w := httptest.NewRecorder()

		listHandler := auth.RequireAPIKey(http.HandlerFunc(handler.ListEvents))
		listHandler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
