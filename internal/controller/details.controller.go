package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"kleio/internal/database"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (c *Controller) GetReleaseDetails(
	release database.Release,
	token string,
) ([]database.Track, error) {
	resourceURL := release.ResourceURL + "?token=" + token

	// Create a new request
	req, err := http.NewRequest("GET", resourceURL, nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	// Check for rate limiting
	if resp.StatusCode == http.StatusTooManyRequests {
		// Get retry after header if available
		retryAfter := resp.Header.Get("Retry-After")
		slog.Warn("Rate limited by Discogs API", "retryAfter", retryAfter)
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("API returned non-200 status",
			"status", resp.StatusCode,
			"body", string(body),
			"url", resourceURL)
		return nil, err
	}

	// Decode the response
	var releaseDetails struct {
		ID        int                     `json:"id"`
		Tracklist []database.DiscogsTrack `json:"tracklist"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&releaseDetails); err != nil {
		slog.Error("Error decoding response", "error", err)
		return nil, err
	}

	var tracks []database.Track
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

		track := database.Track{
			ReleaseID:       release.ID,
			Position:        discTrack.Position,
			Title:           discTrack.Title,
			DurationText:    discTrack.Duration,
			DurationSeconds: durationSeconds,
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (c *Controller) calculateTrackDurations(
	releaseID int,
	tracks []database.Track,
) (int, bool, error) {
	totalDurationSeconds := 0
	isDurationEstimated := false

	for _, track := range tracks {
		totalDurationSeconds += track.DurationSeconds
	}

	// If we don't have valid durations, estimate based on format
	if totalDurationSeconds == 0 {
		// Get release to check formats
		release, err := c.DB.GetReleaseByID(releaseID)
		if err != nil {
			slog.Error(
				"Failed to get release for duration estimation",
				"error",
				err,
				"id",
				releaseID,
			)
			return 0, false, err
		}

		totalDurationSeconds = estimateVinylPlaytime(release)
		isDurationEstimated = true
	}

	return totalDurationSeconds, isDurationEstimated, nil
}

func (c *Controller) SaveReleaseDetails(
	releaseDetails database.Release,
	tracks []database.Track,
) error {
	// Save tracks and update release duration
	err := c.DB.SaveTracks(releaseDetails.ID, tracks)
	if err != nil {
		slog.Error("Failed to save tracks", "error", err, "releaseID", releaseDetails.ID)
		return err
	}

	// Update release duration
	// err = c.DB.UpdateReleaseWithDetails(
	// 	releaseDetails.ID,
	// 	totalDurationSeconds,
	// 	isDurationEstimated,
	// )
	// if err != nil {
	// 	slog.Error(
	// 		"Failed to update release duration",
	// 		"error",
	// 		err,
	// 		"releaseID",
	// 		releaseDetails.ID,
	// 	)
	// 	return err
	// }

	return nil
}

// Helper to convert duration string to seconds
func convertDurationToSeconds(duration string) (int, error) {
	// Handle empty strings
	if duration == "" {
		return 0, nil
	}

	// Check if the duration is already in seconds format (just a number)
	if seconds, err := strconv.Atoi(duration); err == nil {
		return seconds, nil
	}

	// Handle mm:ss format
	if match, _ := regexp.MatchString(`^\d+:\d{2}$`, duration); match {
		parts := strings.Split(duration, ":")
		minutes, _ := strconv.Atoi(parts[0])
		seconds, _ := strconv.Atoi(parts[1])
		return minutes*60 + seconds, nil
	}

	// Handle hh:mm:ss format
	if match, _ := regexp.MatchString(`^\d+:\d{2}:\d{2}$`, duration); match {
		parts := strings.Split(duration, ":")
		hours, _ := strconv.Atoi(parts[0])
		minutes, _ := strconv.Atoi(parts[1])
		seconds, _ := strconv.Atoi(parts[2])
		return hours*3600 + minutes*60 + seconds, nil
	}

	// Handle additional formats as needed

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

		// Check format descriptions
		isLP := false
		is7inch := false
		is10inch := false
		is12inch := false
		is45rpm := false
		is33rpm := false
		isSingle := false
		isAlbum := false

		for _, desc := range format.Descriptions {
			switch strings.ToLower(desc) {
			case "7\"":
				is7inch = true
			case "10\"":
				is10inch = true
			case "12\"":
				is12inch = true
			case "lp":
				isLP = true
				isAlbum = true
			case "album":
				isAlbum = true
			case "single":
				isSingle = true
			case "33 rpm", "33⅓ rpm", "33 1/3 rpm", "33rpm":
				is33rpm = true
			case "45 rpm", "45rpm":
				is45rpm = true
			}
		}

		// If size not specified but is LP/Album, assume 12"
		if (!is7inch && !is10inch && !is12inch) && (isLP || isAlbum) {
			is12inch = true
		}

		// If RPM not specified
		if !is45rpm && !is33rpm {
			if is7inch {
				is45rpm = true // 7" are typically 45 RPM
			} else if is12inch || isLP || isAlbum {
				is33rpm = true // 12" LPs are typically 33⅓ RPM
			}
		}

		// Calculate seconds per side based on format
		secondsPerSide := 0
		sidesPerDisc := 2

		if is7inch && is45rpm {
			if isSingle {
				secondsPerSide = 3 * 60 // 3 minutes per side for 7" singles
			} else {
				secondsPerSide = 5 * 60 // 5 minutes per side for 7" EPs
			}
		} else if is7inch && is33rpm {
			secondsPerSide = 8 * 60 // 8 minutes per side for 7" at 33⅓
		} else if is10inch && is45rpm {
			secondsPerSide = 9 * 60 // 9 minutes per side for 10" at 45rpm
		} else if is10inch && is33rpm {
			secondsPerSide = 15 * 60 // 15 minutes per side for 10" at 33⅓
		} else if is12inch && is45rpm {
			secondsPerSide = 12 * 60 // 12 minutes per side for 12" at 45rpm
		} else if is12inch && is33rpm {
			secondsPerSide = 22 * 60 // 22 minutes per side for 12" at 33⅓
		} else {
			// Default for standard album
			secondsPerSide = 20 * 60 // 20 minutes per side
		}

		totalSeconds += qty * sidesPerDisc * secondsPerSide
	}

	// Fallback if we couldn't determine a valid estimate
	if totalSeconds == 0 {
		totalSeconds = 40 * 60 // Default 40 minutes for an album
	}

	return totalSeconds
}

// func estimateVinylPlaytime(release *database.Release) int {
// 	totalSeconds := 0
//
// 	for _, format := range release.Formats {
// 		if format.Name != "Vinyl" {
// 			continue // Skip non-vinyl formats
// 		}
//
// 		// Extract quantity
// 		qty := format.Qty
// 		if qty == 0 {
// 			qty = 1 // Default to 1 if not specified
// 		}
//
// 		// Check vinyl specifics - look through descriptions
// 		is7inch := false
// 		is10inch := false
// 		is45rpm := false
//
// 		for _, desc := range format.Descriptions {
// 			if desc == "7\"" {
// 				is7inch = true
// 			} else if desc == "10\"" {
// 				is10inch = true
// 			} else if desc == "45 RPM" {
// 				is45rpm = true
// 			}
// 		}
//
// 		secondsPerSide := 0
// 		sidesPerDisc := 2
//
// 		// Apply estimates based on format characteristics
// 		if is7inch && is45rpm {
// 			secondsPerSide = 4 * 60 // 4 minutes per side for 7" singles
// 		} else if is10inch {
// 			secondsPerSide = 15 * 60 // 15 minutes per side for 10"
// 		} else if is45rpm {
// 			secondsPerSide = 15 * 60 // 15 minutes per side for 12" at 45rpm
// 		} else {
// 			// Default for 12" 33⅓ RPM LP
// 			secondsPerSide = 20 * 60 // 20 minutes per side
// 		}
//
// 		totalSeconds += qty * sidesPerDisc * secondsPerSide
// 	}
//
// 	return totalSeconds
// }
