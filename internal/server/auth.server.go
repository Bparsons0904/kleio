package server

import (
	"encoding/json"
	"kleio/internal/controller"
	"log/slog"
	"net/http"
)

func (s *Server) getAuth(w http.ResponseWriter, r *http.Request) {
	lastSync, syncingData, err := s.controller.GetAuth()
	if err != nil {
		http.Error(w, "Failed to get auth", http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"lastSync":    lastSync,
		"syncingData": syncingData,
	}

	writeData(w, resp)
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

	username, err := controller.GetUserIdentity(token)
	if err != nil {
		slog.Error("Failed to get user identity", "error", err)
		http.Error(w, "Failed to get user identity", http.StatusInternalServerError)
		return
	}

	slog.Info("Saving token...", "token", token)

	err = s.DB.SaveToken(token, username)
	if err != nil {
		slog.Error("Failed to save token", "error", err)
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]bool{"success": true}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to write response", "error", err)
	}
}
