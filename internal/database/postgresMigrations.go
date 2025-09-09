package database

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"

	"gorm.io/gorm"
)

// RunPostgreSQLMigrations runs SQL migrations for PostgreSQL
// This is used as a fallback or for custom migrations that can't be handled by GORM AutoMigrate
func RunPostgreSQLMigrations(db *gorm.DB) error {
	slog.Info("Running PostgreSQL migrations...")

	// Check if migrations table exists, create if not
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrationFiles, err := filepath.Glob("internal/database/migrations/postgres/*.sql")
	if err != nil {
		return fmt.Errorf("failed to find PostgreSQL migration files: %w", err)
	}

	if len(migrationFiles) == 0 {
		slog.Info("No PostgreSQL migration files found, skipping SQL migrations")
		return nil
	}

	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		migrationName := filepath.Base(file)

		// Check if migration has already been applied
		var count int64
		err := db.Table("schema_migrations").
			Where("migration = ?", migrationName).
			Count(&count).
			Error
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", migrationName, err)
		}

		if count > 0 {
			slog.Info("Migration already applied, skipping", "migration", migrationName)
			continue
		}

		// Read and execute migration
		migration, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		slog.Info("Applying migration", "migration", migrationName)
		err = db.Exec(string(migration)).Error
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationName, err)
		}

		// Record migration as applied
		err = db.Exec(
			"INSERT INTO schema_migrations (migration, applied_at) VALUES (?, NOW())",
			migrationName,
		).Error
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migrationName, err)
		}

		slog.Info("Migration applied successfully", "migration", migrationName)
	}

	slog.Info("PostgreSQL migrations completed successfully")
	return nil
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist
func createMigrationsTable(db *gorm.DB) error {
	createSQL := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			migration VARCHAR(255) UNIQUE NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`

	return db.Exec(createSQL).Error
}

