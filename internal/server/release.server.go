package server

import "net/http"

func (s *Server) deleteRelease(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := getIDFrom3_2Parts(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	releases, err := s.controller.DeleteRelease(id)
	if err != nil {
		http.Error(w, "Failed to delete release", http.StatusInternalServerError)
		return
	}

	writeData(w, releases)
}

func (s *Server) archiveRelease(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := getIDFrom3_2Parts(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	releases, err := s.controller.ArchiveRelease(id)
	if err != nil {
		http.Error(w, "Failed to archive release", http.StatusInternalServerError)
		return
	}

	writeData(w, releases)
}
