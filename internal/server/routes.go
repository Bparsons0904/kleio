package server

import (
	"net/http"
	"strings"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Auth and Collection routes
	mux.HandleFunc("/auth", s.getAuth)
	mux.HandleFunc("/auth/token", s.SaveToken)
	mux.HandleFunc("/collection/sync", s.checkSync)
	mux.HandleFunc("/discogs/collection", s.updateCollection)
	mux.HandleFunc("/discogs/collection/refresh", s.updateCollection)

	// Stylus routes
	mux.HandleFunc("/styluses/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/styluses/" {
			if r.Method == http.MethodPost {
				s.createStylus(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
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

	return s.corsMiddleware(mux)
}
