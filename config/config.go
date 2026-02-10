package config

import (
	"errors"
	"reflect"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	App  App  `env-prefix:"APP_"`
	Http Http `env-prefix:"HTTP_"`
	Database Database `env-prefix:"DB_"`
}

type App struct {
	Mode     string `env:"MODE" env-default:"local"`
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
}

type Http struct {
	Port           string        `env:"PORT" env-default:"8081"`
	ReadTimeout    time.Duration `env:"READ_TIMEOUT" env-default:"10s"`
	WriteTimeout   time.Duration `env:"WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout    time.Duration `env:"IDLE_TIMEOUT" env-default:"30s"`
	BodyLimit      int           `env:"BODY_LIMIT" env-default:"1048576"`
	ReadBufferSize int           `env:"READ_BUFFER_SIZE" env-default:"4096"`
}

type Database struct {
	Host            string        `env:"HOST" env-default:"localhost"`
	Port            string        `env:"PORT" env-default:"5432"`
	User            string        `env:"USER" env-default:"postgres"`
	Password        string        `env:"PASSWORD" env-default:"postgres"`
	Name            string        `env:"NAME" env-default:"subs"`
	SSLMode         string        `env:"SSL_MODE" env-default:"disable"`

	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" env-default:"10"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" env-default:"1h"`

	MigrationsDir   string        `env:"MIGRATIONS_DIR" env-default:"migrations"`
}


func (c Config) IsNil() bool {
	return reflect.DeepEqual(c, Config{})
}

func (c Config) SetLogLevel() {
	switch c.App.LogLevel {
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func New() (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(".env", cfg)

	if err != nil || cfg.IsNil() {
		log.Info().Err(err).Msg("read config from .env error, attempt to read config from environment")
		err = cleanenv.ReadEnv(cfg)
		if err != nil {
			log.Err(err).Msg("read config from environment error")
			return nil, err
		}
	}

	if cfg.IsNil() {
		err = errors.New("config is nil")
		log.Err(err).Msg("read config from .env and environment error")
		return nil, err
	}

	cfg.SetLogLevel()

	log.Debug().Any("config", cfg).Msg("app configuration")
	return cfg, nil
}
