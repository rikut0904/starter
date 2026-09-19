package main

import (
	"log"

	"github.com/rikut0904/starter/create/next-go/backend/internal/infrastructure/config"
	"github.com/rikut0904/starter/create/next-go/backend/internal/infrastructure/database"
	apphttp "github.com/rikut0904/starter/create/next-go/backend/internal/interface/http"
	"github.com/rikut0904/starter/create/next-go/backend/internal/usecase"
)

func main() {
	cfg := config.Load()
	db, err := database.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	if err := database.Ping(db); err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	handler := apphttp.NewRouter(usecase.NewHealth())
	log.Printf("listening on 0.0.0.0:%s", cfg.Port)
	log.Fatal(handler.Start("0.0.0.0:" + cfg.Port))
}
