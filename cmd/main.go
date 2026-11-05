package main

import (
	"log"

	"github.com/Kittipoom-pan/autopart-service/config"
	_ "github.com/Kittipoom-pan/autopart-service/docs"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/logger"
	"github.com/Kittipoom-pan/autopart-service/internal/server"
)

// @title        Autopart Service API
// @version      1.0
// @BasePath     /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description  "Bearer <token>"
func main() {
	conf, err := config.LoadConfigs()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	conn, queries, err := database.NewMySQLDatabase(conf)
	if err != nil {
		log.Fatalf("Database connection failed: %v\n", err)
	}

	logger.InitLogger()
	server.NewServer(conf, conn, queries).Start()
}
