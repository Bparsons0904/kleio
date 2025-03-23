package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"kleio/internal/database"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type Folder struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Count       int       `json:"count"`
	ResourceURL string    `json:"resource_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastSynced  time.Time `json:"last_synced"`
}

type FoldersResponse struct {
	Folders []Folder `json:"folders"`
}

func GetFolders(user database.User) ([]Folder, error) {
	// Build the URL for the folders endpoint
	url := fmt.Sprintf(
		"%s/users/%s/collection/folders?token=%s",
		BaseURL,
		user.Username,
		user.Token,
	)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return nil, err
	}

	// Set required User-Agent header
	req.Header.Set("User-Agent", UserAgent)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to make request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("API returned non-200 status", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	// Decode the response
	var foldersResp FoldersResponse
	if err := json.NewDecoder(resp.Body).Decode(&foldersResp); err != nil {
		slog.Error("Error decoding response", "error", err)
		return nil, err
	}

	// Log the folders
	log.Printf("Found %d folders", len(foldersResp.Folders))
	// for _, folder := range foldersResp.Folders {
	// 	log.Printf("Folder: ID=%d, Name=%s, Count=%d", folder.ID, folder.Name, folder.Count)
	// }

	return foldersResp.Folders, nil
}

func updateFolders(db *sql.DB, folders []Folder) {
	now := time.Now().Format(time.RFC3339)
	for _, folder := range folders {
		updateFolder(db, folder, now)
	}
	slog.Info("Folders updated")

	// TODO: Check if any folders have been deleted
}

func updateFolder(db *sql.DB, folder Folder, now string) {
	_, err := db.Exec(
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

func getLocalFolderLastSynced(db *sql.DB) (time.Time, error) {
	var lastSynced time.Time
	err := db.QueryRow("SELECT last_synced FROM folders ORDER BY last_synced ASC LIMIT 1").
		Scan(&lastSynced)
	if err != nil {
		if err != sql.ErrNoRows {
			slog.Error("Database query error", "error", err)
		}
		return time.Time{}, err
	}

	return lastSynced, nil
}
