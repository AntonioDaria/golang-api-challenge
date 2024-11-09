package router

import (
	"github.com/AntonioDaria/surfe/src/handlers"
	"github.com/gofiber/fiber/v2"
)

func New() *fiber.App {
	router := fiber.New()
	router.Get("/health", handlers.HealthCheckHandler) // Register health check endpoint
	return router
}
