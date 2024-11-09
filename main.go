package main

import (
	"os"

	"github.com/AntonioDaria/surfe/src/handlers/user"
	user_repo "github.com/AntonioDaria/surfe/src/repository/user"
	"github.com/AntonioDaria/surfe/src/router"
	"github.com/AntonioDaria/surfe/src/server"
	"github.com/AntonioDaria/surfe/src/services"

	"github.com/rs/zerolog"
)

func main() {
	// Set up logger
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	// Load User JSON data
	userRepo, err := user_repo.NewUserRepo("./src/repository/data/users.json")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load user data")
	}

	// Initialize user service and handler
	userService := services.NewUserService(userRepo)
	userHandler := user.NewHandler(userService, logger)

	// Initialize router
	httpRouter := router.New(userHandler)

	// Set up server and run the server
	httpServer := server.New(logger, httpRouter)
	if err := httpServer.Run(); err != nil {
		logger.Fatal().Err(err).Msg("server failure")
	}
}
