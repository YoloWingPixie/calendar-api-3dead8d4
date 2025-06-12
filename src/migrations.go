package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

// Migration represents a database migration
type Migration struct {
	Version     string
	Description string
	SQL         string
	Applied     bool
	AppliedAt   *time.Time
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db *DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *DB) *MigrationManager {
	return &MigrationManager{db: db}
}

// GetMigrations returns all available migrations in order
func (m *MigrationManager) GetMigrations() []Migration {
	return []Migration{
		{
			Version:     "001",
			Description: "Create schema_migrations table",
			SQL: `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version VARCHAR(255) PRIMARY KEY,
				description VARCHAR(500) NOT NULL,
				applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
				checksum VARCHAR(64) NOT NULL
			);
			`,
		},
		{
			Version:     "002",
			Description: "Create events table",
			SQL: `
			CREATE TABLE IF NOT EXISTS events (
				id VARCHAR(36) PRIMARY KEY,
				title VARCHAR(255) NOT NULL,
				description TEXT,
				start_time TIMESTAMP WITH TIME ZONE NOT NULL,
				end_time TIMESTAMP WITH TIME ZONE NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
				CHECK (end_time > start_time)
			);
			`,
		},
		{
			Version:     "003",
			Description: "Create indexes on events table",
			SQL: `
			CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);
			CREATE INDEX IF NOT EXISTS idx_events_end_time ON events(end_time);
			CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at);
			`,
		},
		{
			Version:     "004",
			Description: "Create users table for API authentication",
			SQL: `
			CREATE TABLE IF NOT EXISTS users (
				id VARCHAR(36) PRIMARY KEY,
				username VARCHAR(100) UNIQUE NOT NULL,
				api_key VARCHAR(255) UNIQUE NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
			);
			
			CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
			CREATE INDEX IF NOT EXISTS idx_users_api_key ON users(api_key);
			`,
		},
		{
			Version:     "005",
			Description: "Insert test events",
			SQL: `
			INSERT INTO events (id, title, description, start_time, end_time)
			VALUES 
				('550e8400-e29b-41d4-a716-446655440001', 'Team Standup', 'Daily team standup meeting', '2025-06-13 09:00:00+00', '2025-06-13 09:30:00+00'),
				('550e8400-e29b-41d4-a716-446655440002', 'Sprint Planning', 'Sprint planning session for next iteration', '2025-06-13 10:00:00+00', '2025-06-13 12:00:00+00'),
				('550e8400-e29b-41d4-a716-446655440003', 'Lunch & Learn', 'Tech talk about microservices architecture', '2025-06-13 12:30:00+00', '2025-06-13 13:30:00+00'),
				('550e8400-e29b-41d4-a716-446655440004', 'Code Review', 'Review pull requests and discuss implementation', '2025-06-13 14:00:00+00', '2025-06-13 15:00:00+00'),
				('550e8400-e29b-41d4-a716-446655440005', 'End of Day Sync', 'Quick sync to wrap up the day', '2025-06-13 17:00:00+00', '2025-06-13 17:15:00+00');
			`,
		},
	}
}

// RunMigrations executes all pending migrations
func (m *MigrationManager) RunMigrations() error {
	log.Println("🔄 Starting database migrations...")

	migrations := m.GetMigrations()
	if len(migrations) == 0 {
		log.Println("✅ No migrations to run")
		return nil
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	log.Printf("📊 Total migrations: %d, Applied: %d", len(migrations), len(appliedMigrations))

	pendingCount := 0
	for _, migration := range migrations {
		if _, applied := appliedMigrations[migration.Version]; !applied {
			pendingCount++
			if err := m.applyMigration(migration); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}
		} else {
			log.Printf("⏭️  Migration %s (%s) already applied", migration.Version, migration.Description)
		}
	}

	if pendingCount == 0 {
		log.Println("✅ All migrations are up to date")
	} else {
		log.Printf("✅ Successfully applied %d pending migrations", pendingCount)
	}

	return nil
}

