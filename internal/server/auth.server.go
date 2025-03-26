package server

import (
	"encoding/json"
	"kleio/internal/controller"
	"log/slog"
	"net/http"
)

func (s *Server) getAuth(w http.ResponseWriter, r *http.Request) {
	folderLastSync, collectionLastSync, err := s.controller.GetAuth()
	if err != nil {
		http.Error(w, "Failed to get auth", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"folderLastSync":     folderLastSync.String(),
		"collectionLastSync": collectionLastSync.String(),
	}

	writeData(w, resp)
}

func (s *Server) SaveToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, err := getToken(r)
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		slog.Error("Failed to get token", "error", err)
		return
	}

	username, err := controller.GetUserIdentity(token)
	if err != nil {
		http.Error(w, "Failed to get user identity", http.StatusInternalServerError)
		slog.Error("Failed to get user identity", "error", err)
		return
	}

	slog.Info("Saving token...", "token", token)

	err = s.DB.SaveToken(token, username)
	if err != nil {
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		slog.Error("Failed to save token", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]bool{"success": true}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to write response", "error", err)
	}
}
