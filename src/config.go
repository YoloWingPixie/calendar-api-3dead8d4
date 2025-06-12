package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Host string
	Port string

	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Application configuration
	Debug              bool
	BootstrapAdminKey  string
	APIKeyHeader       string
	Environment        string
	DopplerProject     string
	DopplerEnvironment string
	DopplerConfig      string
}

// LoadConfig loads configuration from environment variables and Doppler secrets
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Load from DOPPLER_SECRETS_JSON if available
	if err := config.loadFromDopplerSecrets(); err != nil {
		log.Printf("Warning: Failed to load Doppler secrets: %v", err)
	}

	// Load from individual environment variables (fallback)
	config.loadFromEnv()

	// Validate required configuration
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Log configuration (without sensitive data)
	config.logConfig()

	return config, nil
}

// loadFromDopplerSecrets loads configuration from DOPPLER_SECRETS_JSON
func (c *Config) loadFromDopplerSecrets() error {
	secretsJSON := os.Getenv("DOPPLER_SECRETS_JSON")
	if secretsJSON == "" {
		return fmt.Errorf("DOPPLER_SECRETS_JSON not found")
	}

	log.Println("🔐 Loading configuration from Doppler secrets...")

	var secrets map[string]interface{}
	if err := json.Unmarshal([]byte(secretsJSON), &secrets); err != nil {
		return fmt.Errorf("failed to parse Doppler secrets JSON: %w", err)
	}

	// Load configuration from secrets
	c.Host = getStringFromSecrets(secrets, "TF_VAR_host", "0.0.0.0")
	c.Port = getStringFromSecrets(secrets, "TF_VAR_port", "8080")

	// Database configuration
	c.DBHost = getStringFromSecrets(secrets, "TF_VAR_database_host", "localhost")
	c.DBPort = getStringFromSecrets(secrets, "TF_VAR_database_port", "5432")
	c.DBUser = getStringFromSecrets(secrets, "TF_VAR_database_username", "postgres")
	c.DBPassword = getStringFromSecrets(secrets, "TF_VAR_database_password", "")
	c.DBName = getStringFromSecrets(secrets, "TF_VAR_database_name", "calendar")
	c.DBSSLMode = "require" // Default for production

	// Application configuration
	c.Debug = getBoolFromSecrets(secrets, "TF_VAR_debug", false)
	c.BootstrapAdminKey = getStringFromSecrets(secrets, "TF_VAR_bootstrap_admin_key", "")
	c.APIKeyHeader = getStringFromSecrets(secrets, "TF_VAR_api_key_header", "X-API-Key")
	c.Environment = getStringFromSecrets(secrets, "TF_VAR_environment", "development")
	c.DopplerProject = getStringFromSecrets(secrets, "TF_VAR_doppler_project", "")
	c.DopplerEnvironment = getStringFromSecrets(secrets, "TF_VAR_doppler_environment", "")
	c.DopplerConfig = getStringFromSecrets(secrets, "TF_VAR_doppler_config", "")

	log.Printf("✅ Loaded configuration for environment: %s", c.Environment)
	return nil
}

// loadFromEnv loads configuration from individual environment variables (fallback)
func (c *Config) loadFromEnv() {
	log.Println("🔧 Loading configuration from environment variables...")

	// Only override if not already set by Doppler
	if c.Host == "" {
		c.Host = getEnv("HOST", "0.0.0.0")
	}
	if c.Port == "" {
		c.Port = getEnv("PORT", "8080")
	}

	// Database configuration
	if c.DBHost == "" {
		c.DBHost = getEnv("DB_HOST", "localhost")
	}
	if c.DBPort == "" {
		c.DBPort = getEnv("DB_PORT", "5432")
	}
	if c.DBUser == "" {
		c.DBUser = getEnv("DB_USER", "postgres")
	}
	if c.DBPassword == "" {
		c.DBPassword = getEnv("DB_PASSWORD", "postgres")
	}
	if c.DBName == "" {
		c.DBName = getEnv("DB_NAME", "calendar")
	}
	if c.DBSSLMode == "" {
		c.DBSSLMode = getEnv("DB_SSLMODE", "disable")
	}

	// Application configuration
	if !c.Debug {
		c.Debug = getEnv("DEBUG", "false") == "true"
	}
	if c.BootstrapAdminKey == "" {
		c.BootstrapAdminKey = getEnv("BOOTSTRAP_ADMIN_KEY", "")
	}
	if c.APIKeyHeader == "" {
		c.APIKeyHeader = getEnv("API_KEY_HEADER", "X-API-Key")
	}
	if c.Environment == "" {
		c.Environment = getEnv("ENVIRONMENT", "development")
	}
}

// validate ensures all required configuration is present
func (c *Config) validate() error {
	if c.DBHost == "" {
		return fmt.Errorf("database host is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("database user is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	return nil
}

// logConfig logs the current configuration (without sensitive data)
func (c *Config) logConfig() {
	log.Println("📋 Application Configuration:")
	log.Printf("  Environment: %s", c.Environment)
	log.Printf("  Host: %s", c.Host)
	log.Printf("  Port: %s", c.Port)
	log.Printf("  Database Host: %s", c.DBHost)
	log.Printf("  Database Name: %s", c.DBName)
	log.Printf("  Database User: %s", c.DBUser)
	log.Printf("  Database SSL Mode: %s", c.DBSSLMode)
	log.Printf("  Debug Mode: %t", c.Debug)
	log.Printf("  API Key Header: %s", c.APIKeyHeader)

	if c.BootstrapAdminKey != "" {
		log.Printf("  Bootstrap Admin Key: %s***", c.BootstrapAdminKey[:8])
	}

	if c.DopplerProject != "" {
		log.Printf("  Doppler Project: %s", c.DopplerProject)
		log.Printf("  Doppler Environment: %s", c.DopplerEnvironment)
		log.Printf("  Doppler Config: %s", c.DopplerConfig)
	}
}

// Helper functions

func getStringFromSecrets(secrets map[string]interface{}, key, defaultValue string) string {
	if val, ok := secrets[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getBoolFromSecrets(secrets map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := secrets[key]; ok {
		if str, ok := val.(string); ok {
			if parsed, err := strconv.ParseBool(str); err == nil {
				return parsed
			}
		}
	}
	return defaultValue
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
