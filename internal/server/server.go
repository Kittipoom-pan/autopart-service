package server

import (
	"fmt"
	"log"

	"github.com/Kittipoom-pan/autopart-service/config"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/gofiber/fiber/v2"
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
		log.Fatalln(err.Error())
		panic(err.Error())
	}

	fiberConnURL := fmt.Sprintf("%s:%d", s.Cfg.Server.Host, s.Cfg.Server.Port)

	log.Printf("Server has been started on %s", fiberConnURL)

	if err := s.App.Listen(fiberConnURL); err != nil {
		log.Fatalln(err.Error())
		panic(err.Error())
	}
}
