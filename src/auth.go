package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// User represents a system user with API access
type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	APIKey    string    `json:"api_key" db:"api_key"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByAPIKey retrieves a user by their API key
func (r *UserRepository) GetByAPIKey(apiKey string) (*User, error) {
	query := `
		SELECT id, username, api_key, created_at, updated_at
		FROM users
		WHERE api_key = $1
	`

	var user User
	err := r.db.QueryRow(query, apiKey).Scan(
		&user.ID,
		&user.Username,
		&user.APIKey,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by API key: %w", err)
	}

	return &user, nil
}

// GetByUsername retrieves a user by their username
func (r *UserRepository) GetByUsername(username string) (*User, error) {
	query := `
		SELECT id, username, api_key, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.APIKey,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

// Create inserts a new user into the database
func (r *UserRepository) Create(user *User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt

	query := `
		INSERT INTO users (id, username, api_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Username,
		user.APIKey,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// AuthMiddleware provides API key authentication
type AuthMiddleware struct {
	userRepo *UserRepository
	config   *Config
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(userRepo *UserRepository, config *Config) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: userRepo,
		config:   config,
	}
}

// RequireAPIKey middleware that validates API key
func (a *AuthMiddleware) RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(a.config.APIKeyHeader)
		if apiKey == "" {
			log.Printf("❌ Missing API key in header %s", a.config.APIKeyHeader)
			errorResponse(w, http.StatusUnauthorized, "API key required", nil)
			return
		}

		// Check if it's the bootstrap admin key
		if a.config.BootstrapAdminKey != "" && apiKey == a.config.BootstrapAdminKey {
			log.Printf("✅ Bootstrap admin key used for authentication")
			// Create a virtual admin user for the context
			adminUser := &User{
				ID:       "bootstrap-admin",
				Username: "admin",
				APIKey:   a.config.BootstrapAdminKey,
			}
			r = r.WithContext(WithUser(r.Context(), adminUser))
			next.ServeHTTP(w, r)
			return
		}

		// Look up user by API key
		user, err := a.userRepo.GetByAPIKey(apiKey)
		if err != nil {
			log.Printf("❌ Failed to lookup user by API key: %v", err)
			errorResponse(w, http.StatusInternalServerError, "Authentication error", err)
			return
		}

		if user == nil {
			log.Printf("❌ Invalid API key provided")
			errorResponse(w, http.StatusUnauthorized, "Invalid API key", nil)
			return
		}

		log.Printf("✅ User authenticated: %s", user.Username)

		// Add user to request context
		r = r.WithContext(WithUser(r.Context(), user))
		next.ServeHTTP(w, r)
	})
}

// CreateBootstrapUser creates the bootstrap admin user if it doesn't exist
func (a *AuthMiddleware) CreateBootstrapUser() error {
	if a.config.BootstrapAdminKey == "" {
		log.Println("⚠️ No bootstrap admin key configured, skipping admin user creation")
		return nil
	}

	log.Println("🔑 Setting up bootstrap admin user...")

	// Check if admin user already exists
	existingUser, err := a.userRepo.GetByUsername("admin")
	if err != nil {
		return fmt.Errorf("failed to check for existing admin user: %w", err)
	}

	if existingUser != nil {
		// Update the API key if it's different
		if existingUser.APIKey != a.config.BootstrapAdminKey {
			log.Println("🔄 Updating admin user API key...")
			// In a real system, you'd want to update the user's API key
			// For now, we'll just log that it exists
		}
		log.Printf("✅ Admin user already exists: %s", existingUser.Username)
		return nil
	}

	// Create admin user
	adminUser := &User{
		Username: "admin",
		APIKey:   a.config.BootstrapAdminKey,
	}

	if err := a.userRepo.Create(adminUser); err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Printf("✅ Created bootstrap admin user with username: %s", adminUser.Username)
	return nil
}

// Helper functions for context

type contextKey string

const userContextKey contextKey = "user"

// WithUser adds a user to the context
func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetUser retrieves the user from the context
func GetUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userContextKey).(*User)
	return user, ok
}

// Helper function for error responses
func errorResponse(w http.ResponseWriter, status int, message string, _ error) {
	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}

	jsonResponse(w, status, response)
}

// jsonResponse helper function
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
