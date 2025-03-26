package server

import (
	"net/http"
)

func (s *Server) updateCollection(w http.ResponseWriter, r *http.Request) {
	// releases, err := controller.UpdateCollection(s.DB)
	// if err != nil {
	// 	http.Error(w, "Failed to update collection", http.StatusInternalServerError)
	// 	return
	// }
	//
	// writeData(w, releases)
}
