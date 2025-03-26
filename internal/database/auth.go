package database

import (
	"database/sql"
	"log/slog"
)

func (s *service) SaveToken(token string, username string) error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM auth").Scan(&count)
	if err != nil {
		slog.Error("Failed to check existing token", "error", err)
		return err
	}

	var sqlQuery string
	if count > 0 {
		// Update existing token
		sqlQuery = "UPDATE auth SET (token, username) = (?, ?)"
	} else {
		// Insert new token
		sqlQuery = "INSERT INTO auth (token, username) VALUES (?, ?)"
	}

	// Execute the query
	_, err = s.db.Exec(sqlQuery, token, username)
	if err != nil {
		slog.Error("Failed to save token", "error", err)
		return err
	}

	return nil
}

func (s *service) GetToken() (string, error) {
	var token string
	err := s.db.QueryRow("SELECT token FROM auth").Scan(&token)
	if err != nil {
		if err != sql.ErrNoRows {
			slog.Error("Database query error", "error", err)
		}
		return "", err
	}

	return token, nil
}
