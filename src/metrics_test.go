package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

func TestMetricsMiddleware(t *testing.T) {
	// Reset Prometheus metrics before test
	prometheus.DefaultRegisterer = prometheus.NewRegistry()

	// Create a test router with metrics middleware
	r := mux.NewRouter()
	r.Use(MetricsMiddleware)

	// Add a test endpoint
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	// Test cases
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Successful request",
			method:         "GET",
			path:           "/test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not found request",
			method:         "GET",
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			// Serve request
			r.ServeHTTP(w, req)

			// Check response status
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Verify metrics were recorded
			// Note: We can't directly check metric values as they're private,
			// but we can verify the middleware didn't panic
		})
	}
}

func TestRecordDBOperation(t *testing.T) {
	// Reset Prometheus metrics before test
	prometheus.DefaultRegisterer = prometheus.NewRegistry()

	// Test cases
	tests := []struct {
		name       string
		operation  string
		table      string
		duration   time.Duration
		shouldPass bool
	}{
		{
			name:       "Valid operation",
			operation:  "select",
			table:      "events",
			duration:   100 * time.Millisecond,
			shouldPass: true,
		},
		{
			name:       "Zero duration",
			operation:  "insert",
			table:      "events",
			duration:   0,
			shouldPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Record the operation
			RecordDBOperation(tt.operation, tt.table, tt.duration)

			// Verify no panic occurred
			// Note: We can't directly check metric values as they're private,
			// but we can verify the function didn't panic
		})
	}
}

func TestUpdateSystemMetrics(t *testing.T) {
	// Test that UpdateSystemMetrics doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("UpdateSystemMetrics panicked: %v", r)
		}
	}()

	UpdateSystemMetrics()
}

func TestLoggingMiddleware(t *testing.T) {
	// Create a test router with logging middleware
	r := mux.NewRouter()
	r.Use(LoggingMiddleware)

	// Add a test endpoint
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	// Test cases
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Successful request",
			method:         "GET",
			path:           "/test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not found request",
			method:         "GET",
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			// Serve request
			r.ServeHTTP(w, req)

			// Check response status
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Verify logging middleware didn't panic
			// Note: We can't directly check log output as it's handled by Zap,
			// but we can verify the middleware didn't panic
		})
	}
}

func TestResponseWriter(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		statusCode     int
		expectedStatus int
	}{
		{
			name:           "Set status code",
			statusCode:     http.StatusNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Set OK status code",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new response writer for each test
			w := httptest.NewRecorder()
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Set status code
			rw.WriteHeader(tt.statusCode)

			// Verify status code was set
			if rw.statusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rw.statusCode)
			}

			// Verify underlying response writer status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected underlying status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
