package database

import (
	"database/sql"
	"log/slog"
	"time"
)

func (s *Database) GetPlayHistory(limit, offset int) ([]PlayHistory, error) {
	query := `
		SELECT 
			ph.id, ph.release_id, ph.stylus_id, ph.played_at, ph.created_at, ph.updated_at
		FROM play_history ph
		ORDER BY ph.played_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		slog.Error("Failed to get play history", "error", err)
		return nil, err
	}
	defer rows.Close()

	var histories []PlayHistory
	for rows.Next() {
		var history PlayHistory
		var stylusID sql.NullInt64

		err := rows.Scan(
			&history.ID,
			&history.ReleaseID,
			&stylusID,
			&history.PlayedAt,
			&history.CreatedAt,
			&history.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan play history", "error", err)
			return nil, err
		}

		if stylusID.Valid {
			id := int(stylusID.Int64)
			history.StylusID = &id
		}

		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating play history rows", "error", err)
		return histories, err
	}

	return histories, nil
}

func (s *Database) GetPlayHistoryWithDetails(limit, offset int) ([]PlayHistory, error) {
	// First get the basic play history records
	histories, err := s.GetPlayHistory(limit, offset)
	if err != nil {
		return nil, err
	}

	// Early return if no histories
	if len(histories) == 0 {
		return histories, nil
	}

	// For each history, load the release and stylus details
	for i := range histories {
		// Load release
		release, err := s.GetReleaseByID(histories[i].ReleaseID)
		if err != nil {
			slog.Error(
				"Failed to get release for play history",
				"error",
				err,
				"release_id",
				histories[i].ReleaseID,
			)
			continue
		}
		histories[i].Release = *release

		// Load stylus if present
		if histories[i].StylusID != nil {
			stylus, err := s.GetStylusByID(*histories[i].StylusID)
			if err != nil {
				slog.Error(
					"Failed to get stylus for play history",
					"error",
					err,
					"stylus_id",
					*histories[i].StylusID,
				)
				continue
			}
			histories[i].Stylus = stylus
		}
	}

	return histories, nil
}

func (s *Database) GetPlayHistoryByReleaseID(
	releaseID int,
	limit, offset int,
) ([]PlayHistory, error) {
	query := `
		SELECT 
			ph.id, ph.release_id, ph.stylus_id, ph.played_at, ph.created_at, ph.updated_at
		FROM play_history ph
		WHERE ph.release_id = ?
		ORDER BY ph.played_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.DB.Query(query, releaseID, limit, offset)
	if err != nil {
		slog.Error("Failed to get play history by release", "error", err, "release_id", releaseID)
		return nil, err
	}
	defer rows.Close()

	var histories []PlayHistory
	for rows.Next() {
		var history PlayHistory
		var stylusID sql.NullInt64

		err := rows.Scan(
			&history.ID,
			&history.ReleaseID,
			&stylusID,
			&history.PlayedAt,
			&history.CreatedAt,
			&history.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan play history", "error", err)
			return nil, err
		}

		if stylusID.Valid {
			id := int(stylusID.Int64)
			history.StylusID = &id
		}

		// Load stylus if present
		if history.StylusID != nil {
			stylus, err := s.GetStylusByID(*history.StylusID)
			if err != nil {
				slog.Error(
					"Failed to get stylus for play history",
					"error",
					err,
					"stylus_id",
					*history.StylusID,
				)
			} else {
				history.Stylus = stylus
			}
		}

		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating play history rows", "error", err)
		return histories, err
	}

	return histories, nil
}

// GetReleaseByID is a helper function to get a release by ID
// This should be added to release.database.go, but we'll include it here for completeness
func (s *Database) GetReleaseByID(id int) (*Release, error) {
	query := `
		SELECT 
			r.id, r.instance_id, r.folder_id, r.rating, r.title, 
			r.year, r.resource_url, r.thumb, r.cover_image, 
			r.created_at, r.updated_at, r.last_synced
		FROM releases r
		WHERE r.id = ?
	`

	row := s.DB.QueryRow(query, id)

	var release Release
	var year sql.NullInt32

	err := row.Scan(
		&release.ID,
		&release.InstanceID,
		&release.FolderID,
		&release.Rating,
		&release.Title,
		&year,
		&release.ResourceURL,
		&release.Thumb,
		&release.CoverImage,
		&release.CreatedAt,
		&release.UpdatedAt,
		&release.LastSynced,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No release found
		}
		slog.Error("Failed to scan release", "error", err)
		return nil, err
	}

	if year.Valid {
		y := int(year.Int32)
		release.Year = &y
	}

	// Load related data like artists, labels, etc. if needed
	// This can be expanded later for more complete release information

	return &release, nil
}

