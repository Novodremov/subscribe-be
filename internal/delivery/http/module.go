package http

import (
	"context"

	"github.com/Novodremov/subscribe-be/config"
	"github.com/Novodremov/subscribe-be/internal/delivery/http/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	fx.Provide(
		handler.NewSubscriptionHandler,
		NewRouter,
	),
	fx.Invoke(
		RunServer,
	),
)

func RunServer(lc fx.Lifecycle, app *fiber.App, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := app.Listen(":" + cfg.Http.Port)
				if err != nil {
					log.Fatal().Err(err).Msg("fail to start http server")
				}
			}()
			log.Info().Str("port", cfg.Http.Port).Msg("http server started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("http server stopped")
			return app.ShutdownWithContext(ctx)
		},
	})
}
