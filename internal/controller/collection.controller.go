package controller

import (
	"log/slog"
)

func (c *Controller) SyncCollection() error {
	if err := c.SyncFolders(); err != nil {
		slog.Error("Failed to sync folders", "error", err)
		return err
	}

	if err := c.SyncReleases(); err != nil {
		slog.Error("Failed to sync collection", "error", err)
		return err
	}

	return nil
}
