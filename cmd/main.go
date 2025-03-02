package main

import (
	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/server"
)

func main() {
	conf := config.GetConfig()

	database.NewMySQLDatabase(conf)

	server.NewServer(conf).Start()
}
