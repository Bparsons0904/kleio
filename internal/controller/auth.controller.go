package controller

import (
	"kleio/internal/database"
	"log/slog"
	"time"
)

type AuthPayload struct {
	LastSync    time.Time          `json:"lastSync,omitzero"`
	SyncingData bool               `json:"syncingData"`
	Releases    []database.Release `json:"releases"`
	Stylus      []database.Stylus  `json:"stylus"`
	Token       string             `json:"token"`
}

func (c *Controller) GetAuth() (payload AuthPayload, err error) {
	payload.Token, err = c.DB.GetToken()
	if err != nil {
		slog.Error("Failed to get token", "error", err)
		return AuthPayload{}, nil
	}

	lastSync, err := c.DB.GetLatestSync()
	if err != nil {
		slog.Error("Failed to get last sync", "error", err)
		return payload, err
	}

	if lastSync.Status == "" && lastSync.Status != "complete" {
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

	payload.Stylus, err = c.DB.GetStyluses()
	if err != nil {
		slog.Error("Failed to get stylus", "error", err)
		return payload, err
	}

	expectedFolderSync := time.Now().Add(-12 * time.Hour)
	// expectedFolderSync := time.Now()
	if lastSync.SyncStart.Before(expectedFolderSync) {
		slog.Info("Last synced is older than 12 hours, updating folders...")
		go c.syncCollection()
		payload.SyncingData = true
		return payload, nil
	}

	payload.LastSync = lastSync.SyncStart
	return payload, nil
}

func (c *Controller) SaveToken(token string) (payload AuthPayload, err error) {
	username, err := GetUserIdentity(token)
	if err != nil {
		slog.Error("Failed to get user identity", "error", err)
		return payload, err
	}

	err = c.DB.SaveToken(token, username)
	if err != nil {
		slog.Error("Failed to save token", "error", err)
		return payload, err
	}

	payload, err = c.GetAuth()
	if err != nil {
		slog.Error("Failed to get auth", "error", err)
		return payload, err
	}

	go c.syncCollection()
	payload.SyncingData = true

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
