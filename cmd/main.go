package main

import (
	"log"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/logger"
	"github.com/Kittipoom-pan/autopart-service/internal/server"
)

func main() {
	conf, err := config.LoadConfigs()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewMySQLDatabase(conf)
	if err != nil {
		log.Fatalf("Database connection failed: %v\n", err)
	}

	logger.InitLogger()
	server.NewServer(conf, db).Start()
}
