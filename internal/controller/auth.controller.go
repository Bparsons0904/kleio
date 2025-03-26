package controller

import (
	"log"
	"log/slog"
	"time"
)

func (c *Controller) GetAuth() (time.Time, time.Time, error) {
	token, err := c.DB.GetToken()
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
		c.SyncFolders()
	}

	expectedCollectionSync := now.Add(-12 * time.Hour)
	if lastFolderSync.Before(expectedCollectionSync) {
		slog.Info("Last synced is older than 12 hours, updating collection...")
		c.SyncReleases()
		// UpdateCollection()
	}

	// go controller.UpdateCollection(s.db)

	log.Printf("token: %s", token)
	return lastFolderSync, now, nil
}
