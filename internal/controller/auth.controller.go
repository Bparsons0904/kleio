package controller

import (
	"log/slog"
	"time"
)

func (c *Controller) GetAuth() (time.Time, bool, error) {
	_, err := c.DB.GetToken()
	if err != nil {
		slog.Error("Failed to get token", "error", err)
		return time.Time{}, false, err
	}

	lastSync, err := c.DB.GetLatestSync()
	if err != nil {
		slog.Error("Failed to get last sync", "error", err)
		return time.Time{}, false, err
	}

	if lastSync.Status != "complete" {
		slog.Error("Last sync failed, re-syncing", "error", err)
		err := c.DB.CompleteSync(lastSync.ID, false)
		if err != nil {
			slog.Error("Failed to complete sync", "error", err)
		}
		go c.syncCollection()
		return time.Time{}, true, nil
	}

	expectedFolderSync := time.Now().Add(-12 * time.Hour)
	// expectedFolderSync := time.Now()
	if lastSync.SyncStart.Before(expectedFolderSync) {
		slog.Info("Last synced is older than 12 hours, updating folders...")
		go c.syncCollection()
		return lastSync.SyncEnd, true, nil
	}

	return lastSync.SyncEnd, false, nil
}

func (c *Controller) syncCollection() {
	id, err := c.DB.StartSync()
	if err != nil {
		slog.Error("Failed to start sync", "error", err)
		return
	}

	if err := c.SyncFolders(); err != nil {
		slog.Error("Failed to sync folders", "error", err)
		return
	}

	if err := c.SyncReleases(); err != nil {
		slog.Error("Failed to sync collection", "error", err)
		return
	}

	if err := c.DB.CompleteSync(id, true); err != nil {
		slog.Error("Failed to complete sync", "error", err)
		err := c.DB.CleanupAbandonedSyncs()
		if err != nil {
			slog.Error("Failed to cleanup abandoned syncs", "error", err)
		}
	}
}
