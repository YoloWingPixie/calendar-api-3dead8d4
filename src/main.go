package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("🚀 Starting Calendar API...")
	log.Printf("📋 %s", GetVersionInfo())

	// Load environment variables from .env file (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("📝 No .env file found, using environment variables")
	}

	// Load configuration from Doppler secrets or environment variables
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database connection
	db, err := NewDB(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations with proper tracking and logging
	migrationManager := NewMigrationManager(db)

	// Print current migration status
	if err := migrationManager.PrintMigrationStatus(); err != nil {
		log.Printf("Warning: Failed to print migration status: %v", err)
	}

	// Run pending migrations
	if err := migrationManager.RunMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	userRepo := NewUserRepository(db)
	eventRepo := NewEventRepository(db)

	// Initialize authentication middleware
	authMiddleware := NewAuthMiddleware(userRepo, config)

	// Create bootstrap admin user if needed
	if err := authMiddleware.CreateBootstrapUser(); err != nil {
		log.Fatal("Failed to create bootstrap user:", err)
	}

	// Initialize handlers
	eventHandler := NewEventHandler(eventRepo)

	// Setup routes
	r := mux.NewRouter()

	// Public routes (no authentication required)
	r.HandleFunc("/health", eventHandler.HealthCheck).Methods("GET")
	r.HandleFunc("/version", eventHandler.VersionCheck).Methods("GET")

	// Protected API routes (require authentication)
	api := r.PathPrefix("/api").Subrouter()
	api.Use(authMiddleware.RequireAPIKey)
	api.HandleFunc("/events", eventHandler.ListEvents).Methods("GET")
	api.HandleFunc("/events", eventHandler.CreateEvent).Methods("POST")
	api.HandleFunc("/events/{id}", eventHandler.GetEvent).Methods("GET")
	api.HandleFunc("/events/{id}", eventHandler.UpdateEvent).Methods("PUT")
	api.HandleFunc("/events/{id}", eventHandler.DeleteEvent).Methods("DELETE")

	// Middleware
	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)

	// Server configuration
	addr := config.Host + ":" + config.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🌐 Starting server on %s (environment: %s)", addr, config.Environment)
		if config.BootstrapAdminKey != "" {
			log.Printf("🔑 Bootstrap admin key is configured for API access")
		}
		log.Printf("🔒 API endpoints require authentication via %s header", config.APIKeyHeader)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

// Middleware functions
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
