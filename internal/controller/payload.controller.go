package controller

import (
	"kleio/internal/database"
	"log/slog"
	"sort"
	"time"
)

type Payload struct {
	LastSync    time.Time              `json:"lastSync,omitzero"`
	SyncingData bool                   `json:"syncingData"`
	Releases    []database.Release     `json:"releases"`
	Stylus      []database.Stylus      `json:"stylus"`
	PlayHistory []database.PlayHistory `json:"playHistory"`
	Folders     []database.Folder      `json:"folders"`
	Token       string                 `json:"token"`
}

func (p *Payload) GetLastSync(controller *Controller) error {
	lastSync, err := controller.DB.GetLatestSync()
	if err != nil {
		slog.Error("Failed to get last sync", "error", err)
		return err
	}

	if lastSync.Status != "complete" {
		slog.Error("Last sync failed, re-syncing", "error", err)
		err := controller.DB.CompleteSync(lastSync.ID, false)
		if err != nil {
			slog.Error("Failed to complete sync", "error", err)
		}
		go controller.AsyncCollection()
		return nil
	}

	p.LastSync = lastSync.SyncStart

	expectedFolderSync := time.Now().Add(-12 * time.Hour)
	// expectedFolderSync := time.Now()
	if lastSync.SyncStart.Before(expectedFolderSync) {
		slog.Info("Last synced is older than 12 hours, updating folders...")
		go controller.AsyncCollection()
		p.SyncingData = true
	}

	return nil
}

func (p *Payload) GetPayload(controller *Controller) (err error) {
	p.Releases, err = controller.DB.GetAllReleases()
	if err != nil {
		slog.Error("Failed to get releases", "error", err)
		return err
	}

	p.GetPlayHistory(controller)

	p.Stylus, err = controller.DB.GetStyluses()
	if err != nil {
		slog.Error("Failed to get stylus", "error", err)
		return err
	}

	err = p.GetLastSync(controller)
	if err != nil {
		slog.Error("Failed to get last sync", "error", err)
		return err
	}

	err = p.GetFolders(controller)
	if err != nil {
		slog.Error("Failed to get folders", "error", err)
		return err
	}

	return nil
}

func (p *Payload) GetFolders(controller *Controller) (err error) {
	p.Folders, err = controller.DB.GetFolders()
	if err != nil {
		slog.Error("Failed to get folders", "error", err)
		return err
	}

	return nil
}

func (p *Payload) GetPlayHistory(controller *Controller) {
	var playHistory []database.PlayHistory
	for _, release := range p.Releases {
		for _, history := range release.PlayHistory {
			history.Release = release
			playHistory = append(playHistory, history)
		}
	}

	sort.Slice(playHistory, func(i, j int) bool {
		return playHistory[i].PlayedAt.After(playHistory[j].PlayedAt)
	})
	p.PlayHistory = playHistory
}
