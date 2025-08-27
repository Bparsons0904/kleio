package database

import (
	"database/sql"
	"log/slog"
)

func (s *Database) GetLatestSync() (Sync, error) {
	slog.Debug("Retrieving latest sync record")

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
			slog.Debug("No sync records found in database")
			return Sync{}, nil
		}
		slog.Error("Database query error for latest sync", "error", err, "query", query)
		return Sync{}, err
	}

	slog.Debug("Retrieved latest sync record", 
		"syncID", sync.ID,
		"status", sync.Status,
		"syncStart", sync.SyncStart)

	return sync, nil
}

func (s *Database) StartSync() (int64, error) {
	slog.Info("Starting new sync operation")

	query := `
        INSERT INTO syncs (sync_start, status)
        VALUES (CURRENT_TIMESTAMP, 'in_progress')
    `
	result, err := s.DB.Exec(query)
	if err != nil {
		slog.Error("Failed to insert new sync record", "error", err, "query", query)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error(
			"Failed to get last insert id for sync",
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

	slog.Info("Created new sync record", "syncID", id)
	return id, nil
}

func (s *Database) CompleteSync(id int64, success bool) error {
	status := "failed"
	if success {
		status = "complete"
	}

	slog.Info("Completing sync operation", "syncID", id, "success", success, "status", status)

	var query string
	if success {
		query = `
			UPDATE syncs
			SET sync_end = CURRENT_TIMESTAMP, status = 'complete'
			WHERE id = ?`
	} else {
		query = `
			UPDATE syncs
			SET sync_end = CURRENT_TIMESTAMP, status = 'failed'
			WHERE id = ?`
	}

	result, err := s.DB.Exec(query, id)
	if err != nil {
		slog.Error(
			"Failed to complete sync in database",
			"error",
			err,
			"syncID",
			id,
			"success",
			success,
			"query",
			query,
		)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		slog.Warn("No sync record was updated", "syncID", id)
	} else {
		slog.Info("Successfully completed sync", "syncID", id, "status", status)
	}

	return nil
}

func (s *Database) CleanupAbandonedSyncs() error {
	slog.Info("Starting cleanup of abandoned syncs")

	// First, check if there are any abandoned syncs
	countQuery := `SELECT COUNT(*) FROM syncs WHERE status = 'in_progress'`
	var count int
	err := s.DB.QueryRow(countQuery).Scan(&count)
	if err != nil {
		slog.Error("Failed to count abandoned syncs", "error", err)
		return err
	}

	if count == 0 {
		slog.Info("No abandoned syncs found")
		return nil
	}

	slog.Info("Found abandoned syncs", "count", count)

	query := `
		UPDATE syncs
		SET status = 'failed', sync_end = CURRENT_TIMESTAMP
		WHERE status = 'in_progress'`

	result, err := s.DB.Exec(query)
	if err != nil {
		slog.Error("Failed to update abandoned syncs, attempting deletion", "error", err)

		deleteQuery := `
			DELETE FROM syncs
			WHERE status = 'in_progress'`

		deleteResult, deleteErr := s.DB.Exec(deleteQuery)
		if deleteErr != nil {
			slog.Error("Failed to delete abandoned syncs", "error", deleteErr)
			return err
		}

		deleted, _ := deleteResult.RowsAffected()
		slog.Info("Successfully deleted abandoned syncs after update failure", "deleted", deleted)
		return nil
	}

	updated, _ := result.RowsAffected()
	slog.Info("Successfully marked abandoned syncs as failed", "updated", updated)
	return nil
}
