package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"kleio/internal/database"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

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
	folders, err := GetFolders(user)
	if err != nil {
		return nil, fmt.Errorf("error fetching folders: %w", err)
	}

	// Map to store all responses by folder ID
	allReleases := make(map[int][]DiscogsResponse)

	// Fetch releases for each folder
	for _, folder := range folders {
		slog.Info("Fetching releases for folder", "folderID", folder.ID, "name", folder.Name)

		// Skip folder 0 (All) since it would duplicate releases
		if folder.ID == 0 {
			slog.Info("Skipping 'All' folder to avoid duplicates", "folderID", folder.ID)
			continue
		}

	}

	return allReleases, nil
}

func SaveReleases(db *sql.DB, response DiscogsResponse) error {
	// Begin a transaction for all insert operations
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			slog.Error("Rolling back transaction due to error", "error", err)
			tx.Rollback()
		}
	}()

	// Set the current time for all last_synced values
	now := time.Now().Format(time.RFC3339)

	// Process each release
	for _, release := range response.Releases {
		// Map the release to our database schema
		// 1. Insert/Update the main release
		err = saveRelease(tx, release, now)
		if err != nil {
			return err
		}

		// 2. Process labels
		err = saveLabels(tx, release)
		if err != nil {
			return err
		}

		// 3. Process artists
		err = saveArtists(tx, release)
		if err != nil {
			return err
		}

		// 4. Process formats
		err = saveFormats(tx, release)
		if err != nil {
			return err
		}

		// 5. Process genres
		err = saveGenres(tx, release)
		if err != nil {
			return err
		}

		// 6. Process styles
		err = saveStyles(tx, release)
		if err != nil {
			return err
		}

		// 7. Process notes
		err = saveNotes(tx, release)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	slog.Info("Successfully saved releases to database", "count", len(response.Releases))
	return nil
}

