package database

import (
	"database/sql"
	"log/slog"
)

func (s *Database) GetLatestSync() (Sync, error) {
	var sync Sync
	query := `
		SELECT id, sync_start, sync_end, status
		FROM syncs
		ORDER BY id DESC
		LIMIT 1`

	err := s.DB.QueryRow(query).Scan(
		&sync.ID,
		&sync.SyncStart,
		&sync.SyncEnd,
		&sync.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Sync{}, nil
		}
		slog.Error("Database query error", "error", err, "query", query)
		return Sync{}, err
	}
	return sync, nil
}

func (s *Database) StartSync() (int64, error) {
	query := `
        INSERT INTO syncs (sync_start, status)
        VALUES (CURRENT_TIMESTAMP, 'in_progress')
    `
	result, err := s.DB.Exec(query)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error(
			"Failed to get last insert id",
			"error",
			err,
			"query",
			query,
			"result",
			result,
			"lastInsertId",
			id,
		)
		return 0, err
	}

	return id, nil
}

func (s *Database) CompleteSync(id int64, success bool) error {
	var query string
	if success {
		query = `
			UPDATE syncs
			SET sync_end = CURRENT_TIMESTAMP, status = 'complete'
			WHERE id = ?`
	} else {
		query = `
			UPDATE syncs
			SET status = 'failed'
			WHERE id = ?`
	}

	_, err := s.DB.Exec(query, id)
	if err != nil {
		slog.Error(
			"Failed to complete sync",
			"error",
			err,
			"id",
			id,
			"success",
			success,
			"query",
			query,
		)
	}
	return err
}

func (s *Database) CleanupAbandonedSyncs() error {
	query := `
		UPDATE syncs
		SET status = 'failed'
		WHERE status = 'in_progress'`

	_, err := s.DB.Exec(query)
	if err != nil {
		slog.Error("Failed to update abandoned syncs, attempting deletion", "error", err)

		deleteQuery := `
			DELETE FROM syncs
			WHERE status = 'in_progress'`

		_, deleteErr := s.DB.Exec(deleteQuery)
		if deleteErr != nil {
			slog.Error("Failed to delete abandoned syncs", "error", deleteErr)
			return err
		}

		slog.Info("Successfully deleted abandoned syncs after update failure")
		return nil
	}

	return nil
}
