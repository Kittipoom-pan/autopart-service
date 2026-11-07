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
	JWT    *JWT
	Upload *Upload
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

type JWT struct {
	Secret string
	Expiry int
}

type Upload struct {
	Dir          string
	BaseURL      string
	MaxSizeBytes int64
	MaxWidth     int
}

var (
	once           sync.Once
	configInstance *Config
)

func LoadConfigs() (*Config, error) {
	var loadErr error

	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error loading .env file")
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
			JWT: &JWT{
				Secret: os.Getenv("JWT_SECRET_KEY"),
				Expiry: utils.GetEnvAsInt("JWT_EXPIRY", 3600), // default 1 hours
			},
			Upload: &Upload{
				Dir:          utils.GetEnv("UPLOAD_DIR", "uploads"),
				BaseURL:      utils.GetEnv("UPLOAD_BASE_URL", "/uploads"),
				MaxSizeBytes: int64(utils.GetEnvAsInt("UPLOAD_MAX_SIZE_MB", 5)) * 1024 * 1024,
				MaxWidth:     utils.GetEnvAsInt("UPLOAD_MAX_WIDTH", 1200),
			},
		}

		if configInstance.Server.Host == "" || configInstance.Db.MySqlHost == "" {
			loadErr = errors.New("missing required configuration values")
		}
	})

	return configInstance, loadErr
}