// saveRelease handles inserting or updating a release in the database
func saveRelease(tx *sql.Tx, release DiscogsRelease, now string,
) error {
	// Prepare statement for release upsert
	stmt, err := tx.Prepare(`
		INSERT INTO releases (
			id, instance_id, folder_id, rating, title, year, 
			resource_url, thumb, cover_image, last_synced
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			instance_id = excluded.instance_id,
			folder_id = excluded.folder_id,
			rating = excluded.rating,
			title = excluded.title,
			year = excluded.year,
			resource_url = excluded.resource_url,
			thumb = excluded.thumb,
			cover_image = excluded.cover_image,
			last_synced = excluded.last_synced,
			updated_at = CURRENT_TIMESTAMP
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare release statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(
		release.ID,
		release.InstanceID,
		release.FolderID,
		release.Rating,
		release.BasicInfo.Title,
		release.BasicInfo.Year,
		release.BasicInfo.ResourceURL,
		release.BasicInfo.Thumb,
		release.BasicInfo.CoverImage,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to execute release statement: %w", err)
	}

	return nil
}

// saveLabels handles inserting or updating labels and their relationship to releases
func saveLabels(tx *sql.Tx, release DiscogsRelease,
) error {
	// Prepare statement for label upsert
	labelStmt, err := tx.Prepare(`
		INSERT INTO labels (id, name, resource_url, entity_type)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			resource_url = excluded.resource_url,
			entity_type = excluded.entity_type
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare label statement: %w", err)
	}
	defer labelStmt.Close()

	// Prepare statement for release_labels upsert
	relLabelStmt, err := tx.Prepare(`
		INSERT INTO release_labels (release_id, label_id, catno)
		VALUES (?, ?, ?)
		ON CONFLICT(release_id, label_id, catno) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare release_labels statement: %w", err)
	}
	defer relLabelStmt.Close()

	// Process each label
	for _, label := range release.BasicInfo.Labels {
		// Insert or update the label
		_, err = labelStmt.Exec(
			label.ID,
			label.Name,
			label.ResourceURL,
			label.EntityType,
		)
		if err != nil {
			return fmt.Errorf("failed to execute label statement: %w", err)
		}

		// Insert the release-label relationship
		_, err = relLabelStmt.Exec(
			release.ID,
			label.ID,
			label.CatNo,
		)
		if err != nil {
			return fmt.Errorf("failed to execute release_labels statement: %w", err)
		}
	}

	return nil
}

// saveArtists handles inserting or updating artists and their relationship to releases
func saveArtists(tx *sql.Tx, release DiscogsRelease,
) error {
	// Prepare statement for artist upsert
	artistStmt, err := tx.Prepare(`
		INSERT INTO artists (id, name, resource_url)
		VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			resource_url = excluded.resource_url
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare artist statement: %w", err)
	}
	defer artistStmt.Close()

	// Prepare statement for release_artists upsert
	relArtistStmt, err := tx.Prepare(`
		INSERT INTO release_artists (release_id, artist_id, join_relation, anv, tracks, role)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(release_id, artist_id, role) DO UPDATE SET
			join_relation = excluded.join_relation,
			anv = excluded.anv,
			tracks = excluded.tracks
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare release_artists statement: %w", err)
	}
	defer relArtistStmt.Close()

	// Process each artist
	for _, artist := range release.BasicInfo.Artists {
		// Insert or update the artist
		_, err = artistStmt.Exec(
			artist.ID,
			artist.Name,
			artist.ResourceURL,
		)
		if err != nil {
			return fmt.Errorf("failed to execute artist statement: %w", err)
		}

		// Insert the release-artist relationship
		_, err = relArtistStmt.Exec(
			release.ID,
			artist.ID,
			artist.Join,
			artist.ANV,
			artist.Tracks,
			artist.Role,
		)
		if err != nil {
			return fmt.Errorf("failed to execute release_artists statement: %w", err)
		}
	}

	return nil
}

// saveFormats handles inserting or updating formats and their descriptions
func saveFormats(tx *sql.Tx, release DiscogsRelease,
) error {
	// First, delete any existing formats for this release to avoid duplicates
	_, err := tx.Exec("DELETE FROM formats WHERE release_id = ?", release.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing formats: %w", err)
	}

	// Prepare statement for format insert
	formatStmt, err := tx.Prepare(`
		INSERT INTO formats (release_id, name, qty)
		VALUES (?, ?, ?)
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare format statement: %w", err)
	}
	defer formatStmt.Close()

	// Prepare statement for format descriptions
	descStmt, err := tx.Prepare(`
		INSERT INTO format_descriptions (format_id, description)
		VALUES (?, ?)
		ON CONFLICT(format_id, description) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare format_descriptions statement: %w", err)
	}
	defer descStmt.Close()

	// Process each format
	for _, format := range release.BasicInfo.Formats {
		// Convert qty from string to int
		qty, err := strconv.Atoi(format.Qty)
		if err != nil {
			// Default to 1 if conversion fails
			qty = 1
			slog.Warn(
				"Failed to convert format quantity",
				"release_id",
				release.ID,
				"qty_string",
				format.Qty,
			)
		}

		// Insert the format
		var formatID int
		err = formatStmt.QueryRow(
			release.ID,
			format.Name,
			qty,
		).Scan(&formatID)
		if err != nil {
			return fmt.Errorf("failed to execute format statement: %w", err)
		}

		// Insert each description
		for _, desc := range format.Descriptions {
			_, err = descStmt.Exec(formatID, desc)
			if err != nil {
				return fmt.Errorf("failed to execute format_descriptions statement: %w", err)
			}
		}
	}

	return nil
}

// saveGenres handles inserting or updating genres and their relationship to releases
func saveGenres(tx *sql.Tx, release DiscogsRelease,
) error {
	// First, delete any existing genre relationships for this release
	_, err := tx.Exec("DELETE FROM release_genres WHERE release_id = ?", release.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing release_genres: %w", err)
	}

	// Prepare statement for genre upsert
	genreStmt, err := tx.Prepare(`
		INSERT INTO genres (name)
		VALUES (?)
		ON CONFLICT(name) DO NOTHING
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare genre statement: %w", err)
	}
	defer genreStmt.Close()

	// Prepare statement to get genre ID if it already exists
	getGenreStmt, err := tx.Prepare(`
		SELECT id FROM genres WHERE name = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare get_genre statement: %w", err)
	}
	defer getGenreStmt.Close()

	// Prepare statement for release_genres insert
	relGenreStmt, err := tx.Prepare(`
		INSERT INTO release_genres (release_id, genre_id)
		VALUES (?, ?)
		ON CONFLICT(release_id, genre_id) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare release_genres statement: %w", err)
	}
	defer relGenreStmt.Close()

	// Process each genre
	for _, genreName := range release.BasicInfo.Genres {
		var genreID int

		// Try to insert the genre
		err = genreStmt.QueryRow(genreName).Scan(&genreID)
		if err != nil {
			// If insert failed, try to get the genre ID
			err = getGenreStmt.QueryRow(genreName).Scan(&genreID)
			if err != nil {
				return fmt.Errorf("failed to get genre ID: %w", err)
			}
		}

		// Insert the release-genre relationship
		_, err = relGenreStmt.Exec(release.ID, genreID)
		if err != nil {
			return fmt.Errorf("failed to execute release_genres statement: %w", err)
		}
	}

	return nil
}

