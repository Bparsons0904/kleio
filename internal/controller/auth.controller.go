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

	lastFolderSync, err := c.DB.GetLastFolderSync()
	if err != nil {
		slog.Error("Failed to get last sync", "error", err)
		return time.Time{}, false, err
	}

	now := time.Now()
	syncingData := false
	expectedFolderSync := now.Add(-12 * time.Hour)
	if lastFolderSync.Before(expectedFolderSync) {
		syncingData = true
		slog.Info("Last synced is older than 12 hours, updating folders...")
		go c.syncCollection()
	}

	return now, syncingData, nil
}

func (c *Controller) syncCollection() {
	if err := c.SyncFolders(); err != nil {
		slog.Error("Failed to sync folders", "error", err)
		return
	}

	if err := c.SyncReleases(); err != nil {
		slog.Error("Failed to sync collection", "error", err)
		return
	}
}
