package database

import (
	"log/slog"
	"time"
)

type CleaningHistory struct {
	ID        int       `json:"id"        db:"id"`
	ReleaseID int       `json:"releaseId" db:"release_id"`
	CleanedAt time.Time `json:"cleanedAt" db:"cleaned_at"`
	Notes     string    `json:"notes"     db:"notes"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	Release   Release   `json:"release"   db:"-"`
}

func (s *Database) CreateCleaningHistory(history *CleaningHistory) error {
	query := `
		INSERT INTO cleaning_history (
			release_id, cleaned_at, notes
		) VALUES (?, ?, ?)
		RETURNING id, created_at, updated_at
	`

	err := s.DB.QueryRow(
		query,
		history.ReleaseID,
		history.CleanedAt.Format("2006-01-02 15:04:05"),
		history.Notes,
	).Scan(&history.ID, &history.CreatedAt, &history.UpdatedAt)
	if err != nil {
		slog.Error("Failed to create cleaning history", "error", err)
		return err
	}

	return nil
}

func (s *Database) UpdateCleaningHistory(history *CleaningHistory) error {
	query := `
		UPDATE cleaning_history SET
			release_id = ?,
			cleaned_at = ?,
			notes = ?
		WHERE id = ?
		RETURNING updated_at
	`

	err := s.DB.QueryRow(
		query,
		history.ReleaseID,
		history.CleanedAt.Format("2006-01-02 15:04:05"),
		history.Notes,
		history.ID,
	).Scan(&history.UpdatedAt)
	if err != nil {
		slog.Error("Failed to update cleaning history", "error", err)
		return err
	}

	return nil
}

func (s *Database) DeleteCleaningHistory(id int) error {
	_, err := s.DB.Exec("DELETE FROM cleaning_history WHERE id = ?", id)
	if err != nil {
		slog.Error("Failed to delete cleaning history", "error", err)
		return err
	}

	return nil
}

// GetCleaningsByTimeRange gets cleanings within a specific time range
func (s *Database) GetCleaningsByTimeRange(start, end time.Time) ([]CleaningHistory, error) {
	query := `
		SELECT 
			ch.id, ch.release_id, ch.cleaned_at, ch.notes, ch.created_at, ch.updated_at
		FROM cleaning_history ch
		WHERE ch.cleaned_at BETWEEN ? AND ?
		ORDER BY ch.cleaned_at DESC
	`

	rows, err := s.DB.Query(query, start, end)
	if err != nil {
		slog.Error("Failed to get cleanings by time range", "error", err)
		return nil, err
	}
	defer rows.Close()

	var histories []CleaningHistory
	for rows.Next() {
		var history CleaningHistory

		err := rows.Scan(
			&history.ID,
			&history.ReleaseID,
			&history.CleanedAt,
			&history.Notes,
			&history.CreatedAt,
			&history.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan cleaning history", "error", err)
			return nil, err
		}

		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating cleaning history rows", "error", err)
		return histories, err
	}

	return histories, nil
}

// CountCleaningsByRelease counts the number of cleanings for each release
func (s *Database) CountCleaningsByRelease() (map[int]int, error) {
	query := `
		SELECT release_id, COUNT(*) as cleaning_count
		FROM cleaning_history
		GROUP BY release_id
		ORDER BY cleaning_count DESC
	`

	rows, err := s.DB.Query(query)
	if err != nil {
		slog.Error("Failed to get cleaning counts", "error", err)
		return nil, err
	}
	defer rows.Close()

	cleaningCounts := make(map[int]int)
	for rows.Next() {
		var releaseID, count int
		if err := rows.Scan(&releaseID, &count); err != nil {
			slog.Error("Failed to scan cleaning count", "error", err)
			return cleaningCounts, err
		}
		cleaningCounts[releaseID] = count
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating cleaning count rows", "error", err)
		return cleaningCounts, err
	}

	return cleaningCounts, nil
}

func (s *Database) GetAllCleaningHistory() ([]CleaningHistory, error) {
	query := `
		SELECT 
			ch.id, ch.release_id, ch.cleaned_at, ch.notes, ch.created_at, ch.updated_at
		FROM cleaning_history ch
		ORDER BY ch.cleaned_at DESC
	`

	rows, err := s.DB.Query(query)
	if err != nil {
		slog.Error("Failed to get all cleaning history", "error", err)
		return nil, err
	}
	defer rows.Close()

	var histories []CleaningHistory
	for rows.Next() {
		var history CleaningHistory

		err := rows.Scan(
			&history.ID,
			&history.ReleaseID,
			&history.CleanedAt,
			&history.Notes,
			&history.CreatedAt,
			&history.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan cleaning history", "error", err)
			return nil, err
		}

		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating cleaning history rows", "error", err)
		return histories, err
	}

	return histories, nil
}
