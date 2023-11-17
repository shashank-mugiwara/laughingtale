package sourceconfig

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shashank-mugiwara/laughingtale/db"
)

type Handler interface {
	AddSourceConfig(ctx *fiber.Ctx) error
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
}

func RegisterRoutes(router *fiber.App, logger Logger) {
	h := NewHandler(logger, db.GetlaughingtaleDb())
	router.Post("/api/v1/loaderSourceConfig", h.AddSourceConfig)
}
