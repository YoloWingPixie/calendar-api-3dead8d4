package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// EventHandler handles HTTP requests for events
type EventHandler struct {
	repo      EventRepositoryInterface
	validator *validator.Validate
}

// NewEventHandler creates a new event handler
func NewEventHandler(repo EventRepositoryInterface) *EventHandler {
	return &EventHandler{
		repo:      repo,
		validator: validator.New(),
	}
}

// ListEvents handles GET /api/events
func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		RecordDBOperation("list", "events", time.Since(start))
	}()

	events, err := h.repo.List()
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to retrieve events", err)
		return
	}

	response := ListEventsResponse{
		Events: events,
		Count:  len(events),
	}

	h.jsonResponse(w, http.StatusOK, response)
}

// GetEvent handles GET /api/events/{id}
func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		RecordDBOperation("get", "events", time.Since(start))
	}()

	vars := mux.Vars(r)
	id := vars["id"]

	event, err := h.repo.Get(id)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to retrieve event", err)
		return
	}

	if event == nil {
		h.errorResponse(w, http.StatusNotFound, "Event not found", nil)
		return
	}

	h.jsonResponse(w, http.StatusOK, event)
}

// sanitizeString removes potentially dangerous characters from input strings
func sanitizeString(s string) string {
	// Replace HTML special characters with their escaped versions
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// CreateEvent handles POST /api/events
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		RecordDBOperation("create", "events", time.Since(start))
	}()

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.validationErrorResponse(w, err)
		return
	}

	// Sanitize input strings
	req.Title = sanitizeString(req.Title)
	if req.Description != nil {
		sanitizedDesc := sanitizeString(*req.Description)
		req.Description = &sanitizedDesc
	}

	// Parse times
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid start_time format", err)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid end_time format", err)
		return
	}

	// Validate time range
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		h.errorResponse(w, http.StatusBadRequest, "end_time must be after start_time", nil)
		return
	}

	// Create event
	event := &Event{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	if err := h.repo.Create(event); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to create event", err)
		return
	}

	// Record metrics
	eventsCreatedTotal.Inc()
	activeEvents.Inc()

	h.jsonResponse(w, http.StatusCreated, event)
}

// UpdateEvent handles PUT /api/events/{id}
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		RecordDBOperation("update", "events", time.Since(start))
	}()

	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.validationErrorResponse(w, err)
		return
	}

	// Sanitize input strings
	req.Title = sanitizeString(req.Title)
	if req.Description != nil {
		sanitizedDesc := sanitizeString(*req.Description)
		req.Description = &sanitizedDesc
	}

	// Check if event exists
	existing, err := h.repo.Get(id)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to retrieve event", err)
		return
	}
	if existing == nil {
		h.errorResponse(w, http.StatusNotFound, "Event not found", nil)
		return
	}

	// Parse times
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid start_time format", err)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid end_time format", err)
		return
	}

	// Validate time range
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		h.errorResponse(w, http.StatusBadRequest, "end_time must be after start_time", nil)
		return
	}

	// Update event
	event := &Event{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		CreatedAt:   existing.CreatedAt,
	}

	if updateErr := h.repo.Update(id, event); updateErr != nil {
		if updateErr == sql.ErrNoRows {
			h.errorResponse(w, http.StatusNotFound, "Event not found", nil)
			return
		}
		h.errorResponse(w, http.StatusInternalServerError, "Failed to update event", updateErr)
		return
	}

	// Get updated event
	updated, err := h.repo.Get(id)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to retrieve updated event", err)
		return
	}

	h.jsonResponse(w, http.StatusOK, updated)
}

// DeleteEvent handles DELETE /api/events/{id}
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		RecordDBOperation("delete", "events", time.Since(start))
	}()

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			h.errorResponse(w, http.StatusNotFound, "Event not found", nil)
			return
		}
		h.errorResponse(w, http.StatusInternalServerError, "Failed to delete event", err)
		return
	}

	// Record metrics
	eventsDeletedTotal.Inc()
	activeEvents.Dec()

	w.WriteHeader(http.StatusNoContent)
}

// HealthCheck handles GET /health
func (h *EventHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
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

	h.jsonResponse(w, statusCode, response)
}

// VersionCheck handles GET /version
func (h *EventHandler) VersionCheck(w http.ResponseWriter, r *http.Request) {
	response := VersionResponse{
		Version:   GetVersion(),
		BuildInfo: GetVersionInfo(),
		Timestamp: time.Now().UTC(),
	}

	h.jsonResponse(w, http.StatusOK, response)
}

// Helper methods

func (h *EventHandler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *EventHandler) errorResponse(w http.ResponseWriter, status int, message string, _ error) {
	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}

	h.jsonResponse(w, status, response)
}

func (h *EventHandler) validationErrorResponse(w http.ResponseWriter, err error) {
	response := ErrorResponse{
		Error:   "Validation Failed",
		Message: "Request validation failed",
		Details: make(map[string]string),
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			response.Details[e.Field()] = e.Tag()
		}
	}

	h.jsonResponse(w, http.StatusBadRequest, response)
}
