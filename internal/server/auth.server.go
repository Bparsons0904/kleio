package server

import (
	"log/slog"
	"net/http"
)

func (s *Server) getAuth(w http.ResponseWriter, r *http.Request) {
	payload, err := s.controller.GetAuth()
	if err != nil {
		http.Error(w, "Failed to get auth", http.StatusInternalServerError)
		return
	}

	writeData(w, payload)
}

func (s *Server) SaveToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Error("Method not allowed", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, err := getToken(r)
	if err != nil {
		slog.Error("Failed to get token", "error", err)
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}

	payload, err := s.controller.SaveToken(token)
	if err != nil {
		http.Error(w, "Failed to get auth", http.StatusInternalServerError)
		return
	}

	writeData(w, payload)
}
