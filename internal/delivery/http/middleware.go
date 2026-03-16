package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// RequestLoggerMiddleware логирует успешные HTTP-запросы
func RequestLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		if err == nil {
			log.Info().
				Str("method", c.Method()).
				Str("route", c.OriginalURL()).
				Int("status", c.Response().StatusCode()).
				Dur("duration", time.Since(start)).
				Send()
		}
		// ошибки логируются в errorHandler
		return err
	}
}
