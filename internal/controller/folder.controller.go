package controller

import (
	"encoding/json"
	"fmt"
	"io"
	. "kleio/internal/database"
	"log"
	"log/slog"
	"net/http"
	"time"
)

func (c *Controller) SyncFolders() error {
	user, err := c.DB.GetUser()
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		return err
	}

	folders, err := c.getDiscogFolders(user)
	if err != nil {
		slog.Error("Failed to get user folders", "error", err)
		return err
	}

	err = c.updateFolders(folders)
	if err != nil {
		slog.Error("Failed to update folders", "error", err)
		return err
	}

	return nil
}

func (c *Controller) getDiscogFolders(user User) ([]Folder, error) {
	url := fmt.Sprintf(
		"%s/users/%s/collection/folders?token=%s",
		BaseURL,
		user.Username,
		user.Token,
	)

	slog.Info("Fetching folders", "url", url)

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

func (c *Controller) updateFolders(folders []Folder) error {
	for _, folder := range folders {
		if err := c.DB.UpdateFolder(folder); err != nil {
			slog.Error("Failed to update folder", "error", err)
			return err
		}
	}
	slog.Info("Folders updated")

	return nil

	// TODO: Check if any folders have been deleted
}
