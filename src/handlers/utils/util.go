package utils

import "github.com/gofiber/fiber/v2"

func JsonError(c *fiber.Ctx, statusCode int, message string) error {
	return c.SendStatus(statusCode)
}
