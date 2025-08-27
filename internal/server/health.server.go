package server

import (
	"net/http"
)

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	// Simple health check - could add database connectivity test if needed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}