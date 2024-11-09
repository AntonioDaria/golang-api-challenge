package main

import (
	"os"

	"github.com/AntonioDaria/surfe/src/handlers/action"
	"github.com/AntonioDaria/surfe/src/handlers/user"
	action_repo "github.com/AntonioDaria/surfe/src/repository/action"
	user_repo "github.com/AntonioDaria/surfe/src/repository/user"
	"github.com/AntonioDaria/surfe/src/router"
	"github.com/AntonioDaria/surfe/src/server"
	action_service "github.com/AntonioDaria/surfe/src/services/action"
	users_service "github.com/AntonioDaria/surfe/src/services/user"

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

	// Load Action JSON data
	actionRepo, err := action_repo.NewActionRepo("./src/repository/data/actions.json")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load action data")
	}

	// Initialize user service and handler
	userService := users_service.NewUserService(userRepo)
	userHandler := user.NewHandler(userService, logger)

	actionService := action_service.NewActionService(actionRepo)
	actionHandler := action.NewHandler(actionService, logger)

	// Group handlers
	handlers := &router.Handlers{
		UserHandler:   userHandler,
		ActionHandler: actionHandler,
	}

	// Initialize router
	httpRouter := router.New(handlers)

	// Set up server and run the server
	httpServer := server.New(logger, httpRouter)
	if err := httpServer.Run(); err != nil {
		logger.Fatal().Err(err).Msg("server failure")
	}
}
