package main

import (
	"log"

	"github.com/Novodremov/subscribe-be/config"
	"github.com/Novodremov/subscribe-be/internal/db"
	"github.com/Novodremov/subscribe-be/internal/db/migrator"
	"github.com/Novodremov/subscribe-be/internal/delivery/http"
	"github.com/Novodremov/subscribe-be/internal/logging"
	"github.com/Novodremov/subscribe-be/internal/repo"
	"github.com/Novodremov/subscribe-be/internal/service"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("fail to read .env-file")
	}
	fx.New(
		config.Module,
		logging.Module,
		db.Module,
		migrator.Module,
		http.Module,
		service.Module,
		repo.Module,
	).Run()
}
