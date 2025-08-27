package server

import (
	"log/slog"
	"net/http"
	"time"
)

func (s *Server) getCollection(w http.ResponseWriter, r *http.Request) {
	releases, err := s.controller.GetCollection()
	if err != nil {
		http.Error(w, "Failed to get collection", http.StatusInternalServerError)
		return
	}

	writeData(w, releases)
}

func (s *Server) updateCollection(w http.ResponseWriter, r *http.Request) {
	slog.Info("Collection sync requested by user")

	// Check if there's already a sync in progress
	latestSync, err := s.DB.GetLatestSync()
	if err != nil {
		slog.Error("Failed to check latest sync before starting new one", "error", err)
		http.Error(w, "Failed to check sync status", http.StatusInternalServerError)
		return
	}

	if latestSync.Status == "in_progress" {
		slog.Warn("Sync already in progress, not starting new one", 
			"syncID", latestSync.ID,
			"syncStarted", latestSync.SyncStart)
		response := map[string]any{
			"isSyncing": true,
			"message": "Sync already in progress",
			"syncID": latestSync.ID,
		}
		writeData(w, response)
		return
	}

	slog.Info("Starting new collection sync in background")
	go func() {
		err := s.controller.AsyncCollection()
		if err != nil {
			slog.Error("Background sync failed", "error", err)
		} else {
			slog.Info("Background sync completed successfully")
		}
	}()

	response := map[string]any{"isSyncing": true, "message": "Sync started"}
	writeData(w, response)
}

func (s *Server) checkSync(w http.ResponseWriter, r *http.Request) {
	sync, err := s.DB.GetLatestSync()
	if err != nil {
		http.Error(w, "Failed to check sync", http.StatusInternalServerError)
		return
	}

	if sync.Status == "in_progress" && time.Since(sync.SyncStart) > 1*time.Minute {
		err := s.DB.CompleteSync(sync.ID, false)
		if err != nil {
			slog.Error("Failed to complete sync", "error", err)
			err := s.DB.CleanupAbandonedSyncs()
			if err != nil {
				slog.Error("Failed to cleanup abandoned syncs", "error", err)
			}
			http.Error(w, "Failed to complete sync", http.StatusInternalServerError)
			return
		}
	}

	writeData(w, sync)
}
