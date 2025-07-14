package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"gityard-api/config"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func WithoutAuthInfoProtection(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Authorization header exists",
		})
	}
	refreshTokenCookie := c.Cookies("refresh_token", "")
	if refreshTokenCookie != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Authorization cookie exists",
		})
	}
	return c.Next()
}

func AuthInfoProtection(c *fiber.Ctx) error {
	// 1. "Authorization"ヘッダーを取得
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header is missing",
		})
	}

	// 2. "Bearer <token>"の形式かチェック
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token format",
		})
	}
	tokenString := parts[1] // [0] == "Bearer"

	// 3. トークンを検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムが期待通りかチェック (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
		}
		return []byte(config.Config("SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired token",
		})
	}

	// 4. トークンからユーザー情報を取得し、Contextに保存
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		c.Locals("user_id", claims["user_id"])
	}

	// 5. 次のミドルウェアまたはハンドラへ処理を移す
	return c.Next()
}
