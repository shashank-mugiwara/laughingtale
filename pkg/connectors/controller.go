package connectors

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shashank-mugiwara/laughingtale/client"
)

type Handler interface {
	AddConnectorConfig(ctx *fiber.Ctx) error
}

type Logger interface {
}

func RegisterRoutes(router *fiber.App, logger Logger) {
	h := NewHandler(logger, client.HttpClient())
	router.Post("/api/v1/connectors", h.AddConnectorConfig)
}
