package migrator

import (
	"github.com/Novodremov/subscribe-be/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

// NewMigrator создает новый экземпляр Migrator с подключением к БД и конфигурацией.
// Используется для запуска миграций базы данных.
func NewMigrator(db *pgxpool.Pool, cfg *config.Config) *Migrator {
	return &Migrator{
		db:  db,
		cfg: cfg,
	}
}

type Migrator struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

// Run выполняет миграции базы данных, используя конфигурацию Migrator.
// Логирует начало процесса и возможные ошибки при применении миграций.
func (m *Migrator) Run() error {
	log.Info().Msg("starting database migrations")
	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		log.Error().Err(err).Msg("failed to set goose dialect")
		return err
	}

	db := stdlib.OpenDBFromPool(m.db)
	if err := goose.Up(db, m.cfg.Database.MigrationsDir); err != nil {
		log.Error().Err(err).Msg("failed to up migrations")
		return err
	}
	if err := db.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close connection")
		return err
	}

	return nil
}
