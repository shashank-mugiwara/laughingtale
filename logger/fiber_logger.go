package logger

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

func GetFiberLoggerConfig() logger.Config {
	return logger.Config{Next: nil,
		Done:          nil,
		Format:        "[${time}] ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat:    "15:04:05",
		TimeZone:      "Local",
		TimeInterval:  500 * time.Millisecond,
		Output:        os.Stdout,
		DisableColors: false}
}
