package sourceconfig

import (
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	AddSourceConfig(ctx *fiber.Ctx) error
}

type Logger interface {
	Info(args ...interface{})
}

func RegisterRoutes(router *fiber.App, logger Logger) {
	h := NewHandler(logger)
	router.Get("/api/v1/sourceConfig/:name", h.AddSourceConfig)
}
