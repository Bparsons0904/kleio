package database

import (
	"fmt"
)

func (s *Database) SaveTracks(
	releaseID int,
	tracks []Track,
) error {
	_, err := s.DB.Exec("DELETE FROM tracks WHERE release_id = ?", releaseID)
	if err != nil {
		return fmt.Errorf("failed to delete existing tracks: %w", err)
	}

	for _, track := range tracks {
		_, err := s.DB.Exec(
			"INSERT INTO tracks (release_id, position, title, duration_text, duration_seconds) VALUES (?, ?, ?, ?, ?)",
			releaseID,
			track.Position,
			track.Title,
			track.DurationText,
			track.DurationSeconds,
		)
		if err != nil {
			return fmt.Errorf("failed to insert track: %w", err)
		}
	}

	return nil
}

func (s *Database) UpdateReleaseWithDetails(
	releaseID int,
	totalDuration int,
	estimated bool,
) error {
	_, err := s.DB.Exec(
		"UPDATE releases SET play_duration = ?, play_duration_estimated = ? WHERE id = ?",
		totalDuration,
		estimated,
		releaseID,
	)
	if err != nil {
		return fmt.Errorf("failed to update release duration: %w", err)
	}

	return nil
}
