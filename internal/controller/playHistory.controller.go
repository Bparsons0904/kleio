package controller

import (
	"kleio/internal/database"
	"log/slog"
	"time"
)

func (c *Controller) CreatePlayHistory(history *database.PlayHistory) error {
	err := c.DB.CreatePlayHistory(history)
	if err != nil {
		slog.Error("Failed to create play history", "error", err)
		return err
	}

	return nil
}

func (c *Controller) UpdatePlayHistory(history *database.PlayHistory) error {
	err := c.DB.UpdatePlayHistory(history)
	if err != nil {
		slog.Error("Failed to update play history", "error", err)
		return err
	}

	return nil
}

func (c *Controller) DeletePlayHistory(id int) error {
	err := c.DB.DeletePlayHistory(id)
	if err != nil {
		slog.Error("Failed to delete play history", "error", err)
		return err
	}

	return nil
}

func (c *Controller) GetPlayCountByRelease() (map[int]int, error) {
	playCounts, err := c.DB.GetPlayCountByRelease()
	if err != nil {
		slog.Error("Failed to get play counts", "error", err)
		return nil, err
	}

	return playCounts, nil
}

func (c *Controller) GetRecentPlays(limit int) ([]database.PlayHistory, error) {
	plays, err := c.DB.GetRecentPlays(limit)
	if err != nil {
		slog.Error("Failed to get recent plays", "error", err)
		return nil, err
	}

	return plays, nil
}

func (c *Controller) GetPlaysByTimeRange(start, end time.Time) ([]database.PlayHistory, error) {
	plays, err := c.DB.GetPlaysByTimeRange(start, end)
	if err != nil {
		slog.Error("Failed to get plays by time range", "error", err)
		return nil, err
	}

	return plays, nil
}
