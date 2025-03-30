package controller

import (
	"kleio/internal/database"
	"log"
	"log/slog"
	"time"
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

func (c *Controller) SyncTracksAndDuration() {
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

	count := 0
	for _, release := range releases {
		count++
		slog.Info("Processing release", "releaseID", release.ID)
		err := c.processReleaseTracks(release, user)
		if err != nil {
			slog.Error("Failed to process release tracks", "error", err)
			break
		}

		if c.RateLimit.ShouldThrottle() {
			time.Sleep(15 * time.Second)
		}
	}
}

func (c *Controller) processReleaseTracks(release database.Release, user database.User) error {
	tracks, err := c.GetReleaseDetails(release, user.Token)
	if err != nil {
		slog.Error("processReleaseTracks: Failed to get track and duration", "error", err)
		return err
	}

	err = c.DB.SaveTracks(release.ID, tracks)
	if err != nil {
		slog.Error("processReleaseTracks: Failed to save tracks", "error", err)
		return err
	}

	durationSeconds, isDurationEstimated, err := c.calculateTrackDurations(release.ID, tracks)
	if err != nil {
		slog.Error("processReleaseTracks: Failed to calculate track durations", "error", err)
		return err
	}

	log.Println("duration", "seconds", durationSeconds, "estimated", isDurationEstimated)
	if durationSeconds == 0 {
		slog.Error(
			"processReleaseTracks: Duration",
			"seconds",
			durationSeconds,
			"estimated",
			isDurationEstimated,
			"resourceURL",
			release.ResourceURL,
		)
		return err
	}

	err = c.DB.UpdateReleaseWithDetails(release.ID, durationSeconds, isDurationEstimated)
	if err != nil {
		slog.Error("processReleaseTracks: Failed to update release duration", "error", err)
		return err
	}

	return nil
}
