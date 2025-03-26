package database

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type Service interface {
	Close() error
	GetDB() *sql.DB
	GetToken() (string, error)
	GetUser() (User, error)
	SaveToken(token string, username string) error
}

type service struct {
	db *sql.DB
}

var (
	dburl      = "sqlite.db"
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
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

	dbInstance = &service{
		db: db,
	}
	return dbInstance
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

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	return s.db.Close()
}

func (s *service) GetDB() *sql.DB {
	return s.db
}
