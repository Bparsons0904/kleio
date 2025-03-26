package server

import (
	"kleio/internal/controller"
	"net/http"
)

func (s *Server) updateCollection(w http.ResponseWriter, r *http.Request) {
	releases, err := controller.UpdateCollection(s.db)
	if err != nil {
		http.Error(w, "Failed to update collection", http.StatusInternalServerError)
		return
	}

	writeData(w, releases)
}
