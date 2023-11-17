package cacheconfig

import "github.com/gofiber/fiber/v2"

type handler struct {
	Logger Logger
}

func NewHandler(logger Logger) Handler {
	return &handler{
		Logger: logger,
	}
}

func (h handler) AddSourceConfig(c *fiber.Ctx) error {
	return nil
}
