package database

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

var (
	dburl      = "sqlite.db"
	dbInstance *Database
)

func New() Database {
	// Reuse Connection
	if dbInstance != nil {
		return *dbInstance
	}

	// Check if we should use PostgreSQL
	if os.Getenv("DB_TYPE") == "postgres" {
		if err := InitializeGORM(); err != nil {
			log.Fatalf("Failed to initialize PostgreSQL with GORM: %v", err)
		}
		// Create wrapper for compatibility - GORM instance accessed via global DB variable
		dbInstance = &Database{DB: nil} // GORM doesn't use sql.DB directly
		slog.Info("Using PostgreSQL with GORM")
		return *dbInstance
	}

	// Default SQLite behavior
	if os.Getenv("APP_ENV") != "development" {
		dburl = getDBPath()
	}

	if err := Initialize(dburl); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	slog.Info("Connecting to database...", "dburl", dburl)
	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	dbInstance = &Database{
		DB: db,
	}
	slog.Info("Using SQLite")
	return *dbInstance
}

func Initialize(dbPath string) error {
	slog.Info("Initializing database...", "dbPath", dbPath)
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	files, err := filepath.Glob("internal/database/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to find migration files: %w", err)
	}

	sort.Strings(files)

	for _, file := range files {
		migration, err := os.ReadFile(file)
		slog.Info("Applying migration...", "file", file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		_, err = db.Exec(string(migration))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		fmt.Printf("Applied migration: %s\n", file)
	}

	return nil
}

func (s *Database) Close() error {
	// Handle GORM PostgreSQL cleanup
	if os.Getenv("DB_TYPE") == "postgres" {
		log.Printf("Disconnecting from PostgreSQL with GORM")
		return CloseGORM()
	}
	
	// Handle SQLite cleanup
	log.Printf("Disconnected from database: %s", dburl)
	return s.DB.Close()
}

func (s *Database) GetDB() *sql.DB {
	return s.DB
}

func getDBPath() string {
	if os.Getenv("APP_ENV") == "production" {
		return "/data/db/kleio.db"
	}

	if os.Getenv("APP_ENV") == "development" {
		return "kleio.db"
	}

	return "sqlite.db"
}
