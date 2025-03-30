package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"time"
)

func (s *Database) GetAllReleases() ([]Release, error) {
	// Query to get all releases with their related data as JSON
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
		) AS artists,
		
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
		) AS labels,
		
		-- Formats with descriptions (JSON array)
		(
			SELECT json_group_array(
				json_object(
					'format_id', f.id,
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
		) AS formats,
		
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
		) AS genres,
		
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
		) AS styles,
		
		-- Notes (JSON array)
		(
			SELECT json_group_array(
				json_object(
					'field_id', rn.field_id,
					'value', rn.value
				)
			)
			FROM release_notes rn
			WHERE rn.release_id = r.id
		) AS notes,
		
		-- Play History (JSON array)
		(
			SELECT json_group_array(
				json_object(
					'id', ph.id,
					'release_id', ph.release_id,
					'stylus_id', ph.stylus_id,
					'played_at', ph.played_at,
					'created_at', ph.created_at,
					'updated_at', ph.updated_at,
					'notes', ph.notes,
					'stylus', CASE WHEN ph.stylus_id IS NOT NULL THEN (
						SELECT json_object(
							'id', s.id,
							'name', s.name,
							'manufacturer', s.manufacturer,
							'expected_lifespan_hours', s.expected_lifespan_hours,
							'purchase_date', s.purchase_date,
							'active', s.active,
							'primary_stylus', s.primary_stylus,
							'created_at', s.created_at,
							'updated_at', s.updated_at,
							'owned', s.owned,
							'base_model', s.base_model
						)
						FROM styluses s
						WHERE s.id = ph.stylus_id
					) ELSE NULL END
				)
			)
			FROM play_history ph
			WHERE ph.release_id = r.id
			ORDER BY ph.played_at DESC
		) AS play_history,
		
		-- Cleaning History (JSON array)
		(
			SELECT json_group_array(
				json_object(
					'id', ch.id,
					'release_id', ch.release_id,
					'cleaned_at', ch.cleaned_at,
					'notes', ch.notes,
					'created_at', ch.created_at,
					'updated_at', ch.updated_at
				)
			)
			FROM cleaning_history ch
			WHERE ch.release_id = r.id
			ORDER BY ch.cleaned_at DESC
		) AS cleaning_history
	FROM releases r
	ORDER BY r.title`

	// Execute the query
	rows, err := s.DB.Query(query)
	if err != nil {
		slog.Error("Error querying releases", "error", err)
		return nil, err
	}
	defer rows.Close()

	var releases []Release

	// Process each row
	for rows.Next() {
		var release Release
		var artistsJSON, labelsJSON, formatsJSON, genresJSON, stylesJSON, notesJSON, playHistoryJSON, cleaningHistoryJSON []byte

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
			&artistsJSON,
			&labelsJSON,
			&formatsJSON,
			&genresJSON,
			&stylesJSON,
			&notesJSON,
			&playHistoryJSON,
			&cleaningHistoryJSON,
		)
		if err != nil {
			slog.Error("Error scanning release row", "error", err)
			continue // Skip this release but continue with others
		}

		// Unmarshal the JSON data for artists
		var artistsData []ArtistData
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

		// Unmarshal the JSON data for labels
		var labelsData []LabelData
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

		// Unmarshal the JSON data for formats
		var formatsData []FormatData
		if err := json.Unmarshal(formatsJSON, &formatsData); err == nil {
			for _, f := range formatsData {
				format := Format{
					ID:           f.FormatID,
					ReleaseID:    release.ID,
					Name:         f.Name,
					Qty:          f.Qty,
					Descriptions: f.Descriptions,
				}
				release.Formats = append(release.Formats, format)
			}
		}

		// Unmarshal the JSON data for genres
		var genresData []Genre
		if err := json.Unmarshal(genresJSON, &genresData); err == nil {
			for _, g := range genresData {
				genre := Genre{
					ID:   g.ID,
					Name: g.Name,
				}
				release.Genres = append(release.Genres, genre)
			}
		}

		// Unmarshal the JSON data for styles
		var stylesData []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		if err := json.Unmarshal(stylesJSON, &stylesData); err == nil {
			for _, s := range stylesData {
				style := Style{
					ID:   s.ID,
					Name: s.Name,
				}
				release.Styles = append(release.Styles, style)
			}
		}

		// Unmarshal the JSON data for notes
		var notesData []NoteData
		if err := json.Unmarshal(notesJSON, &notesData); err == nil {
			for _, n := range notesData {
				note := ReleaseNote{
					ReleaseID: release.ID,
					FieldID:   n.FieldID,
					Value:     n.Value,
				}
				release.Notes = append(release.Notes, note)
			}
		}

		// Unmarshal the JSON data for play history
		var playHistoryData []struct {
			ID        int             `json:"id"`
			ReleaseID int             `json:"release_id"`
			StylusID  *int            `json:"stylus_id"`
			PlayedAt  string          `json:"played_at"`
			CreatedAt string          `json:"created_at"`
			UpdatedAt string          `json:"updated_at"`
			Stylus    json.RawMessage `json:"stylus"`
			Notes     string          `json:"notes"`
		}

		if err := json.Unmarshal(playHistoryJSON, &playHistoryData); err == nil {
			for _, ph := range playHistoryData {
				playHistory := PlayHistory{
					ID:        ph.ID,
					ReleaseID: ph.ReleaseID,
					StylusID:  ph.StylusID,
					PlayedAt:  parseTime(ph.PlayedAt),
					CreatedAt: parseTime(ph.CreatedAt),
					UpdatedAt: parseTime(ph.UpdatedAt),
					Notes:     ph.Notes,
				}

				// If stylus data is present, unmarshal it
				if len(ph.Stylus) > 2 { // Check if not null or empty {}
					// Temporary struct with string fields for dates and int for bools
					var tempStylus struct {
						ID               int     `json:"id"`
						Name             string  `json:"name"`
						Manufacturer     string  `json:"manufacturer"`
						ExpectedLifespan int     `json:"expected_lifespan_hours"`
						PurchaseDate     *string `json:"purchase_date"`
						Active           int     `json:"active"`
						Primary          int     `json:"primary_stylus"`
						CreatedAt        string  `json:"created_at"`
						UpdatedAt        string  `json:"updated_at"`
						Owned            int     `json:"owned"`
						BaseModel        int     `json:"base_model"`
					}

					if err := json.Unmarshal(ph.Stylus, &tempStylus); err == nil {
						// Parse time fields
						var createdAt, updatedAt time.Time
						var purchaseDate *time.Time

						// Parse created_at and updated_at
						if tempStylus.CreatedAt != "" {
							t, err := time.Parse("2006-01-02 15:04:05", tempStylus.CreatedAt)
							if err != nil {
								slog.Error(
									"Failed to parse created_at",
									"error",
									err,
									"value",
									tempStylus.CreatedAt,
								)
							} else {
								createdAt = t
							}
						}

						if tempStylus.UpdatedAt != "" {
							t, err := time.Parse("2006-01-02 15:04:05", tempStylus.UpdatedAt)
							if err != nil {
								slog.Error(
									"Failed to parse updated_at",
									"error",
									err,
									"value",
									tempStylus.UpdatedAt,
								)
							} else {
								updatedAt = t
							}
						}

						// Parse purchase_date (which might be null)
						if tempStylus.PurchaseDate != nil && *tempStylus.PurchaseDate != "" {
							t, err := time.Parse("2006-01-02 15:04:05", *tempStylus.PurchaseDate)
							if err != nil {
								slog.Error(
									"Failed to parse purchase_date",
									"error",
									err,
									"value",
									*tempStylus.PurchaseDate,
								)
							} else {
								purchaseDate = &t
							}
						}

						// Convert to actual Stylus struct
						stylusData := Stylus{
							ID:               tempStylus.ID,
							Name:             tempStylus.Name,
							Manufacturer:     tempStylus.Manufacturer,
							ExpectedLifespan: tempStylus.ExpectedLifespan,
							PurchaseDate:     purchaseDate,
							Active:           tempStylus.Active != 0,
							Primary:          tempStylus.Primary != 0,
							CreatedAt:        createdAt,
							UpdatedAt:        updatedAt,
							Owned:            tempStylus.Owned != 0,
							BaseModel:        tempStylus.BaseModel != 0,
						}
						playHistory.Stylus = &stylusData
					} else {
						slog.Error("Failed to unmarshal stylus data", "error", err, "data", string(ph.Stylus))
					}
				}

				// Add to release's play history
				release.PlayHistory = append(release.PlayHistory, playHistory)
			}
		}

		// Sort play history by played_at (most recent first)
		sort.Slice(release.PlayHistory, func(i, j int) bool {
			return release.PlayHistory[i].PlayedAt.After(release.PlayHistory[j].PlayedAt)
		})

		// Unmarshal the JSON data for cleaning history
		var cleaningHistoryData []struct {
			ID        int    `json:"id"`
			ReleaseID int    `json:"release_id"`
			CleanedAt string `json:"cleaned_at"`
			Notes     string `json:"notes"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		}

		if err := json.Unmarshal(cleaningHistoryJSON, &cleaningHistoryData); err == nil {
			for _, ch := range cleaningHistoryData {
				cleaningHistory := CleaningHistory{
					ID:        ch.ID,
					ReleaseID: ch.ReleaseID,
					CleanedAt: parseTime(ch.CleanedAt),
					Notes:     ch.Notes,
					CreatedAt: parseTime(ch.CreatedAt),
					UpdatedAt: parseTime(ch.UpdatedAt),
				}

				// Add to release's cleaning history
				release.CleaningHistory = append(release.CleaningHistory, cleaningHistory)
			}
		}

		// Sort cleaning history by cleaned_at (most recent first)
		sort.Slice(release.CleaningHistory, func(i, j int) bool {
			return release.CleaningHistory[i].CleanedAt.After(release.CleaningHistory[j].CleanedAt)
		})

		releases = append(releases, release)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating over release rows", "error", err)
		return releases, err
	}

	slog.Info("Successfully fetched all releases", "count", len(releases))
	return releases, nil
}

