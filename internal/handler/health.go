package handler

import (
	"github.com/gofiber/fiber/v2"
)

// HealthCheck handles health check requests
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"version": "1.0.0",
	})
}
