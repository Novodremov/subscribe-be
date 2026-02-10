package logging

import (
	"os"

	"github.com/Novodremov/subscribe-be/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"logging",
	fx.Provide(
		NewLogger,
	),
)

func NewLogger(cfg *config.Config) zerolog.Logger {
	level := zerolog.InfoLevel
	switch cfg.App.LogLevel {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	}

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "subscribe-be").
		Logger().
		Level(level)

	log.Logger = logger

	return logger
}
