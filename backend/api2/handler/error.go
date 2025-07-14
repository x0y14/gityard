package handler

import "github.com/gofiber/fiber/v2"

func InternalError(c *fiber.Ctx) error {
	return c.Status(500).JSON(fiber.Map{"message": "internal error"})
}
