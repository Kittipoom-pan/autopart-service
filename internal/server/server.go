package server

import (
	"fmt"
	"os"

	"github.com/Kittipoom-pan/autopart-service/config"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Server struct {
	App *fiber.App
	Cfg *config.Config
	Db  *db.Queries
}

func NewServer(cfg *config.Config, db *db.Queries) *Server {
	return &Server{
		App: fiber.New(),
		Cfg: cfg,
		Db:  db,
	}
}

func (s *Server) Start() {
	if err := s.MapHandlers(); err != nil {
		log.Error().Err(err).Msg("Failed to map handlers")
		os.Exit(1)
	}

	fiberConnURL := fmt.Sprintf("%s:%d", s.Cfg.Server.Host, s.Cfg.Server.Port)

	log.Printf("Server has been started on %s", fiberConnURL)

	if err := s.App.Listen(fiberConnURL); err != nil {
		log.Error().Err(err).Str("url", fiberConnURL).Msg("Failed to start server")
		os.Exit(1)
	}
}
