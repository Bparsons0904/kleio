package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"kleio/internal/database"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (c *Controller) GetReleaseDetails(resourceURL string, token string) error {
	resourceURL += "?token=" + token

	// Create a new request
	req, err := http.NewRequest("GET", resourceURL, nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return err
	}

	req.Header.Set("User-Agent", UserAgent)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to make request", "error", err, "url", resourceURL)
		return err
	}
	defer resp.Body.Close()

	// Check for rate limiting
	if resp.StatusCode == http.StatusTooManyRequests {
		// Get retry after header if available
		retryAfter := resp.Header.Get("Retry-After")
		slog.Warn("Rate limited by Discogs API", "retryAfter", retryAfter)
		return fmt.Errorf("rate limited by Discogs API")
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("API returned non-200 status",
			"status", resp.StatusCode,
			"body", string(body),
			"url", resourceURL)
		return fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	// Decode the response
	var releaseDetails struct {
		ID        int                     `json:"id"`
		Tracklist []database.DiscogsTrack `json:"tracklist"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&releaseDetails); err != nil {
		slog.Error("Error decoding response", "error", err)
		return err
	}

	log.Println("Fetching release details...", resourceURL)
	log.Printf("Fetched release details for %v\n", releaseDetails)

	// Process tracks and calculate durations
	var tracks []database.Track
	totalDurationSeconds := 0
	isDurationEstimated := false

	for _, discTrack := range releaseDetails.Tracklist {
		// Skip non-track items (headers, etc.)
		if discTrack.Type != "track" && discTrack.Type != "" {
			continue
		}

		durationSeconds, err := convertDurationToSeconds(discTrack.Duration)
		if err != nil {
			slog.Error("Failed to parse track duration",
				"error", err,
				"track", discTrack.Title,
				"duration", discTrack.Duration)
			durationSeconds = 0
		}

		// Only count tracks with valid durations
		if durationSeconds > 0 {
			totalDurationSeconds += durationSeconds
		}

		track := database.Track{
			ReleaseID:       releaseDetails.ID,
			Position:        discTrack.Position,
			Title:           discTrack.Title,
			DurationText:    discTrack.Duration,
			DurationSeconds: durationSeconds,
		}

		tracks = append(tracks, track)
	}

	// If we don't have valid durations, estimate based on format
	if totalDurationSeconds == 0 {
		// Get release to check formats
		release, err := c.DB.GetReleaseByID(releaseDetails.ID)
		if err != nil {
			slog.Error(
				"Failed to get release for duration estimation",
				"error",
				err,
				"id",
				releaseDetails.ID,
			)
			return err
		}

		totalDurationSeconds = estimateVinylPlaytime(release)
		isDurationEstimated = true
	}
	slog.Info(
		"Successfully fetched release details",
		"releaseID",
		releaseDetails.ID,
		"trackCount",
		len(tracks),
		"totalDuration",
		totalDurationSeconds,
		"isEstimated",
		isDurationEstimated,
	)

	// Save tracks and update release duration
	// err = c.DB.SaveTracks(releaseDetails.ID, tracks, totalDurationSeconds, isDurationEstimated)
	// if err != nil {
	// 	slog.Error("Failed to save tracks", "error", err, "releaseID", releaseDetails.ID)
	// 	return err
	// }
	//
	// slog.Info("Successfully fetched and saved release details",
	// 	"releaseID", releaseDetails.ID,
	// 	"trackCount", len(tracks),
	// 	"totalDuration", totalDurationSeconds,
	// 	"isEstimated", isDurationEstimated)

	return nil
}

// Helper to convert duration string to seconds
func convertDurationToSeconds(duration string) (int, error) {
	if duration == "" {
		return 0, nil
	}

	parts := strings.Split(duration, ":")
	if len(parts) == 2 {
		// mm:ss format
		minutes, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		seconds, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		return minutes*60 + seconds, nil
	} else if len(parts) == 3 {
		// hh:mm:ss format
		hours, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		minutes, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			return 0, err
		}
		return hours*3600 + minutes*60 + seconds, nil
	}

	return 0, fmt.Errorf("unsupported duration format: %s", duration)
}

// Estimate vinyl playtime based on format information
func estimateVinylPlaytime(release *database.Release) int {
	totalSeconds := 0

	for _, format := range release.Formats {
		if format.Name != "Vinyl" {
			continue // Skip non-vinyl formats
		}

		// Extract quantity
		qty := format.Qty
		if qty == 0 {
			qty = 1 // Default to 1 if not specified
		}

		// Check vinyl specifics - look through descriptions
		is7inch := false
		is10inch := false
		is45rpm := false

		for _, desc := range format.Descriptions {
			if desc == "7\"" {
				is7inch = true
			} else if desc == "10\"" {
				is10inch = true
			} else if desc == "45 RPM" {
				is45rpm = true
			}
		}

		secondsPerSide := 0
		sidesPerDisc := 2

		// Apply estimates based on format characteristics
		if is7inch && is45rpm {
			secondsPerSide = 4 * 60 // 4 minutes per side for 7" singles
		} else if is10inch {
			secondsPerSide = 15 * 60 // 15 minutes per side for 10"
		} else if is45rpm {
			secondsPerSide = 15 * 60 // 15 minutes per side for 12" at 45rpm
		} else {
			// Default for 12" 33⅓ RPM LP
			secondsPerSide = 20 * 60 // 20 minutes per side
		}

		totalSeconds += qty * sidesPerDisc * secondsPerSide
	}

	return totalSeconds
}
