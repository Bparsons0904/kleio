package controller

import (
	"kleio/internal/database"
	"log/slog"
	"time"
)

func (c *Controller) CreateCleaningHistory(
	history *database.CleaningHistory,
) (payload Payload, err error) {
	err = c.DB.CreateCleaningHistory(history)
	if err != nil {
		slog.Error("Failed to create cleaning history", "error", err)
		return
	}

	err = payload.GetPayload(c)
	if err != nil {
		slog.Error("Failed to get payload for cleaning history", "error", err)
	}

	return
}

func (c *Controller) UpdateCleaningHistory(
	history *database.CleaningHistory,
) (payload Payload, err error) {
	err = c.DB.UpdateCleaningHistory(history)
	if err != nil {
		slog.Error("Failed to update cleaning history", "error", err)
		return
	}

	err = payload.GetPayload(c)
	if err != nil {
		slog.Error("Failed to get payload for cleaning history", "error", err)
	}

	return
}

func (c *Controller) DeleteCleaningHistory(id int) (payload Payload, err error) {
	err = c.DB.DeleteCleaningHistory(id)
	if err != nil {
		slog.Error("Failed to delete cleaning history", "error", err)
		return
	}

	err = payload.GetPayload(c)
	if err != nil {
		slog.Error("Failed to get payload for cleaning history", "error", err)
	}

	return
}

func (c *Controller) GetCleaningsByTimeRange(
	start, end time.Time,
) ([]database.CleaningHistory, error) {
	cleanings, err := c.DB.GetCleaningsByTimeRange(start, end)
	if err != nil {
		slog.Error("Failed to get cleanings by time range", "error", err)
		return nil, err
	}

	return cleanings, nil
}

func (c *Controller) CountCleaningsByRelease() (map[int]int, error) {
	counts, err := c.DB.CountCleaningsByRelease()
	if err != nil {
		slog.Error("Failed to count cleanings by release", "error", err)
		return nil, err
	}

	return counts, nil
}
