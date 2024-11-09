package handlers

import (
	"github.com/AntonioDaria/surfe/src/services"
	"github.com/gofiber/fiber/v2"
)

func HealthCheckHandler(c *fiber.Ctx) error {
	status := services.HealthCheck()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": status})
}
