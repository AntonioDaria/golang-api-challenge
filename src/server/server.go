package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// Run the server in a separate goroutine
	go func() {
		s.logger.Info().Msg("🚀 Starting HTTP Server")
		if err := s.app.Listen(":3000"); err != nil {
			s.logger.Fatal().Err(err).Msg("Server failure")
		}
	}()

	// Set up channel to listen for shutdown signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Block until a signal is received
	<-osSignals
	s.logger.Info().Msg("🔴 Shutting down HTTP Server")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		s.logger.Info().Msg("🔴 HTTP Server shutdown complete")
	}

	return nil
}
