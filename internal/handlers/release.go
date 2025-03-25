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

func GetAllReleases(db *sql.DB) ([]Release, error) {
	// Use the query that fetches all releases with JSON-aggregated related data
	query := `
	SELECT 
		r.id,
		r.instance_id,
		r.folder_id,
		r.rating,
		r.title,
		r.year,
		r.resource_url,
		r.thumb,
		r.cover_image,
		r.created_at,
		r.updated_at,
		r.last_synced,
		
		-- Artists (JSON array)
		(
			SELECT json_group_array(
				json_object(
					'artist_id', a.id,
					'name', a.name,
					'resource_url', a.resource_url,
					'join_relation', ra.join_relation,
					'anv', ra.anv,
					'tracks', ra.tracks,
					'role', ra.role
				)
			)
			FROM release_artists ra
			JOIN artists a ON ra.artist_id = a.id
			WHERE ra.release_id = r.id
		) AS artists_json,
		
		-- Labels (JSON array)
		(
			SELECT json_group_array(
				json_object(
					'label_id', l.id,
					'name', l.name,
					'resource_url', l.resource_url,
					'entity_type', l.entity_type,
					'catno', rl.catno
				)
			)
			FROM release_labels rl
			JOIN labels l ON rl.label_id = l.id
			WHERE rl.release_id = r.id
		) AS labels_json,
		
		-- Formats with descriptions (JSON array)
		(
			SELECT json_group_array(
				json_object(
					'id', f.id,
					'release_id', f.release_id,
					'name', f.name,
					'qty', f.qty,
					'descriptions', (
						SELECT json_group_array(fd.description)
						FROM format_descriptions fd
						WHERE fd.format_id = f.id
					)
				)
			)
			FROM formats f
			WHERE f.release_id = r.id
		) AS formats_json,
		
		-- Genres (JSON array)
		(
			SELECT json_group_array(
				json_object(
					'id', g.id,
					'name', g.name
				)
			)
			FROM release_genres rg
			JOIN genres g ON rg.genre_id = g.id
			WHERE rg.release_id = r.id
		) AS genres_json,
		
		-- Styles (JSON array)
		(
			SELECT json_group_array(
				json_object(
					'id', s.id,
					'name', s.name
				)
			)
			FROM release_styles rs
			JOIN styles s ON rs.style_id = s.id
			WHERE rs.release_id = r.id
		) AS styles_json,
		
		-- Notes (JSON array)
		(
			SELECT json_group_array(
				json_object(
					'release_id', rn.release_id,
					'field_id', rn.field_id,
					'value', rn.value
				)
			)
			FROM release_notes rn
			WHERE rn.release_id = r.id
		) AS notes_json
	FROM releases r
	ORDER BY r.title`

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying releases: %w", err)
	}
	defer rows.Close()

	// Prepare slice to hold all releases
	var releases []Release

	// Process each row
	for rows.Next() {
		var release Release
		var artistsJSON, labelsJSON, formatsJSON, genresJSON, stylesJSON, notesJSON []byte

		// Scan the row into our variables
		err := rows.Scan(
			&release.ID,
			&release.InstanceID,
			&release.FolderID,
			&release.Rating,
			&release.Title,
			&release.Year,
			&release.ResourceURL,
			&release.Thumb,
			&release.CoverImage,
			&release.CreatedAt,
			&release.UpdatedAt,
			&release.LastSynced,
			&artistsJSON,
			&labelsJSON,
			&formatsJSON,
			&genresJSON,
			&stylesJSON,
			&notesJSON,
		)
		if err != nil {
			slog.Error("Error scanning release row", "error", err)
			continue // Skip this release but continue with others
		}

		// Unmarshal the JSON arrays into our struct fields
		// Artists
		var artistsData []struct {
			ArtistID     int    `json:"artist_id"`
			Name         string `json:"name"`
			ResourceURL  string `json:"resource_url"`
			JoinRelation string `json:"join_relation"`
			ANV          string `json:"anv"`
			Tracks       string `json:"tracks"`
			Role         string `json:"role"`
		}
		if err := json.Unmarshal(artistsJSON, &artistsData); err == nil {
			for _, a := range artistsData {
				artist := &Artist{
					ID:          a.ArtistID,
					Name:        a.Name,
					ResourceURL: a.ResourceURL,
				}
				releaseArtist := ReleaseArtist{
					ReleaseID:    release.ID,
					ArtistID:     a.ArtistID,
					JoinRelation: a.JoinRelation,
					ANV:          a.ANV,
					Tracks:       a.Tracks,
					Role:         a.Role,
					Artist:       artist,
				}
				release.Artists = append(release.Artists, releaseArtist)
			}
		}

		// Labels
		var labelsData []struct {
			LabelID     int    `json:"label_id"`
			Name        string `json:"name"`
			ResourceURL string `json:"resource_url"`
			EntityType  string `json:"entity_type"`
			CatNo       string `json:"catno"`
		}
		if err := json.Unmarshal(labelsJSON, &labelsData); err == nil {
			for _, l := range labelsData {
				label := &Label{
					ID:          l.LabelID,
					Name:        l.Name,
					ResourceURL: l.ResourceURL,
					EntityType:  l.EntityType,
				}
				releaseLabel := ReleaseLabel{
					ReleaseID: release.ID,
					LabelID:   l.LabelID,
					CatNo:     l.CatNo,
					Label:     label,
				}
				release.Labels = append(release.Labels, releaseLabel)
			}
		}

		// Formats with descriptions
		var formatsData []Format
		if err := json.Unmarshal(formatsJSON, &formatsData); err == nil {
			for _, f := range formatsData {
				format := Format{
					ID:           f.ID,
					ReleaseID:    f.ReleaseID,
					Name:         f.Name,
					Qty:          f.Qty,
					Descriptions: f.Descriptions,
				}
				release.Formats = append(release.Formats, format)
			}
		}

		// Genres
		var genresData []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		if err := json.Unmarshal(genresJSON, &genresData); err == nil {
			for _, g := range genresData {
				genre := Genre{
					ID:   g.ID,
					Name: g.Name,
				}
				release.Genres = append(release.Genres, genre)
			}
		}

		// Styles
		var stylesData []Style
		if err := json.Unmarshal(stylesJSON, &stylesData); err == nil {
			for _, s := range stylesData {
				style := Style{
					ID:   s.ID,
					Name: s.Name,
				}
				release.Styles = append(release.Styles, style)
			}
		}

		// Notes
		var notesData []ReleaseNote
		if err := json.Unmarshal(notesJSON, &notesData); err == nil {
			for _, n := range notesData {
				note := ReleaseNote{
					ReleaseID: n.ReleaseID,
					FieldID:   n.FieldID,
					Value:     n.Value,
				}
				release.Notes = append(release.Notes, note)
			}
		}

		releases = append(releases, release)
	}

	if err := rows.Err(); err != nil {
		return releases, fmt.Errorf("error iterating releases: %w", err)
	}

	slog.Info("Retrieved all releases", "count", len(releases))
	return releases, nil
}

// GetAllReleasesAsJSON retrieves all releases with their complete related data as a JSON string
func GetAllReleasesAsJSON(db *sql.DB) (string, error) {
	releases, err := GetAllReleases(db)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(releases, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling releases to JSON: %w", err)
	}

	return string(jsonData), nil
}
