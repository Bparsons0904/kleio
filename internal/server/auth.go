package server

import (
	"database/sql"
	"encoding/json"
	"kleio/internal/handlers"
	"log"
	"log/slog"
	"net/http"
)

func (s *Server) GetAuth(w http.ResponseWriter, r *http.Request) {
	db := s.db.GetDB()

	var token string
	err := db.QueryRow("SELECT token FROM auth").Scan(&token)
	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			log.Printf("Database query error: %v", err)
		}
		return
	}

	go handlers.UpdateCollection(s.db)

	resp := map[string]string{"token": token}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
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

	username, err := handlers.GetUserIdentity(token)
	if err != nil {
		http.Error(w, "Failed to get user identity", http.StatusInternalServerError)
		slog.Error("Failed to get user identity", "error", err)
		return
	}

	slog.Info("Saving token...", "token", token)

	s.db.SaveToken(token, username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]bool{"success": true}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to write response", "error", err)
	}
}
