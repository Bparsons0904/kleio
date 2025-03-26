package database

import (
	"database/sql"
	"log/slog"
	"time"
)

func (s *Database) GetLastFolderSync() (time.Time, error) {
	var lastSynced time.Time
	err := s.DB.QueryRow("SELECT last_synced FROM folders ORDER BY last_synced ASC LIMIT 1").
		Scan(&lastSynced)
	if err != nil {
		if err != sql.ErrNoRows {
			slog.Error("Database query error", "error", err)
		}
		return time.Time{}, err
	}

	return lastSynced, nil
}

func (s *Database) GetFolders() ([]Folder, error) {
	var folders []Folder
	rows, err := s.DB.Query("SELECT * FROM folders")
	if err != nil {
		slog.Error("Failed to get folders", "error", err)
		return nil, err
	}

	for rows.Next() {
		var folder Folder
		err := rows.Scan(
			&folder.ID,
			&folder.Name,
			&folder.Count,
			&folder.ResourceURL,
			&folder.LastSynced,
		)
		if err != nil {
			slog.Error("Failed to scan folder", "error", err)
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}