// saveStyles handles inserting or updating styles and their relationship to releases
func saveStyles(tx *sql.Tx, release DiscogsRelease,
) error {
	// First, delete any existing style relationships for this release
	_, err := tx.Exec("DELETE FROM release_styles WHERE release_id = ?", release.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing release_styles: %w", err)
	}

	// Prepare statement for style upsert
	styleStmt, err := tx.Prepare(`
		INSERT INTO styles (name)
		VALUES (?)
		ON CONFLICT(name) DO NOTHING
		RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare style statement: %w", err)
	}
	defer styleStmt.Close()

	// Prepare statement to get style ID if it already exists
	getStyleStmt, err := tx.Prepare(`
		SELECT id FROM styles WHERE name = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare get_style statement: %w", err)
	}
	defer getStyleStmt.Close()

	// Prepare statement for release_styles insert
	relStyleStmt, err := tx.Prepare(`
		INSERT INTO release_styles (release_id, style_id)
		VALUES (?, ?)
		ON CONFLICT(release_id, style_id) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare release_styles statement: %w", err)
	}
	defer relStyleStmt.Close()

	// Process each style
	for _, styleName := range release.BasicInfo.Styles {
		var styleID int

		// Try to insert the style
		err = styleStmt.QueryRow(styleName).Scan(&styleID)
		if err != nil {
			// If insert failed, try to get the style ID
			err = getStyleStmt.QueryRow(styleName).Scan(&styleID)
			if err != nil {
				return fmt.Errorf("failed to get style ID: %w", err)
			}
		}

		// Insert the release-style relationship
		_, err = relStyleStmt.Exec(release.ID, styleID)
		if err != nil {
			return fmt.Errorf("failed to execute release_styles statement: %w", err)
		}
	}

	return nil
}

// saveNotes handles inserting or updating release notes
func saveNotes(tx *sql.Tx, release DiscogsRelease,
) error {
	// First, delete any existing notes for this release
	_, err := tx.Exec("DELETE FROM release_notes WHERE release_id = ?", release.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing release_notes: %w", err)
	}

	// Skip if no notes
	if len(release.Notes) == 0 {
		return nil
	}

	// Prepare statement for notes insert
	noteStmt, err := tx.Prepare(`
		INSERT INTO release_notes (release_id, field_id, value)
		VALUES (?, ?, ?)
		ON CONFLICT(release_id, field_id) DO UPDATE SET
			value = excluded.value
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare release_notes statement: %w", err)
	}
	defer noteStmt.Close()

	// Process each note
	for _, note := range release.Notes {
		_, err = noteStmt.Exec(
			release.ID,
			note.FieldID,
			note.Value,
		)
		if err != nil {
			return fmt.Errorf("failed to execute release_notes statement: %w", err)
		}
	}

	return nil
}
