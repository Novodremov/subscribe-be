package db

import (
	"context"
	"fmt"

	"github.com/Novodremov/subscribe-be/config"

	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(lc fx.Lifecycle, cfg *config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse Postgres DSN: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Postgres pool: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := pool.Ping(ctx); err != nil {
				return fmt.Errorf("failed to connect to postgres: %w", err)
			}
			log.Info().Msg("successfully connected to postgres")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("closing connection to postgres")
			pool.Close()
			return nil
		},
	})

	return pool, nil
}
