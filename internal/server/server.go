package server

import (
	"fmt"
	"kleio/internal/controller"
	"kleio/internal/database"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
)

func (s *Server) setupApp() http.Handler {
	app := fiber.New()

	app.Use(cors.New())

	prometheus := NewMetrics()
	RegisterMetrics(app, prometheus)

	s.RegisterRoutes(app)

	return adaptor.FiberApp(app)
}

type Server struct {
	port       int
	DB         database.Database
	controller *controller.Controller
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
	slog.Info("Starting server...", "port", port)
	NewServer := &Server{
		port:       port,
		DB:         database.New(),
		controller: controller.InitNewController(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.setupApp(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	slog.Info("Server started")

	return server
}
