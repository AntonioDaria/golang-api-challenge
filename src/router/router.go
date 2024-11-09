package router

import (
	"github.com/AntonioDaria/surfe/src/handlers/action"
	"github.com/AntonioDaria/surfe/src/handlers/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Handlers struct {
	UserHandler   *user.Handler
	ActionHandler *action.Handler
}

func New(handlers *Handlers) *fiber.App {
	router := fiber.New()

	// Add Recover middleware to handle panics
	router.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// User endpoint
	router.Get("/user/:id", handlers.UserHandler.GetUserByIDHandler)

	// Action endpoint
	router.Get("/users/:id/actions/count", handlers.ActionHandler.GetActionCountByUserIDHandler)
	return router
}
