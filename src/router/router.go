package router

import (
	"github.com/AntonioDaria/surfe/src/handlers/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func New(userHandler *user.Handler) *fiber.App {
	router := fiber.New()

	// Add Recover middleware to handle panics
	router.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// User endpoint
	router.Get("/user/:id", userHandler.GetUserByIDHandler)
	return router
}
