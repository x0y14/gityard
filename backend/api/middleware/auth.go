package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gityard-api/security"
	"strings"
)

func WithoutAuthInfoProtection(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "authorization header exists",
		})
	}
	refreshToken := c.Cookies("refresh_token", "")
	if refreshToken != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "authorization cookie exists",
		})
	}
	return c.Next()
}

func AuthHeaderProtection(c *fiber.Ctx) error {
	// 1. "Authorization"ヘッダーを取得
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "authorization header is missing",
		})
	}

	// 2. "Bearer <token>"の形式かチェック
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token format",
		})
	}
	accessToken := parts[1] // [0] == "Bearer"

	// 3. トークンを検証
	userId, ok := security.VerifyAccessToken(accessToken)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid access_token"})
	}

	// 4. 情報取り出す
	c.Locals("user_id", userId)

	return c.Next()
}
