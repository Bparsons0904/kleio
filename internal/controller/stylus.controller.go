package controller

import (
	"kleio/internal/database"
	"log/slog"
)

func (c *Controller) GetStyluses() ([]database.Stylus, error) {
	styluses, err := c.DB.GetStyluses()
	if err != nil {
		slog.Error("Failed to get styluses", "error", err)
		return nil, err
	}

	return styluses, nil
}

func (c *Controller) CreateStylus(stylus *database.Stylus) error {
	err := c.DB.CreateStylus(stylus)
	if err != nil {
		slog.Error("Failed to create stylus", "error", err)
		return err
	}

	return nil
}

func (c *Controller) UpdateStylus(stylus *database.Stylus) error {
	err := c.DB.UpdateStylus(stylus)
	if err != nil {
		slog.Error("Failed to update stylus", "error", err)
		return err
	}

	return nil
}

func (c *Controller) DeleteStylus(id int) error {
	err := c.DB.DeleteStylus(id)
	if err != nil {
		slog.Error("Failed to delete stylus", "error", err)
		return err
	}

	return nil
}
