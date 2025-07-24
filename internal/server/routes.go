package server

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func (s *Server) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth and Collection routes
	api.Get("/auth", adaptor.HTTPHandlerFunc(s.getAuth))
	api.Post("/auth/token", adaptor.HTTPHandlerFunc(s.SaveToken))
	api.Get("/collection", adaptor.HTTPHandlerFunc(s.getCollection))
	api.Get("/collection/sync", adaptor.HTTPHandlerFunc(s.checkSync))
	api.Post("/collection/resync", adaptor.HTTPHandlerFunc(s.updateCollection))
	api.Post("/discogs/collection/refresh", adaptor.HTTPHandlerFunc(s.updateCollection))
	api.Delete("/releases/:id/delete", adaptor.HTTPHandlerFunc(s.deleteRelease))
	api.Post("/releases/:id/archive", adaptor.HTTPHandlerFunc(s.archiveRelease))

	api.Post("/styluses", adaptor.HTTPHandlerFunc(s.createStylus))
	api.Put("/styluses/:id", adaptor.HTTPHandlerFunc(s.updateStylus))
	api.Delete("/styluses/:id", adaptor.HTTPHandlerFunc(s.deleteStylus))

	// Play history routes
	api.Post("/plays", adaptor.HTTPHandlerFunc(s.createPlayHistory))
	api.Get("/plays/counts", adaptor.HTTPHandlerFunc(s.getPlayCounts))
	api.Get("/plays/recent", adaptor.HTTPHandlerFunc(s.getRecentPlays))
	api.Get("/plays/range", adaptor.HTTPHandlerFunc(s.getPlaysByTimeRange))
	api.Put("/plays/:id", adaptor.HTTPHandlerFunc(s.updatePlayHistory))
	api.Delete("/plays/:id", adaptor.HTTPHandlerFunc(s.deletePlayHistory))

	// Cleaning history routes
	api.Post("/cleanings", adaptor.HTTPHandlerFunc(s.createCleaningHistory))
	api.Get("/cleanings/counts", adaptor.HTTPHandlerFunc(s.getCleaningCountsByRelease))
	api.Get("/cleanings/range", adaptor.HTTPHandlerFunc(s.getCleaningsByTimeRange))
	api.Put("/cleanings/:id", adaptor.HTTPHandlerFunc(s.updateCleaningHistory))
	api.Delete("/cleanings/:id", adaptor.HTTPHandlerFunc(s.deleteCleaningHistory))

	api.Get("/export/history", adaptor.HTTPHandlerFunc(s.exportHistory))

	// Setup static file server for SPA
	distDir := "./clio/dist"
	app.Static("/", distDir)

	// For all other paths, serve index.html to support client-side routing
	app.Get("/*", func(c *fiber.Ctx) error {
		// Check if this is a request for a static file
		path := filepath.Join(distDir, c.Path())
		_, err := os.Stat(path)
		if err == nil {
			// File exists, serve it directly
			return c.SendFile(path)
		}

		// For all other paths, serve index.html to support client-side routing
		indexPath := filepath.Join(distDir, "index.html")
		return c.SendFile(indexPath)
	})
}
