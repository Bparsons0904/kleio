package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"kleio/internal/database"
	"log/slog"
	"net/http"
	"reflect"
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

	c.RateLimit.UpdateLimits(resp)

	// Check for rate limiting first
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		slog.Warn("Rate limited while fetching release details", 
			"url", resourceURL, 
			"retryAfter", retryAfter,
			"releaseID", release.ID)
		// Return error to allow caller to handle retry
		return nil, fmt.Errorf("rate limited: retry after %s seconds", retryAfter)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("API returned non-200 status for release details",
			"status", resp.StatusCode,
			"body", string(body),
			"url", resourceURL,
			"releaseID", release.ID,
			"releaseTitle", release.Title,
		)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode the response
	var releaseDetails struct {
		ID        int                     `json:"id"`
		Tracklist []database.DiscogsTrack `json:"tracklist"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&releaseDetails); err != nil {
		slog.Error("Error decoding release details response", 
			"error", err,
			"url", resourceURL,
			"releaseID", release.ID,
			"releaseTitle", release.Title)
		return nil, err
	}

	var tracks []database.Track
	slog.Info("Processing tracklist", 
		"releaseID", release.ID,
		"releaseTitle", release.Title,
		"trackCount", len(releaseDetails.Tracklist))

	for i, discTrack := range releaseDetails.Tracklist {
		// Skip non-track items (headers, etc.)
		if discTrack.Type != "track" && discTrack.Type != "" {
			slog.Debug("Skipping non-track item", 
				"releaseID", release.ID,
				"position", i,
				"type", discTrack.Type,
				"title", discTrack.Title)
			continue
		}

		durationSeconds := convertDurationToSeconds(discTrack.Duration)
		track := database.Track{
			ReleaseID:       release.ID,
			Position:        discTrack.Position,
			Title:           discTrack.Title,
			DurationText:    discTrack.Duration,
			DurationSeconds: durationSeconds,
		}

		tracks = append(tracks, track)
	}

	slog.Info("Finished processing tracklist", 
		"releaseID", release.ID,
		"validTracks", len(tracks))

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

func calculateTimeToDuration(parts []string, index int, multiplier int) int {
	if len(parts) > index {
		seconds, err := strconv.Atoi(parts[index])
		if err != nil {
			slog.Error(
				"Failed to parse seconds",
				"error",
				err,
				"seconds",
				parts[index],
				"index",
				index,
			)
			return 0
		}
		return seconds * multiplier
	}
	return 0
}

func convertDurationToSeconds(duration string) int {
	if duration == "" {
		return 0
	}

	totalSeconds := 0

	switch {
	case isNumeric(duration):
		return calculateTimeToDuration([]string{duration}, 0, 1)

	case isMinutesSeconds(duration):
		parts := strings.Split(duration, ":")
		totalSeconds += calculateTimeToDuration(parts, 0, 60)
		totalSeconds += calculateTimeToDuration(parts, 1, 1)
		return totalSeconds

	case isHoursMinutesSeconds(duration):
		parts := strings.Split(duration, ":")
		totalSeconds += calculateTimeToDuration(parts, 0, 3600)
		totalSeconds += calculateTimeToDuration(parts, 1, 60)
		totalSeconds += calculateTimeToDuration(parts, 2, 1)
		return totalSeconds

	default:
		slog.Error(
			"Failed to parse duration",
			"duration",
			duration,
			"type",
			reflect.TypeOf(duration),
		)
		return 0
	}
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isMinutesSeconds(s string) bool {
	// Match pattern like "3:45" or "12:30"
	pattern := `^\d+:\d{2}$`
	match, err := regexp.MatchString(pattern, s)
	return err == nil && match
}

func isHoursMinutesSeconds(s string) bool {
	// Match pattern like "1:23:45"
	pattern := `^\d+:\d{2}:\d{2}$`
	match, err := regexp.MatchString(pattern, s)
	return err == nil && match
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
