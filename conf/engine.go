package conf

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	customLogger "github.com/shashank-mugiwara/laughingtale/logger"
)

var router *fiber.App

func GetLaughingTaleEngine() *fiber.App {
	return router
}

func InitEngine() {
	app := fiber.New()
	app.Use(logger.New(customLogger.GetFiberLoggerConfig()))
	router = app
}
