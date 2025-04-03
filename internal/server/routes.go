package server

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Auth and Collection routes
	mux.HandleFunc("/auth", s.getAuth)
	mux.HandleFunc("/auth/token", s.SaveToken)
	mux.HandleFunc("/collection", s.getCollection)
	mux.HandleFunc("/collection/sync", s.checkSync)
	mux.HandleFunc("/collection/resync", s.updateCollection)
	mux.HandleFunc("/discogs/collection/refresh", s.updateCollection)

	mux.HandleFunc("/styluses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.createStylus(w, r)
		default:
			slog.Warn("Method not allowed", "method", r.Method, "path", r.URL.Path)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/styluses/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/styluses/")
		if path == "" {
			http.Error(w, "Missing stylus ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodPut:
			s.updateStylus(w, r)
		case http.MethodDelete:
			s.deleteStylus(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Play history routes
	mux.HandleFunc("/plays", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.createPlayHistory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/plays/counts", s.getPlayCounts)
	mux.HandleFunc("/plays/recent", s.getRecentPlays)
	mux.HandleFunc("/plays/range", s.getPlaysByTimeRange)

	mux.HandleFunc("/plays/", func(w http.ResponseWriter, r *http.Request) {
		// Skip if it's one of the special routes
		path := r.URL.Path
		if path == "/plays/" ||
			strings.HasPrefix(path, "/plays/counts") ||
			strings.HasPrefix(path, "/plays/recent") ||
			strings.HasPrefix(path, "/plays/range") ||
			strings.HasPrefix(path, "/plays/release/") {
			return
		}

		// Handling routes with IDs
		switch r.Method {
		case http.MethodPut:
			s.updatePlayHistory(w, r)
		case http.MethodDelete:
			s.deletePlayHistory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Cleaning history routes
	mux.HandleFunc("/cleanings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.createCleaningHistory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/cleanings/counts", s.getCleaningCountsByRelease)
	mux.HandleFunc("/cleanings/range", s.getCleaningsByTimeRange)

	mux.HandleFunc("/cleanings/", func(w http.ResponseWriter, r *http.Request) {
		// Skip if it's one of the special routes
		path := r.URL.Path
		if path == "/cleanings/" ||
			strings.HasPrefix(path, "/cleanings/counts") ||
			strings.HasPrefix(path, "/cleanings/range") {
			return
		}

		// Handling routes with IDs
		switch r.Method {
		case http.MethodPut:
			s.updateCleaningHistory(w, r)
		case http.MethodDelete:
			s.deleteCleaningHistory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/export/history", s.exportHistory)

	// Setup static file server for SPA
	distDir := "./clio/dist"
	fileServer := http.FileServer(http.Dir(distDir))

	// Create a special handler for SPA routing
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// First check if the path is an API endpoint
		if strings.HasPrefix(r.URL.Path, "/auth") ||
			strings.HasPrefix(r.URL.Path, "/collection") ||
			strings.HasPrefix(r.URL.Path, "/styluses") ||
			strings.HasPrefix(r.URL.Path, "/plays") ||
			strings.HasPrefix(r.URL.Path, "/cleanings") ||
			strings.HasPrefix(r.URL.Path, "/export") {
			http.NotFound(w, r)
			return
		}

		// Check if this is a request for a static file
		path := filepath.Join(distDir, r.URL.Path)
		_, err := os.Stat(path)
		if err == nil {
			// File exists, serve it directly
			fileServer.ServeHTTP(w, r)
			return
		}

		// For all other paths, serve index.html to support client-side routing
		indexPath := filepath.Join(distDir, "index.html")
		http.ServeFile(w, r, indexPath)
	})

	return s.corsMiddleware(mux)
}
