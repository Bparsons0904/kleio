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

	for _, folder := range folders {
		page := 1
		perPage := 100
		for {
			response, err := fetchReleasesPage(user, folder.ID, page, perPage)
			if err != nil {
				return err
			}

			if len(response.Releases) == 0 {
				break
			}

			err = c.DB.SaveReleases(response)
			if err != nil {
				return err
			}

			page++

			if page >= response.Pagination.Pages {
				break
			}

			time.Sleep(1 * time.Second)
		}
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

	slog.Info("Fetching releases page", "url", url)

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
		slog.Warn("Rate limited by Discogs API", "retryAfter", retryAfter)

		// Default to 60 seconds if no Retry-After header
		waitTime := 60 * time.Second
		if retryAfter != "" {
			if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
				waitTime = seconds
			}
		}

		// Wait and retry once
		slog.Info("Waiting before retry", "waitTime", waitTime)
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

	slog.Info("Successfully fetched releases page",
		"folderID", folderID,
		"page", page,
		"totalPages", response.Pagination.Pages,
		"itemsOnPage", len(response.Releases),
		"totalItems", response.Pagination.Items)

	return response, nil
}

func GetAllFolderReleases(user database.User) (map[int][]DiscogsResponse, error) {
	// First get all folders
	// folders, err := getDiscogFolders(user)
	// if err != nil {
	// 	return nil, fmt.Errorf("error fetching folders: %w", err)
	// }

	// Map to store all responses by folder ID
	allReleases := make(map[int][]DiscogsResponse)

	// Fetch releases for each folder
	// for _, folder := range folders {
	// 	slog.Info("Fetching releases for folder", "folderID", folder.ID, "name", folder.Name)
	//
	// 	// Skip folder 0 (All) since it would duplicate releases
	// 	if folder.ID == 0 {
	// 		slog.Info("Skipping 'All' folder to avoid duplicates", "folderID", folder.ID)
	// 		continue
	// 	}
	//
	// }

	return allReleases, nil
}

// GetAllReleasesAsJSON retrieves all releases with their complete related data as a JSON string
// func GetAllReleasesAsJSON(db *sql.DB) (string, error) {
// 	releases, err := GetAllReleases(db)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	jsonData, err := json.MarshalIndent(releases, "", "  ")
// 	if err != nil {
// 		return "", fmt.Errorf("error marshalling releases to JSON: %w", err)
// 	}
//
// 	return string(jsonData), nil
// }
