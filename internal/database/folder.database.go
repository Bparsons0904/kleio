package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

func (s *Database) GetLastFolderSync() (time.Time, error) {
	// var lastSynced time.Time
	// err := s.DB.QueryRow("SELECT last_synced FROM folders ORDER BY last_synced ASC LIMIT 1").
	// 	Scan(&lastSynced)
	// if err != nil {
	// 	if err != sql.ErrNoRows {
	// 		slog.Error("Database query error", "error", err)
	// 	}
	// 	return time.Time{}, err
	// }
	//
	// return lastSynced, nil

	return s.GetLastSync("folders")
}

func (s *Database) GetLastReleaseSync() (time.Time, error) {
	return s.GetLastSync("releases")
}

func (s *Database) GetLastSync(table string) (time.Time, error) {
	query := fmt.Sprintf("SELECT last_synced FROM %s ORDER BY last_synced ASC LIMIT 1", table)
	var lastSynced time.Time
	err := s.DB.QueryRow(query).Scan(&lastSynced)
	if err != nil {
		if err != sql.ErrNoRows {
			slog.Error("Database query error", "error", err, "table", table)
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
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan folder", "error", err)
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

func (database *Database) UpdateFolder(folder Folder, now string) error {
	_, err := dbInstance.DB.Exec(
		"INSERT OR REPLACE INTO folders (id, name, count, resource_url, last_synced) VALUES (?, ?, ?, ?, ?)",
		folder.ID,
		folder.Name,
		folder.Count,
		folder.ResourceURL,
		now,
	)
	if err != nil {
		slog.Error("Failed to update folder", "error", err, "folder", folder)
	}

	return err
}
