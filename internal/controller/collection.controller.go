package controller

import (
	"fmt"
	"kleio/internal/database"
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

	if err := c.syncTracksAndDuration(); err != nil {
		slog.Error("Failed to sync tracks and duration", "error", err)
		return err
	}

	return nil
}

func (c *Controller) AsyncCollection() (err error) {
	slog.Info("AsyncCollection started")
	defer func() {
		if r := recover(); r != nil {
			slog.Error("AsyncCollection panicked", "panic", r)
			err = fmt.Errorf("sync panicked: %v", r)
		}
		if err != nil {
			slog.Error("AsyncCollection failed", "error", err)
			err := c.DB.CleanupAbandonedSyncs()
			if err != nil {
				slog.Error("Failed to cleanup abandoned syncs", "error", err)
			}
		} else {
			slog.Info("AsyncCollection completed successfully")
		}
	}()

	// Check if sync is already in progress
	latestSync, err := c.DB.GetLatestSync()
	if err != nil {
		slog.Error("Failed to check latest sync", "error", err)
		return err
	}

	if latestSync.Status == "in_progress" {
		slog.Warn("Another sync is already in progress", 
			"syncID", latestSync.ID,
			"syncStarted", latestSync.SyncStart)
		return fmt.Errorf("sync already in progress (ID: %d)", latestSync.ID)
	}

	id, err := c.DB.StartSync()
	if err != nil {
		slog.Error("Failed to start sync", "error", err)
		return err
	}

	slog.Info("Starting full collection sync", "syncID", id)

	if err = c.SyncCollection(); err != nil {
		slog.Error("Failed to sync collection", "error", err, "syncID", id)
		return c.DB.CompleteSync(id, false)
	}

	if err := c.DB.CompleteSync(id, true); err != nil {
		slog.Error("Failed to complete sync", "error", err, "syncID", id)
		return err
	}

	slog.Info("Full collection sync completed", "syncID", id)
	return nil
}

func (c *Controller) syncTracksAndDuration() error {
	releases, err := c.DB.GetReleasesWithoutDuration()
	if err != nil {
		slog.Error("Failed to get releases without duration", "error", err)
		return err
	}

	slog.Info("Starting track and duration sync", "totalReleases", len(releases))

	user, err := c.DB.GetUser()
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		return err
	}

	processed := 0
	failed := 0

	for i, release := range releases {
		slog.Info("Processing release", 
			"releaseID", release.ID, 
			"title", release.Title,
			"progress", fmt.Sprintf("%d/%d", i+1, len(releases)),
			"processed", processed,
			"failed", failed)

		err := c.processReleaseTracks(release, user)
		if err != nil {
			slog.Error("Failed to process release tracks", 
				"error", err, 
				"releaseID", release.ID,
				"title", release.Title)
			failed++
			// Continue processing other releases instead of failing completely
			continue
		}
		processed++

		if c.RateLimit.ShouldThrottle() {
			current := c.RateLimit.GetCurrent()
			slog.Info("Rate limit throttling", 
				"remaining", current.Remaining,
				"used", current.Used,
				"limit", current.Limit)
			time.Sleep(15 * time.Second)
		}
	}

	slog.Info("Completed track and duration sync", 
		"totalReleases", len(releases),
		"processed", processed,
		"failed", failed)

	if failed > 0 && processed == 0 {
		return fmt.Errorf("failed to process any releases: %d failures", failed)
	}

	return nil
}

func (c *Controller) processReleaseTracks(release database.Release, user database.User) error {
	slog.Debug("Starting track processing", 
		"releaseID", release.ID,
		"title", release.Title,
		"resourceURL", release.ResourceURL)

	tracks, err := c.GetReleaseDetails(release, user.Token)
	if err != nil {
		slog.Error("processReleaseTracks: Failed to get track and duration", 
			"error", err,
			"releaseID", release.ID,
			"title", release.Title,
			"resourceURL", release.ResourceURL)
		return err
	}

	slog.Debug("Retrieved tracks from API", 
		"releaseID", release.ID,
		"trackCount", len(tracks))

	err = c.DB.SaveTracks(release.ID, tracks)
	if err != nil {
		slog.Error("processReleaseTracks: Failed to save tracks", 
			"error", err,
			"releaseID", release.ID,
			"trackCount", len(tracks))
		return err
	}

	slog.Debug("Saved tracks to database", "releaseID", release.ID)

	durationSeconds, isDurationEstimated, err := c.calculateTrackDurations(release.ID, tracks)
	if err != nil {
		slog.Error("processReleaseTracks: Failed to calculate track durations", 
			"error", err,
			"releaseID", release.ID)
		return err
	}

	slog.Info("Calculated release duration", 
		"releaseID", release.ID,
		"title", release.Title,
		"durationSeconds", durationSeconds, 
		"durationMinutes", fmt.Sprintf("%.1f", float64(durationSeconds)/60),
		"estimated", isDurationEstimated)

	if durationSeconds == 0 {
		slog.Warn("Zero duration calculated for release",
			"releaseID", release.ID,
			"title", release.Title,
			"trackCount", len(tracks),
			"estimated", isDurationEstimated,
			"resourceURL", release.ResourceURL)
		// Don't return error for zero duration - this might be expected for some releases
	}

	err = c.DB.UpdateReleaseWithDetails(release.ID, durationSeconds, isDurationEstimated)
	if err != nil {
		slog.Error("processReleaseTracks: Failed to update release duration", 
			"error", err,
			"releaseID", release.ID,
			"durationSeconds", durationSeconds)
		return err
	}

	slog.Debug("Successfully processed release tracks", "releaseID", release.ID)
	return nil
}
