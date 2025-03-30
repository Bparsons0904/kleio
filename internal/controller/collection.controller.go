package controller

import (
	"log/slog"
)

func (c *Controller) GetCollection() (payload Payload, err error) {
	err = payload.GetPayload(c)
	if err != nil {
		slog.Error("Failed to get payload for collection", "error", err)
	}

	return payload, err
}

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

func (c *Controller) asyncCollection() {
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

func (c *Controller) SyncTrackAndDuration() {
	releases, err := c.DB.GetReleasesWithoutDuration()
	if err != nil {
		slog.Error("Failed to get releases without duration", "error", err)
		return
	}

	user, err := c.DB.GetUser()
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		return
	}

	err = c.GetReleaseDetails(releases[0].ResourceURL, user.Token)
	if err != nil {
		slog.Error("Failed to get track and duration", "error", err)
		return
	}

	//  for _, release := range releases {
	//    err := c.GetReleaseDetails(release.ResourceURL, user.Token)
	//    if err != nil {
	//      slog.Error("Failed to get track and duration", "error", err)
	//      continue
	//    }
	// }
}
