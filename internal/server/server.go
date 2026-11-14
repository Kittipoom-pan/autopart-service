package server

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Kittipoom-pan/autopart-service/config"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	App    *fiber.App
	Cfg    *config.Config
	Db     *db.Queries
	Conn   *sql.DB
	Logger zerolog.Logger
}

func NewServer(
	cfg *config.Config,
	conn *sql.DB,
	queries *db.Queries,
) *Server {
	return &Server{
		App: fiber.New(fiber.Config{
			BodyLimit: uploadBodyLimit(cfg),
		}),
		Cfg:    cfg,
		Db:     queries,
		Conn:   conn,
		Logger: log.Logger,
	}
}

// uploadBodyLimit allows a full batch of uploads plus multipart overhead.
func uploadBodyLimit(cfg *config.Config) int {
	const overheadBytes = 1 * 1024 * 1024

	return int(cfg.Upload.MaxSizeBytes)*cfg.Upload.MaxFiles + overheadBytes
}

func (s *Server) Start() {
	if err := s.MapHandlers(); err != nil {
		log.Error().Err(err).Msg("Failed to map handlers")
		os.Exit(1)
	}

	fiberConnURL := fmt.Sprintf(
		"%s:%d",
		s.Cfg.Server.Host,
		s.Cfg.Server.Port,
	)

	log.Info().Str("url", fiberConnURL).Msg("Server has been started")

	if err := s.App.Listen(fiberConnURL); err != nil {
		log.Error().Err(err).Str("url", fiberConnURL).Msg("Failed to start server")
		os.Exit(1)
	}
}
