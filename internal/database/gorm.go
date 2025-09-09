package database

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"kleio/internal/models"
)

var DB *gorm.DB

func InitializeGORM() error {
	dsn := buildPostgreSQLDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// PostgreSQL 18 has native UUID v7 support - no extension needed
	// Verify gen_random_uuid7() function is available
	var hasUUIDv7 bool
	err = db.Raw("SELECT EXISTS(SELECT 1 FROM pg_proc WHERE proname = 'gen_random_uuid7')").Scan(&hasUUIDv7).Error
	if err != nil {
		return fmt.Errorf("failed to check for UUID v7 support: %w", err)
	}
	if !hasUUIDv7 {
		return fmt.Errorf("PostgreSQL UUID v7 support not available - ensure you're using PostgreSQL 18+")
	}

	// Run custom SQL migrations first (for schema setup)
	if err := RunPostgreSQLMigrations(db); err != nil {
		return fmt.Errorf("failed to run PostgreSQL migrations: %w", err)
	}

	// Auto-migrate schemas (for any model changes not covered by SQL migrations)
	if err := autoMigrate(db); err != nil {
		return fmt.Errorf("failed to auto-migrate schemas: %w", err)
	}

	DB = db
	slog.Info("Successfully connected to PostgreSQL with GORM")
	return nil
}

func buildPostgreSQLDSN() string {
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "kleio")
	password := getEnvOrDefault("DB_PASSWORD", "")
	dbname := getEnvOrDefault("DB_NAME", "kleio")
	sslmode := getEnvOrDefault("DB_SSL_MODE", "disable")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.AuthToken{},
		&models.Release{},
		&models.Artist{},
		&models.Label{},
		&models.Genre{},
		&models.Style{},
		&models.Track{},
		&models.UserRelease{},
		&models.PlayHistory{},
		&models.Stylus{},
		&models.CleaningHistory{},
		&models.Folder{},
		&models.Sync{},
	)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func CloseGORM() error {
	if DB == nil {
		return nil
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	
	return sqlDB.Close()
}