// getAppliedMigrations returns a map of applied migration versions
func (m *MigrationManager) getAppliedMigrations() (map[string]bool, error) {
	// First, ensure the schema_migrations table exists
	firstMigration := m.GetMigrations()[0] // The schema_migrations creation
	if err := m.executeSQL(firstMigration.SQL, firstMigration.Version, firstMigration.Description); err != nil {
		return nil, fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	query := `
		SELECT version, description, applied_at 
		FROM schema_migrations 
		ORDER BY version
	`

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version, description string
		var appliedAt time.Time

		if err := rows.Scan(&version, &description, &appliedAt); err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}

		applied[version] = true
		log.Printf("📋 Previously applied: %s (%s) at %s", version, description, appliedAt.Format(time.RFC3339))
	}

	return applied, nil
}

// applyMigration applies a single migration
func (m *MigrationManager) applyMigration(migration Migration) error {
	log.Printf("🚀 Applying migration %s: %s", migration.Version, migration.Description)

	// Start a transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Warning: failed to rollback transaction: %v", err)
		}
	}()

	// Execute the migration SQL
	start := time.Now()
	if err := m.executeSQLInTx(tx, migration.SQL, migration.Version, migration.Description); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}
	duration := time.Since(start)

	// Skip recording the first migration (schema_migrations table creation) to avoid recursion
	if migration.Version != "001" {
		// Record the migration as applied
		checksum := m.calculateChecksum(migration.SQL)
		recordSQL := `
			INSERT INTO schema_migrations (version, description, applied_at, checksum)
			VALUES ($1, $2, NOW(), $3)
			ON CONFLICT (version) DO NOTHING
		`

		if _, err := tx.Exec(recordSQL, migration.Version, migration.Description, checksum); err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration transaction: %w", err)
	}

	log.Printf("✅ Migration %s completed successfully in %v", migration.Version, duration)
	return nil
}

// executeSQL executes SQL with proper error handling and logging
func (m *MigrationManager) executeSQL(sql, version, description string) error {
	log.Printf("🔧 Executing SQL for migration %s (%s)", version, description)

	// Clean up the SQL for logging (remove extra whitespace)
	cleanSQL := strings.TrimSpace(strings.ReplaceAll(sql, "\n", " "))
	if len(cleanSQL) > 200 {
		cleanSQL = cleanSQL[:200] + "..."
	}
	log.Printf("📝 SQL: %s", cleanSQL)

	if _, err := m.db.Exec(sql); err != nil {
		log.Printf("❌ Failed to execute SQL for migration %s: %v", version, err)
		return err
	}

	return nil
}

// executeSQLInTx executes SQL within a transaction
func (m *MigrationManager) executeSQLInTx(tx *sql.Tx, sql, version, description string) error {
	log.Printf("🔧 Executing SQL for migration %s (%s) in transaction", version, description)

	// Clean up the SQL for logging (remove extra whitespace)
	cleanSQL := strings.TrimSpace(strings.ReplaceAll(sql, "\n", " "))
	if len(cleanSQL) > 200 {
		cleanSQL = cleanSQL[:200] + "..."
	}
	log.Printf("📝 SQL: %s", cleanSQL)

	if _, err := tx.Exec(sql); err != nil {
		log.Printf("❌ Failed to execute SQL for migration %s: %v", version, err)
		return err
	}

	return nil
}

// calculateChecksum creates a simple checksum for the migration SQL
func (m *MigrationManager) calculateChecksum(sql string) string {
	// Simple checksum - in production you might want to use SHA256
	return fmt.Sprintf("%x", len(sql))
}

// GetMigrationStatus returns the status of all migrations
func (m *MigrationManager) GetMigrationStatus() ([]Migration, error) {
	migrations := m.GetMigrations()
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Mark which migrations have been applied
	for i := range migrations {
		if _, applied := appliedMigrations[migrations[i].Version]; applied {
			migrations[i].Applied = true
		}
	}

	return migrations, nil
}

// PrintMigrationStatus prints a detailed status of all migrations
func (m *MigrationManager) PrintMigrationStatus() error {
	log.Println("📊 Migration Status Report")
	log.Println("=" + strings.Repeat("=", 80))

	migrations, err := m.GetMigrationStatus()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		status := "❌ PENDING"
		if migration.Applied {
			status = "✅ APPLIED"
		}

		log.Printf("%s | %s | %s", status, migration.Version, migration.Description)
	}

	appliedCount := 0
	for _, m := range migrations {
		if m.Applied {
			appliedCount++
		}
	}

	log.Println("=" + strings.Repeat("=", 80))
	log.Printf("📈 Summary: %d/%d migrations applied", appliedCount, len(migrations))

	return nil
}