// CreatePlayHistory creates a new play history record
func (s *Database) CreatePlayHistory(history *PlayHistory) error {
	query := `
		INSERT INTO play_history (
			release_id, stylus_id, played_at
		) VALUES (?, ?, ?)
		RETURNING id, created_at, updated_at
	`

	var stylusID interface{}
	if history.StylusID != nil {
		stylusID = *history.StylusID
	} else {
		stylusID = nil
	}

	err := s.DB.QueryRow(
		query,
		history.ReleaseID,
		stylusID,
		history.PlayedAt,
	).Scan(&history.ID, &history.CreatedAt, &history.UpdatedAt)
	if err != nil {
		slog.Error("Failed to create play history", "error", err)
		return err
	}

	return nil
}

// UpdatePlayHistory updates an existing play history record
func (s *Database) UpdatePlayHistory(history *PlayHistory) error {
	query := `
		UPDATE play_history SET
			release_id = ?,
			stylus_id = ?,
			played_at = ?
		WHERE id = ?
		RETURNING updated_at
	`

	var stylusID interface{}
	if history.StylusID != nil {
		stylusID = *history.StylusID
	} else {
		stylusID = nil
	}

	err := s.DB.QueryRow(
		query,
		history.ReleaseID,
		stylusID,
		history.PlayedAt,
		history.ID,
	).Scan(&history.UpdatedAt)
	if err != nil {
		slog.Error("Failed to update play history", "error", err)
		return err
	}

	return nil
}

// DeletePlayHistory deletes a play history record by ID
func (s *Database) DeletePlayHistory(id int) error {
	_, err := s.DB.Exec("DELETE FROM play_history WHERE id = ?", id)
	if err != nil {
		slog.Error("Failed to delete play history", "error", err)
		return err
	}

	return nil
}

// GetPlayCountByRelease gets the number of plays for each release
func (s *Database) GetPlayCountByRelease() (map[int]int, error) {
	query := `
		SELECT release_id, COUNT(*) as play_count
		FROM play_history
		GROUP BY release_id
		ORDER BY play_count DESC
	`

	rows, err := s.DB.Query(query)
	if err != nil {
		slog.Error("Failed to get play counts", "error", err)
		return nil, err
	}
	defer rows.Close()

	playCounts := make(map[int]int)
	for rows.Next() {
		var releaseID, count int
		if err := rows.Scan(&releaseID, &count); err != nil {
			slog.Error("Failed to scan play count", "error", err)
			return playCounts, err
		}
		playCounts[releaseID] = count
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating play count rows", "error", err)
		return playCounts, err
	}

	return playCounts, nil
}

// GetRecentPlays gets the most recent plays
func (s *Database) GetRecentPlays(limit int) ([]PlayHistory, error) {
	return s.GetPlayHistory(limit, 0)
}

// GetPlaysByTimeRange gets plays within a specific time range
func (s *Database) GetPlaysByTimeRange(start, end time.Time) ([]PlayHistory, error) {
	query := `
		SELECT 
			ph.id, ph.release_id, ph.stylus_id, ph.played_at, ph.created_at, ph.updated_at
		FROM play_history ph
		WHERE ph.played_at BETWEEN ? AND ?
		ORDER BY ph.played_at DESC
	`

	rows, err := s.DB.Query(query, start, end)
	if err != nil {
		slog.Error("Failed to get plays by time range", "error", err)
		return nil, err
	}
	defer rows.Close()

	var histories []PlayHistory
	for rows.Next() {
		var history PlayHistory
		var stylusID sql.NullInt64

		err := rows.Scan(
			&history.ID,
			&history.ReleaseID,
			&stylusID,
			&history.PlayedAt,
			&history.CreatedAt,
			&history.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan play history", "error", err)
			return nil, err
		}

		if stylusID.Valid {
			id := int(stylusID.Int64)
			history.StylusID = &id
		}

		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating play history rows", "error", err)
		return histories, err
	}

	return histories, nil
}
