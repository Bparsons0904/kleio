package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.HelloWorldHandler)

	mux.HandleFunc("/health", s.GetAuth)
	mux.HandleFunc("/auth", s.GetAuth)
	mux.HandleFunc("/discogs/token", s.SaveToken)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().
			Set("Access-Control-Allow-Origin", "*")
			// Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().
			Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().
			Set("Access-Control-Allow-Credentials", "false")
			// Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) SaveToken(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var requestBody struct {
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		slog.Error("Failed to decode request body", "error", err)
		return
	}

	if requestBody.Token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	slog.Info("Saving token...", "token", requestBody.Token)

	// Get database connection
	db := s.db.GetDB()

	// Check if a token already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM auth").Scan(&count)
	if err != nil {
		http.Error(w, "Failed to check existing token", http.StatusInternalServerError)
		slog.Error("Failed to check existing token", "error", err)
		return
	}

	var sqlQuery string
	if count > 0 {
		// Update existing token
		sqlQuery = "UPDATE auth SET token = ?"
	} else {
		// Insert new token
		sqlQuery = "INSERT INTO auth (token) VALUES (?)"
	}

	// Execute the query
	_, err = db.Exec(sqlQuery, requestBody.Token)
	if err != nil {
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		slog.Error("Failed to save token", "error", err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]bool{"success": true}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to write response", "error", err)
	}
}

func (s *Server) GetAuth(w http.ResponseWriter, r *http.Request) {
	db := s.db.GetDB()

	// Use QueryRow for queries that return a single row
	var token string
	err := db.QueryRow("SELECT token FROM auth").Scan(&token)
	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			log.Printf("Database query error: %v", err)
		}
		return
	}

	// Create and send the response
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

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
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

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
