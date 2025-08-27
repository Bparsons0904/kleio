package controller

import (
	"encoding/json"
	"fmt"
	"io"
	. "kleio/internal/database"
	"log/slog"
	"net/http"
	"time"
)

func (c *Controller) SyncFolders() error {
	slog.Info("Starting folder sync")

	user, err := c.DB.GetUser()
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		return err
	}

	slog.Debug("Retrieved user for folder sync", "username", user.Username)

	folders, err := c.getDiscogFolders(user)
	if err != nil {
		slog.Error("Failed to get user folders from Discogs API", "error", err)
		return err
	}

	slog.Info("Retrieved folders from Discogs", "folderCount", len(folders))

	err = c.updateFolders(folders)
	if err != nil {
		slog.Error("Failed to update folders in database", "error", err, "folderCount", len(folders))
		return err
	}

	slog.Info("Folder sync completed successfully", "folderCount", len(folders))
	return nil
}

func (c *Controller) getDiscogFolders(user User) ([]Folder, error) {
	url := fmt.Sprintf(
		"%s/users/%s/collection/folders?token=%s",
		BaseURL,
		user.Username,
		user.Token,
	)

	slog.Debug("Making API request for folders", "url", url, "username", user.Username)

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

	// Check for rate limiting
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		slog.Warn("Rate limited while fetching folders", 
			"retryAfter", retryAfter,
			"username", user.Username)
		return nil, fmt.Errorf("rate limited while fetching folders: retry after %s seconds", retryAfter)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("API returned non-200 status for folders", 
			"status", resp.StatusCode, 
			"body", string(body),
			"username", user.Username,
			"url", url)
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	var foldersResp FoldersResponse
	if err := json.NewDecoder(resp.Body).Decode(&foldersResp); err != nil {
		slog.Error("Error decoding response", "error", err)
		return nil, err
	}

	slog.Info("Successfully retrieved folders from API", 
		"folderCount", len(foldersResp.Folders),
		"username", user.Username)

	for i, folder := range foldersResp.Folders {
		slog.Debug("Retrieved folder", 
			"index", i,
			"folderID", folder.ID,
			"folderName", folder.Name,
			"itemCount", folder.Count)
	}

	return foldersResp.Folders, nil
}

func (c *Controller) updateFolders(folders []Folder) error {
	updated := 0
	failed := 0

	for i, folder := range folders {
		slog.Debug("Updating folder", 
			"index", i,
			"folderID", folder.ID,
			"folderName", folder.Name)

		if err := c.DB.UpdateFolder(folder); err != nil {
			slog.Error("Failed to update folder", 
				"error", err,
				"folderID", folder.ID,
				"folderName", folder.Name)
			failed++
			continue // Continue with other folders
		}
		updated++
	}

	slog.Info("Folder update completed", 
		"totalFolders", len(folders),
		"updated", updated,
		"failed", failed)

	if failed > 0 && updated == 0 {
		return fmt.Errorf("failed to update any folders: %d failures", failed)
	}

	return nil

	// TODO: Check if any folders have been deleted
}
