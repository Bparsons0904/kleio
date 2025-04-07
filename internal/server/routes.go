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

	// create api base path
	apiMux := http.NewServeMux()

	// Auth and Collection routes
	apiMux.HandleFunc("/auth", s.getAuth)
	apiMux.HandleFunc("/auth/token", s.SaveToken)
	apiMux.HandleFunc("/collection", s.getCollection)
	apiMux.HandleFunc("/collection/sync", s.checkSync)
	apiMux.HandleFunc("/collection/resync", s.updateCollection)
	apiMux.HandleFunc("/discogs/collection/refresh", s.updateCollection)

	apiMux.HandleFunc("/styluses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.createStylus(w, r)
		default:
			slog.Warn("Method not allowed", "method", r.Method, "path", r.URL.Path)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	apiMux.HandleFunc("/styluses/", func(w http.ResponseWriter, r *http.Request) {
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
	apiMux.HandleFunc("/plays", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.createPlayHistory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	apiMux.HandleFunc("/plays/counts", s.getPlayCounts)
	apiMux.HandleFunc("/plays/recent", s.getRecentPlays)
	apiMux.HandleFunc("/plays/range", s.getPlaysByTimeRange)

	apiMux.HandleFunc("/plays/", func(w http.ResponseWriter, r *http.Request) {
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
	apiMux.HandleFunc("/cleanings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.createCleaningHistory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	apiMux.HandleFunc("/cleanings/counts", s.getCleaningCountsByRelease)
	apiMux.HandleFunc("/cleanings/range", s.getCleaningsByTimeRange)

	apiMux.HandleFunc("/cleanings/", func(w http.ResponseWriter, r *http.Request) {
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

	apiMux.HandleFunc("/export/history", s.exportHistory)

	mux.Handle("/api/", http.StripPrefix("/api", s.corsMiddleware(apiMux)))
	// Setup static file server for SPA
	distDir := "./clio/dist"
	fileServer := http.FileServer(http.Dir(distDir))

	// Create a special handler for SPA routing
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
