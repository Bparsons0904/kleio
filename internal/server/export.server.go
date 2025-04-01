// internal/server/export.server.go
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func (s *Server) exportHistory(w http.ResponseWriter, r *http.Request) {
	// Get export data from controller
	exportData, err := s.controller.ExportHistory()
	if err != nil {
		slog.Error("Failed to get export data", "error", err)
		http.Error(w, "Failed to export data", http.StatusInternalServerError)
		return
	}

	// Set headers for file download
	filename := "kleio_history_export_" + time.Now().Format("2006-01-02") + ".json"
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	// Write JSON response directly to the response writer
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ") // Pretty print the JSON
	if err := encoder.Encode(exportData); err != nil {
		slog.Error("Failed to encode export data", "error", err)
		http.Error(w, "Failed to export data", http.StatusInternalServerError)
		return
	}
}
