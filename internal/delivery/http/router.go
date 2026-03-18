package http

import (
	"errors"
	"time"

	_ "github.com/Novodremov/subscribe-be/assets/docs"
	"github.com/Novodremov/subscribe-be/config"
	"github.com/Novodremov/subscribe-be/internal/delivery/http/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rs/zerolog/log"
	fiberswag "github.com/swaggo/fiber-swagger"
)

const (
	apiPrefix = "/api/v1"
)

// NewRouter создаёт новый fiber-роутер.
//
// @title           Subscribe API
// @version         1.0
// @description     Сервис для работы с подписками.
// @BasePath        /api/v1
// @host            localhost:8080
func NewRouter(
	sh *handler.SubscriptionHandler,
	cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:           cfg.Http.ReadTimeout,
		WriteTimeout:          cfg.Http.WriteTimeout,
		IdleTimeout:           cfg.Http.IdleTimeout,
		BodyLimit:             cfg.Http.BodyLimit,
		ReadBufferSize:        cfg.Http.ReadBufferSize,
		AppName:               "user-client",
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler,
	})

	app.Use(RequestLoggerMiddleware())

	app.Use(limiter.New(limiter.Config{
		Max:        cfg.Http.RateLimitMaxRequests,
		Expiration: cfg.Http.RateLimitWindow,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			log.Warn().
				Str("method", c.Method()).
				Str("route", c.OriginalURL()).
				Str("ip", c.IP()).
				Msg("rate limit exceeded")

			return fiber.ErrTooManyRequests
		},
	}))

	api := app.Group(apiPrefix)

	if cfg.App.Mode != "prod" {
		initSwaggerRoutes(api)
	}

	initSubscriptionRoutes(api, sh, "/subscription")

	return app
}

func errorHandler(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	status := fiber.StatusInternalServerError
	var httpErr *handler.HTTPError
	var fiberErr *fiber.Error

	switch {
	case errors.As(err, &httpErr):
		status = httpErr.Code
	case errors.As(err, &fiberErr):
		status = fiberErr.Code
	}

	log.Error().
		Err(err).
		Str("method", c.Method()).
		Str("route", c.OriginalURL()).
		Int("status", status).
		Str("type", "request").
		Time("time", time.Now().UTC()).
		Send()

	if httpErr != nil {
		return c.Status(httpErr.Code).JSON(httpErr)
	}

	if fiberErr != nil {
		return fiber.DefaultErrorHandler(c, err)
	}

	return c.Status(fiber.StatusInternalServerError).SendString(handler.ErrMsgInternalServerError)
}

func initSubscriptionRoutes(r fiber.Router, sh *handler.SubscriptionHandler, prefix string) fiber.Router {
	subs := r.Group(prefix)
	subs.Get("/total-cost", sh.SubscriptionsTotalCost)
	subs.Post("/", sh.CreateSubscription)
	subs.Get("/:id", sh.GetSubscription)
	subs.Put("/:id", sh.UpdateSubscription)
	subs.Delete("/:id", sh.DeleteSubscription)
	subs.Get("/", sh.ListSubscriptions)
	return subs
}

func initSwaggerRoutes(r fiber.Router) {
	r.Static("/swagger", "./static/swagger")
	r.Get("/swagger/*", fiberswag.WrapHandler)
}
