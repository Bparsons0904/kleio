// internal/controller/export.controller.go
package controller

import (
	"kleio/internal/database"
	"log/slog"
	"time"
)

type ExportPlayHistory struct {
	ReleaseID int       `json:"releaseId"`
	StylusID  *int      `json:"stylusId,omitempty"`
	PlayedAt  time.Time `json:"playedAt"`
	Notes     string    `json:"notes,omitempty"`
}

type ExportCleaningHistory struct {
	ReleaseID int       `json:"releaseId"`
	CleanedAt time.Time `json:"cleanedAt"`
	Notes     string    `json:"notes,omitempty"`
}

type ExportData struct {
	ExportDate      time.Time               `json:"exportDate"`
	PlayHistory     []ExportPlayHistory     `json:"playHistory"`
	CleaningHistory []ExportCleaningHistory `json:"cleaningHistory"`
	Styluses        []database.Stylus       `json:"styluses"`
}

func (c *Controller) ExportHistory() (ExportData, error) {
	// Get all play history
	playHistory, err := c.DB.GetAllPlayHistory()
	if err != nil {
		slog.Error("Failed to get play history for export", "error", err)
		return ExportData{}, err
	}

	// Get all cleaning history
	cleaningHistory, err := c.DB.GetAllCleaningHistory()
	if err != nil {
		slog.Error("Failed to get cleaning history for export", "error", err)
		return ExportData{}, err
	}

	// Get all styluses (owned only, not base models)
	styluses, err := c.DB.GetStyluses()
	if err != nil {
		slog.Error("Failed to get styluses for export", "error", err)
		return ExportData{}, err
	}

	// Filter styluses to only include owned and not base models
	var filteredStyluses []database.Stylus
	for _, stylus := range styluses {
		if stylus.Owned && !stylus.BaseModel {
			filteredStyluses = append(filteredStyluses, stylus)
		}
	}

	// Create simplified play history records
	var exportPlayHistory []ExportPlayHistory
	for _, play := range playHistory {
		exportPlayHistory = append(exportPlayHistory, ExportPlayHistory{
			ReleaseID: play.ReleaseID,
			StylusID:  play.StylusID,
			PlayedAt:  play.PlayedAt,
			Notes:     play.Notes,
		})
	}

	// Create simplified cleaning history records
	var exportCleaningHistory []ExportCleaningHistory
	for _, cleaning := range cleaningHistory {
		exportCleaningHistory = append(exportCleaningHistory, ExportCleaningHistory{
			ReleaseID: cleaning.ReleaseID,
			CleanedAt: cleaning.CleanedAt,
			Notes:     cleaning.Notes,
		})
	}

	// Create export data structure
	exportData := ExportData{
		ExportDate:      time.Now(),
		PlayHistory:     exportPlayHistory,
		CleaningHistory: exportCleaningHistory,
		Styluses:        filteredStyluses,
	}

	return exportData, nil
}
