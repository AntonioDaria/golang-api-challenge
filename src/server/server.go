package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Server struct {
	app    *fiber.App
	logger zerolog.Logger
}

func New(logger zerolog.Logger, httpRouter *fiber.App) *Server {
	return &Server{
		app:    httpRouter,
		logger: logger,
	}
}

func (s *Server) Run() error {
	return s.app.Listen(":3000")
}
