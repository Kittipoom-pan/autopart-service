package server

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App *fiber.App
	Cfg *config.Config
	Db  *sql.DB
}

func NewServer(cfg *config.Config, db *sql.DB) *Server {
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
