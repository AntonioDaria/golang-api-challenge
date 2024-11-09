package main

import (
	"os"

	"github.com/AntonioDaria/surfe/src/router"
	"github.com/AntonioDaria/surfe/src/server"

	"github.com/rs/zerolog"
)

func main() {
	// Set up logger
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	// Initialize router
	httpRouter := router.New()

	// Set up server
	httpServer := server.New(logger, httpRouter)

	// Run server
	if err := httpServer.Run(); err != nil {
		logger.Fatal().Err(err).Msg("server failure")
	}
}
