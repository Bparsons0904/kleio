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

func (s *Server) createCleaningHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var history database.CleaningHistory
	err := json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if history.CleanedAt.IsZero() {
		history.CleanedAt = time.Now()
	}

	if cleanedAtStr := r.URL.Query().Get("cleaned_at"); cleanedAtStr != "" {
		cleanedAt, err := time.Parse(time.RFC3339, cleanedAtStr)
		if err != nil {
			slog.Error("Failed to parse cleaned_at", "error", err)
			http.Error(w, "Invalid cleaned_at format", http.StatusBadRequest)
			return
		}
		history.CleanedAt = cleanedAt
	}

	payload, err := s.controller.CreateCleaningHistory(&history)
	if err != nil {
		http.Error(w, "Failed to create cleaning history", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeData(w, payload)
}

func (s *Server) updateCleaningHistory(w http.ResponseWriter, r *http.Request) {
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

	var history database.CleaningHistory
	err = json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID in path matches ID in body
	history.ID = id

	// Process cleaned_at from query param if provided
	if cleanedAtStr := r.URL.Query().Get("cleaned_at"); cleanedAtStr != "" {
		cleanedAt, err := time.Parse(time.RFC3339, cleanedAtStr)
		if err != nil {
			slog.Error("Failed to parse cleaned_at", "error", err)
			http.Error(w, "Invalid cleaned_at format", http.StatusBadRequest)
			return
		}
		history.CleanedAt = cleanedAt
	}

	payload, err := s.controller.UpdateCleaningHistory(&history)
	if err != nil {
		http.Error(w, "Failed to update cleaning history", http.StatusInternalServerError)
		return
	}

	writeData(w, payload)
}

func (s *Server) deleteCleaningHistory(w http.ResponseWriter, r *http.Request) {
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

	payload, err := s.controller.DeleteCleaningHistory(id)
	if err != nil {
		http.Error(w, "Failed to delete cleaning history", http.StatusInternalServerError)
		return
	}

	writeData(w, payload)
}

func (s *Server) getCleaningCountsByRelease(w http.ResponseWriter, r *http.Request) {
	counts, err := s.controller.CountCleaningsByRelease()
	if err != nil {
		http.Error(w, "Failed to get cleaning counts", http.StatusInternalServerError)
		return
	}

	writeData(w, counts)
}

func (s *Server) getCleaningsByTimeRange(w http.ResponseWriter, r *http.Request) {
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

	cleanings, err := s.controller.GetCleaningsByTimeRange(start, end)
	if err != nil {
		http.Error(w, "Failed to get cleanings by time range", http.StatusInternalServerError)
		return
	}

	writeData(w, cleanings)
}
