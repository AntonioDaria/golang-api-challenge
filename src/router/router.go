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

	// Action endpoints
	router.Get("/users/:id/actions/count", handlers.ActionHandler.GetActionCountByUserIDHandler)
	router.Get("/actions/:actionType/next", handlers.ActionHandler.GetNextActionProbabilitiesHandler)
	router.Get("/actions/referral", handlers.ActionHandler.GetReferralIndexHandler)

	return router
}
