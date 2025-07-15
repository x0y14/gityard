package middleware

import (
	"gityard-api/handler"
	"gityard-api/repository"
	"gityard-api/security"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func WithoutAuthInfoProtection(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "authorization header exists",
		})
	}
	refreshTokenCookie := c.Cookies("refresh_token", "")
	if refreshTokenCookie != "" {
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
	tokenString := parts[1] // [0] == "Bearer"

	// 3. トークンを検証
	userId, ok := security.VerifyAccessToken(tokenString)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid access_token"})
	}

	// 4. 情報取り出す
	c.Locals("user_id", userId)

	return c.Next()
}

func AuthCookieProtection(c *fiber.Ctx) error {
	// 1. "Authorization"クッキーを取得
	authCookie := c.Cookies("refresh_token", "")
	if authCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "authorization cookie is missing"})
	}

	userId, ok := security.VerifyRefreshToken(authCookie)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid refresh_token"})
	}

	refreshToken, err := repository.GetUserRefreshTokenById(userId)
	if err != nil {
		return handler.InternalError(c)
	}

	// 失効チェック
	if refreshToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid refresh_token"})
	}
	if refreshToken.RefreshToken != authCookie {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid refresh_token"})
	}

	// 4. 情報取り出す
	c.Locals("user_id", userId)

	return c.Next()
}
