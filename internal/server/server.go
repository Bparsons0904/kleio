package server

import (
	"fmt"
	"kleio/internal/controller"
	"kleio/internal/database"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port       int
	DB         database.Database
	controller *controller.Controller
}

func NewServer() *http.Server {
	slog.Info("Starting server...", "port", 38080)
	NewServer := &Server{
		port:       38080,
		DB:         database.New(),
		controller: controller.InitNewController(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	slog.Info("Server started")

	return server
}
