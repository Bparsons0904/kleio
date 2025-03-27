package server

import (
	"log/slog"
	"net/http"
	"time"
)

func (s *Server) updateCollection(w http.ResponseWriter, r *http.Request) {
	err := s.controller.SyncCollection()
	if err != nil {
		http.Error(w, "Failed to update collection", http.StatusInternalServerError)
		return
	}

	// writeData(w, releases)
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
