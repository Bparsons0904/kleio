package controller

import (
	"kleio/internal/database"
	"log/slog"
	"time"
)

type AuthPayload struct {
	LastSync    time.Time          `json:"lastSync"`
	SyncingData bool               `json:"syncingData"`
	Releases    []database.Release `json:"releases"`
}

func (c *Controller) GetAuth() (AuthPayload, error) {
	var payload AuthPayload
	_, err := c.DB.GetToken()
	if err != nil {
		slog.Error("Failed to get token", "error", err)
		return payload, err
	}

	lastSync, err := c.DB.GetLatestSync()
	if err != nil {
		slog.Error("Failed to get last sync", "error", err)
		return payload, err
	}

	if lastSync.Status != "complete" {
		slog.Error("Last sync failed, re-syncing", "error", err)
		err := c.DB.CompleteSync(lastSync.ID, false)
		if err != nil {
			slog.Error("Failed to complete sync", "error", err)
		}
		go c.syncCollection()
		return payload, nil
	}

	payload.Releases, err = c.DB.GetAllReleases()
	if err != nil {
		slog.Error("Failed to get releases", "error", err)
		return payload, err
	}

	expectedFolderSync := time.Now().Add(-12 * time.Hour)
	if lastSync.SyncStart.Before(expectedFolderSync) {
		slog.Info("Last synced is older than 12 hours, updating folders...")
		go c.syncCollection()
		payload.SyncingData = true
		return payload, nil
	}

	payload.LastSync = lastSync.SyncStart
	return payload, nil
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
