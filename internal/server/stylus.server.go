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

func (s *Server) createStylus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var stylus database.Stylus
	err := json.NewDecoder(r.Body).Decode(&stylus)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if purchaseDateStr, ok := r.URL.Query()["purchase_date"]; ok && len(purchaseDateStr) > 0 {
		purchaseDate, err := time.Parse(time.RFC3339, purchaseDateStr[0])
		if err != nil {
			slog.Error("Failed to parse purchase date", "error", err)
			http.Error(w, "Invalid purchase date format", http.StatusBadRequest)
			return
		}
		stylus.PurchaseDate = &purchaseDate
	}

	styluses, err := s.controller.CreateStylus(&stylus)
	if err != nil {
		http.Error(w, "Failed to create stylus", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeData(w, styluses)
}

func (s *Server) updateStylus(w http.ResponseWriter, r *http.Request) {
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

	var stylus database.Stylus
	err = json.NewDecoder(r.Body).Decode(&stylus)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stylus.ID = id

	if purchaseDateStr, ok := r.URL.Query()["purchase_date"]; ok && len(purchaseDateStr) > 0 {
		purchaseDate, err := time.Parse(time.RFC3339, purchaseDateStr[0])
		if err != nil {
			slog.Error("Failed to parse purchase date", "error", err)
			http.Error(w, "Invalid purchase date format", http.StatusBadRequest)
			return
		}
		stylus.PurchaseDate = &purchaseDate
	}

	styluses, err := s.controller.UpdateStylus(&stylus)
	if err != nil {
		http.Error(w, "Failed to update stylus", http.StatusInternalServerError)
		return
	}

	writeData(w, styluses)
}

func (s *Server) deleteStylus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	styluses, err := s.controller.DeleteStylus(id)
	if err != nil {
		http.Error(w, "Failed to delete stylus", http.StatusInternalServerError)
		return
	}

	writeData(w, styluses)
}