// Helper function to parse time strings to time.Time
func parseTime(timeStr string) time.Time {
	if timeStr == "" {
		return time.Time{}
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		// Try alternative formats if RFC3339 fails
		t, err = time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			return time.Time{}
		}
	}

	return t
}

func (s *Database) SaveReleases(response DiscogsResponse) error {
	tx, err := dbInstance.DB.Begin()
	if err != nil {
		slog.Error("Failed to begin transaction", "error", err)
		return err
	}

	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			slog.Error("Rolling back transaction due to error", "error", err)
			if err = tx.Rollback(); err != nil {
				slog.Error("Failed to rollback transaction", "error", err)
				return
			}
		}
	}()

	// Process each release
	for _, release := range response.Releases {
		// Map the release to our database schema
		// 1. Insert/Update the main release
		err = saveRelease(tx, release)
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
func saveRelease(tx *sql.Tx, release DiscogsRelease) error {
	// Prepare statement for release upsert
	stmt, err := tx.Prepare(`
		INSERT INTO releases (
			id, instance_id, folder_id, rating, title, year, 
			resource_url, thumb, cover_image
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			instance_id = excluded.instance_id,
			folder_id = excluded.folder_id,
			rating = excluded.rating,
			title = excluded.title,
			year = excluded.year,
			resource_url = excluded.resource_url,
			thumb = excluded.thumb,
			cover_image = excluded.cover_image,
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

func (s *Database) GetReleasesWithoutDuration() ([]Release, error) {
	query := `
        SELECT 
            r.id, r.resource_url,
            (
                SELECT json_group_array(
                    json_object(
                        'format_id', f.id,
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
            ) AS formats
        FROM releases r
        WHERE r.play_duration IS NULL
        ORDER BY RANDOM() -- Randomize to distribute across collection
        
    `

	rows, err := s.DB.Query(query)
	if err != nil {
		slog.Error("Failed to query releases without duration", "error", err)
		return nil, err
	}
	defer rows.Close()

	var releases []Release
	for rows.Next() {
		var release Release
		var formatsJSON []byte

		err := rows.Scan(&release.ID, &release.ResourceURL, &formatsJSON)
		if err != nil {
			slog.Error("Failed to scan release", "error", err)
			continue
		}

		// Parse formats to help with estimation
		var formatsData []FormatData
		if err := json.Unmarshal(formatsJSON, &formatsData); err == nil {
			for _, f := range formatsData {
				format := Format{
					ID:           f.FormatID,
					ReleaseID:    release.ID,
					Name:         f.Name,
					Qty:          f.Qty,
					Descriptions: f.Descriptions,
				}
				release.Formats = append(release.Formats, format)
			}
		}

		releases = append(releases, release)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating release rows", "error", err)
		return releases, err
	}

	slog.Info("Found releases without duration information", "count", len(releases))
	return releases, nil
}
