package config

import (
	"errors"
	"log"
	"os"
	"sync"

	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	Server *Server
	Db     *Db
}

type Server struct {
	Host string
	Port int
}

type Db struct {
	MySqlHost     string
	MySqlPort     int
	MySqlDatabase string
	MySqlUser     string
	MySqlPassword string
}

var (
	once           sync.Once
	configInstance *Config
)

func LoadConfigs() (*Config, error) {
	var loadErr error

	once.Do(func() {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("Warning: No .env file found")
		}

		configInstance = &Config{
			Server: &Server{
				Host: os.Getenv("SERVER_HOST"),
				Port: utils.GetEnvAsInt("SERVER_PORT", 3000),
			},
			Db: &Db{
				MySqlHost:     os.Getenv("MYSQL_HOST"),
				MySqlPort:     utils.GetEnvAsInt("MYSQL_PORT", 3306),
				MySqlDatabase: os.Getenv("MYSQL_DATABASE"),
				MySqlUser:     os.Getenv("MYSQL_USER"),
				MySqlPassword: os.Getenv("MYSQL_PASSWORD"),
			},
		}

		if configInstance.Server.Host == "" || configInstance.Db.MySqlHost == "" {
			loadErr = errors.New("missing required configuration values")
		}
	})

	return configInstance, loadErr
}
