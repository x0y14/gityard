package handler

import "github.com/gofiber/fiber/v2"

// HealthCheck handler for /healthcheck
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}
