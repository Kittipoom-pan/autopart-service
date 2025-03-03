package main

import (
	"log"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/server"
)

func main() {
	conf := config.LoadConfigs()

	db, err := database.NewMySQLDatabase(conf)

	if err != nil {
		log.Fatalf("Database connection failed: %v\n", err)
	}

	server.NewServer(conf, db).Start()
}
