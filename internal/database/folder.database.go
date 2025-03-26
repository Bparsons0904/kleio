package database

import "log/slog"

func (database *Database) UpdateFolder(folder Folder, now string) {
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
}
