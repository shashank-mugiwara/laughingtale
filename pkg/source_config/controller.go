package sourceconfig

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shashank-mugiwara/laughingtale/db"
)

type Handler interface {
	AddLoaderSourceConfig(ctx *fiber.Ctx) error
	GetLoaderSourceConfig(ctx *fiber.Ctx) error
	DeleteLoaderSourceConfig(ctx *fiber.Ctx) error
	UpdateLoaderSourceConfig(ctx *fiber.Ctx) error
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
}

func RegisterRoutes(router *fiber.App, logger Logger) {
	h := NewHandler(logger, db.GetlaughingtaleDb())
	router.Post("/api/v1/loaderSourceConfig", h.AddLoaderSourceConfig)
	router.Get("/api/v1/loaderSourceConfig/:name", h.GetLoaderSourceConfig)
	router.Delete("/api/v1/loaderSourceConfig/:name", h.DeleteLoaderSourceConfig)
	router.Put("/api/v1/loaderSourceConfig/:name", h.UpdateLoaderSourceConfig)
}
