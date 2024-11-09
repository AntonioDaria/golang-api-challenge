package router

import (
	"github.com/AntonioDaria/surfe/src/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func New() *fiber.App {
	router := fiber.New()

	// Add Recover middleware to handle panics
	router.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	router.Get("/health", handlers.HealthCheckHandler) // Register health check endpoint
	return router
}
