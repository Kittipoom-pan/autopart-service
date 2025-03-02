package config

import (
	"log"
	"os"
	"sync"

	utils "github.com/Kittipoom-pan/autopart-service/pkg/utils"
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

func GetConfig() *Config {
	once.Do(func() {
		// Load .env file
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
	})

	return configInstance
}
