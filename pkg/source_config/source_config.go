package sourceconfig

import (
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	Logger Logger
}

func NewHandler(logger Logger) Handler {
	return &handler{
		Logger: logger,
	}
}

func (h handler) AddSourceConfig(c *fiber.Ctx) error {
	h.Logger.Info("Got request for adding source config with name: ", c.Params("name"))

	return nil
}
