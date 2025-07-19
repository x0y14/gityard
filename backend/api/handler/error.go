package handler

import "github.com/gofiber/fiber/v2"

func InternalError(c *fiber.Ctx) error {
	return c.Status(500).JSON(fiber.Map{"message": "internal error"})
}

func UnauthorizedError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
}

func BadRequestError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "bad request"})
}

func ConflictError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "conflict"})
}
