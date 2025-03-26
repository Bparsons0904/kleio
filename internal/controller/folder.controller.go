package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"kleio/internal/database"
	. "kleio/internal/database"
	"log"
	"log/slog"
	"net/http"
	"time"
)

func (c *Controller) SyncFolders() {
	folders, err := c.getDiscogFolders()
	if err != nil {
		slog.Error("Failed to get user folders", "error", err)
		return
	}

	c.updateFolders(folders)
}

func (c *Controller) getDiscogFolders() ([]Folder, error) {
	url := fmt.Sprintf(
		"%s/users/%s/collection/folders?token=%s",
		BaseURL,
		c.User.Username,
		c.User.Token,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to make request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("API returned non-200 status", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	var foldersResp FoldersResponse
	if err := json.NewDecoder(resp.Body).Decode(&foldersResp); err != nil {
		slog.Error("Error decoding response", "error", err)
		return nil, err
	}

	log.Printf("Found %d folders", len(foldersResp.Folders))
	return foldersResp.Folders, nil
}

func (c *Controller) updateFolders(folders []Folder) {
	now := time.Now().Format(time.RFC3339)
	for _, folder := range folders {
		c.DB.UpdateFolder(folder, now)
	}
	slog.Info("Folders updated")

	// TODO: Check if any folders have been deleted
}

// func getLocalFolderLastSynced(db *sql.DB) (time.Time, error) {
// 	var lastSynced time.Time
// 	err := db.QueryRow("SELECT last_synced FROM folders ORDER BY last_synced ASC LIMIT 1").
// 		Scan(&lastSynced)
// 	if err != nil {
// 		if err != sql.ErrNoRows {
// 			slog.Error("Database query error", "error", err)
// 		}
// 		return time.Time{}, err
// 	}
//
// 	return lastSynced, nil
// }

func updateCollectionByFolder(user database.User, db *sql.DB, folder Folder) {
	// Query collection by folder

	// url := fmt.Sprintf(
	// 	"%s/users/%s/collection/folders/%d/releases?token=%s",
	// 	BaseURL,
	// 	user.Username,
	// 	folder.ID,
	// 	user.Token,
	// )
	//
	// // Create a new request
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	slog.Error("Failed to create request", "error", err)
	// 	return
	// }
	//
	// // Set required User-Agent header
	// req.Header.Set("User-Agent", UserAgent)
	//
	// // Create HTTP client with timeout
	// client := &http.Client{
	// 	Timeout: 10 * time.Second,
	// }
	//
	// // Make the request
	// resp, err := client.Do(req)
	// if err != nil {
	// 	slog.Error("Failed to make request", "error", err)
	// 	return
	// }
	// defer resp.Body.Close()
	//
	// // Check response status
	// if resp.StatusCode != http.StatusOK {
	// 	body, _ := io.ReadAll(resp.Body)
	// 	slog.Error("API returned non-200 status", "status", resp.StatusCode, "body", string(body))
	// 	return
	// }

	collection, err := fetchReleasesPage(user, folder.ID, 1, 10)
	if err != nil {
		slog.Error("Failed to fetch collection", "error", err)
		return
	}
	log.Printf("Fetched %d releases for folder %d", len(collection.Releases), folder.ID)

	for _, release := range collection.Releases {
		log.Printf(
			"Release: ID=%d, Title=%s, Year=%d",
			release.ID,
			release.BasicInfo.Title,
			release.BasicInfo.Year,
		)
	}

	// SaveReleases(db, collection)

	// Decode the response
	// var collection CollectionResponse
	// if err := json.NewDecoder(resp.Body).Decode(&collection); err != nil {
	// 	slog.Error("Error decoding response", "error", err)
	// 	return
	// }

	// Save collection to DB
}
