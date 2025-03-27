package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/auth", s.getAuth)
	mux.HandleFunc("/auth/token", s.SaveToken)
	mux.HandleFunc("/collection/sync", s.checkSync)
	mux.HandleFunc("/discogs/collection", s.updateCollection)
	mux.HandleFunc("/discogs/collection/refresh", s.updateCollection)

	return s.corsMiddleware(mux)
}
