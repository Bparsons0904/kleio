package server

import (
	"encoding/json"
	"kleio/internal/database"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *Server) createPlayHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var history database.PlayHistory
	err := json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if history.PlayedAt.IsZero() {
		history.PlayedAt = time.Now()
	}

	if playedAtStr := r.URL.Query().Get("played_at"); playedAtStr != "" {
		playedAt, err := time.Parse(time.RFC3339, playedAtStr)
		if err != nil {
			slog.Error("Failed to parse played_at", "error", err)
			http.Error(w, "Invalid played_at format", http.StatusBadRequest)
			return
		}
		history.PlayedAt = playedAt
	}

	err = s.controller.CreatePlayHistory(&history)
	if err != nil {
		http.Error(w, "Failed to create play history", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeData(w, history)
}

func (s *Server) updatePlayHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var history database.PlayHistory
	err = json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID in path matches ID in body
	history.ID = id

	// Process played_at from query param if provided
	if playedAtStr := r.URL.Query().Get("played_at"); playedAtStr != "" {
		playedAt, err := time.Parse(time.RFC3339, playedAtStr)
		if err != nil {
			slog.Error("Failed to parse played_at", "error", err)
			http.Error(w, "Invalid played_at format", http.StatusBadRequest)
			return
		}
		history.PlayedAt = playedAt
	}

	err = s.controller.UpdatePlayHistory(&history)
	if err != nil {
		http.Error(w, "Failed to update play history", http.StatusInternalServerError)
		return
	}

	writeData(w, history)
}

func (s *Server) deletePlayHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = s.controller.DeletePlayHistory(id)
	if err != nil {
		http.Error(w, "Failed to delete play history", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getPlayCounts(w http.ResponseWriter, r *http.Request) {
	playCounts, err := s.controller.GetPlayCountByRelease()
	if err != nil {
		http.Error(w, "Failed to get play counts", http.StatusInternalServerError)
		return
	}

	writeData(w, playCounts)
}

func (s *Server) getRecentPlays(w http.ResponseWriter, r *http.Request) {
	// Parse limit parameter
	limit := 50 // Default limit

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	plays, err := s.controller.GetRecentPlays(limit)
	if err != nil {
		http.Error(w, "Failed to get recent plays", http.StatusInternalServerError)
		return
	}

	writeData(w, plays)
}

func (s *Server) getPlaysByTimeRange(w http.ResponseWriter, r *http.Request) {
	// Parse start and end times from query parameters
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		http.Error(w, "Missing start or end parameter", http.StatusBadRequest)
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	plays, err := s.controller.GetPlaysByTimeRange(start, end)
	if err != nil {
		http.Error(w, "Failed to get plays by time range", http.StatusInternalServerError)
		return
	}

	writeData(w, plays)
}
