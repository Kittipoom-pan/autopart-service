package server

import (
	"fmt"
	"log"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	App *fiber.App
	Cfg *config.Config
}

func NewServer(cfg *config.Config) Server {
	return &FiberServer{
		App: fiber.New(),
		Cfg: cfg,
	}
}

func (s *FiberServer) Start() {
	fiberConnURL := fmt.Sprintf("%s:%d", s.Cfg.Server.Host, s.Cfg.Server.Port)

	log.Printf("Server has been started on %s", fiberConnURL)

	if err := s.App.Listen(fiberConnURL); err != nil {
		log.Fatalln(err.Error())
		panic(err.Error())
	}
}
