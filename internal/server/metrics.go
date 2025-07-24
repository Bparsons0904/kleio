
package server

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

func NewMetrics() *fiberprometheus.FiberPrometheus {
	prometheus := fiberprometheus.New("kleio")
	return prometheus
}

func RegisterMetrics(app *fiber.App, prometheus *fiberprometheus.FiberPrometheus) {
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
}
