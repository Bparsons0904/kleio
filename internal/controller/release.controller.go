package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"kleio/internal/database"
	. "kleio/internal/database"
	"log/slog"
	"net/http"
	"time"
)

func (c *Controller) SyncReleases() error {
	user, err := c.DB.GetUser()
	if err != nil {
		slog.Error("Failed to get user from database", "error", err)
		return err
	}

	folders, err := c.DB.GetFolders()
	if err != nil {
		slog.Error("Failed to get user folders from database", "error", err)
		return err
	}

	slog.Info("Starting release sync", 
		"username", user.Username,
		"folderCount", len(folders))

	totalReleases := 0
	totalPages := 0
	failedFolders := 0

	for folderIdx, folder := range folders {
		slog.Info("Syncing folder", 
			"folderID", folder.ID, 
			"folderName", folder.Name,
			"progress", fmt.Sprintf("%d/%d", folderIdx+1, len(folders)))

		folderReleases := 0
		page := 1
		perPage := 100

		for {
			slog.Debug("Fetching releases page", 
				"folderID", folder.ID,
				"page", page,
				"perPage", perPage)

			response, err := fetchReleasesPage(user, folder.ID, page, perPage)
			if err != nil {
				slog.Error("Failed to fetch releases page", 
					"error", err,
					"folderID", folder.ID,
					"folderName", folder.Name,
					"page", page)
				failedFolders++
				break // Move to next folder instead of failing completely
			}

			if len(response.Releases) == 0 {
				slog.Debug("No releases found on page", 
					"folderID", folder.ID,
					"page", page)
				break
			}

			err = c.DB.SaveReleases(response)
			if err != nil {
				slog.Error("Failed to save releases", 
					"error", err,
					"folderID", folder.ID,
					"page", page,
					"releaseCount", len(response.Releases))
				return err // Database save failure should stop sync
			}

			folderReleases += len(response.Releases)
			totalReleases += len(response.Releases)
			totalPages++

			slog.Debug("Saved releases page", 
				"folderID", folder.ID,
				"page", page,
				"releasesOnPage", len(response.Releases),
				"totalPagesInFolder", response.Pagination.Pages)

			page++

			if page > response.Pagination.Pages {
				break
			}

			time.Sleep(1 * time.Second)
		}

		slog.Info("Completed folder sync", 
			"folderID", folder.ID,
			"folderName", folder.Name,
			"releasesInFolder", folderReleases)
	}

	slog.Info("Release sync completed", 
		"totalReleases", totalReleases,
		"totalPages", totalPages,
		"successfulFolders", len(folders)-failedFolders,
		"failedFolders", failedFolders)

	if failedFolders > 0 && totalReleases == 0 {
		return fmt.Errorf("failed to sync any folders: %d failures", failedFolders)
	}

	return nil
}

func fetchReleasesPage(user database.User, folderID, page, perPage int) (DiscogsResponse, error) {
	var response DiscogsResponse

	// Build the URL for the folder releases endpoint with pagination
	url := fmt.Sprintf(
		"%s/users/%s/collection/folders/%d/releases?token=%s&page=%d&per_page=%d",
		BaseURL,
		user.Username,
		folderID,
		user.Token,
		page,
		perPage,
	)

	slog.Debug("Making API request", "url", url)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return response, err
	}

	// Set required User-Agent header
	req.Header.Set("User-Agent", UserAgent)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second, // Longer timeout for pagination requests
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to make request", "error", err, "url", url)
		return response, err
	}
	defer resp.Body.Close()

	// Check for rate limiting
	if resp.StatusCode == http.StatusTooManyRequests {
		// Get retry after header if available
		retryAfter := resp.Header.Get("Retry-After")
		slog.Warn("Rate limited by Discogs API", 
			"retryAfter", retryAfter,
			"folderID", folderID,
			"page", page,
			"url", url)

		// Default to 60 seconds if no Retry-After header
		waitTime := 60 * time.Second
		if retryAfter != "" {
			if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
				waitTime = seconds
			}
		}

		// Wait and retry once
		slog.Info("Waiting before retry due to rate limit", 
			"waitTime", waitTime,
			"folderID", folderID,
			"page", page)
		time.Sleep(waitTime)
		return fetchReleasesPage(user, folderID, page, perPage)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("API returned non-200 status",
			"status", resp.StatusCode,
			"body", string(body),
			"url", url)
		return response, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	// Decode the response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Error("Error decoding response", "error", err)
		return response, err
	}

	slog.Debug("Successfully fetched releases page",
		"folderID", folderID,
		"page", page,
		"totalPages", response.Pagination.Pages,
		"itemsOnPage", len(response.Releases),
		"totalItems", response.Pagination.Items)

	return response, nil
}

func (c *Controller) DeleteRelease(releaseID int) (payload Payload, err error) {
	err = c.DB.DeleteRelease(releaseID)
	if err != nil {
		slog.Error("Failed to delete release", "error", err)
		return
	}

	err = payload.GetPayload(c)
	if err != nil {
		slog.Error("Failed to get payload for play history", "error", err)
	}

	return
}

func (c *Controller) ArchiveRelease(releaseID int) (payload Payload, err error) {
	err = c.DB.ArchiveRelease(releaseID)
	if err != nil {
		slog.Error("Failed to archive release", "error", err)
		return
	}

	err = payload.GetPayload(c)
	if err != nil {
		slog.Error("Failed to get payload for play history", "error", err)
	}

	return
}
