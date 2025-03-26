package controller

import (
	"log/slog"
	"time"
)

func (c *Controller) GetAuth() (time.Time, time.Time, error) {
	_, err := c.DB.GetToken()
	if err != nil {
		slog.Error("Failed to get token", "error", err)
		return time.Time{}, time.Time{}, err
	}

	lastFolderSync, err := c.DB.GetLastFolderSync()
	if err != nil {
		slog.Error("Failed to get last sync", "error", err)
		return time.Time{}, time.Time{}, err
	}

	now := time.Now()

	expectedFolderSync := now.Add(-24 * time.Hour)
	if lastFolderSync.Before(expectedFolderSync) {
		slog.Info("Last synced is older than 24 hours, updating folders...")
		if err := c.SyncFolders(); err != nil {
			slog.Error("Failed to sync folders", "error", err)
			return time.Time{}, time.Time{}, err
		}
		lastFolderSync = now
	}

	lastReleaseSync, err := c.DB.GetLastReleaseSync()
	expectedCollectionSync := now.Add(-12 * time.Hour)
	if lastFolderSync.Before(expectedCollectionSync) {
		slog.Info("Last synced is older than 12 hours, updating collection...")
		if err := c.SyncReleases(); err != nil {
			slog.Error("Failed to sync collection", "error", err)
			return time.Time{}, time.Time{}, err
		}
		lastReleaseSync = now
	}

	return lastFolderSync, lastReleaseSync, nil
}
