package database

import (
	"log/slog"
)

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

func (database *Database) UpdateFolder(folder Folder) error {
	_, err := dbInstance.DB.Exec(
		"INSERT OR REPLACE INTO folders (id, name, count, resource_url) VALUES (?, ?, ?, ?)",
		folder.ID,
		folder.Name,
		folder.Count,
		folder.ResourceURL,
	)
	if err != nil {
		slog.Error("Failed to update folder", "error", err, "folder", folder)
	}

	return err
}
