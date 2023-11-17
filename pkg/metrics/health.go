package metrics

import "github.com/gofiber/fiber/v2"

type handler struct {
	Logger Logger
}

func NewHandler(logger Logger) Handler {
	return handler{
		Logger: logger,
	}
}

func (h handler) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "healthy",
	})
}
