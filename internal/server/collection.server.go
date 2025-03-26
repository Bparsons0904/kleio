package server

import (
	"net/http"
)

func (s *Server) updateCollection(w http.ResponseWriter, r *http.Request) {
	err := s.controller.SyncCollection()
	if err != nil {
		http.Error(w, "Failed to update collection", http.StatusInternalServerError)
		return
	}

	// writeData(w, releases)
}
