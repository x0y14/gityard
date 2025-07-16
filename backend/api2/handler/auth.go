package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gityard-api/model"
	"gityard-api/security"
	"gityard-api/service"
	"log/slog"
	"time"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func setTokensAndRespond(c *fiber.Ctx, userId uint, refreshToken *model.RefreshToken) error {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Body,
		MaxAge:   int(refreshToken.ExpiresIn.Seconds()),
		Secure:   false, // 本番環境ではtrueにすべき
		HTTPOnly: true,
		SameSite: "strict",
		Path:     "/api/v1/auth/refresh",
	})

	accessToken, err := security.GenerateAccessToken(userId)
	if err != nil {
		slog.Error("failed to generate access token", "detail", err)
		return InternalError(c)
	}

	type Response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	res := Response{
		AccessToken: accessToken.Body,
		TokenType:   accessToken.Type,
		ExpiresIn:   int64(accessToken.ExpiresIn.Seconds()),
	}
	return c.JSON(res)
}

// SignUp handler for /signup
func SignUp(c *fiber.Ctx) error {
	type Request struct {
		Email      string `json:"email" validate:"required,email"`
		Password   string `json:"password" validate:"required,min=8"`
		HandleName string `json:"handlename" validate:"required,alphanum"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		slog.Debug("failed to parse", "detail", err)
		return c.Status(422).JSON(fiber.Map{"message": "invalid request"})
	}

	// validation
	err := validate.Struct(req)
	if err != nil {
		slog.Debug("failed to validate", "detail", err)
		return c.Status(422).JSON(fiber.Map{"message": "invalid request"})
	}

	user, refreshToken, err := service.SignUp(req.Email, req.Password, req.HandleName)
	if err != nil {
		slog.Error("failed to sign up", "detail", err)
		return InternalError(c)
	}

	return setTokensAndRespond(c, user.ID, refreshToken)
}

// Login handler for /login
func Login(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		slog.Debug("failed to parse", "detail", err)
		return c.Status(422).JSON(fiber.Map{"message": "invalid request"})
	}
	// validation
	err := validate.Struct(req)
	if err != nil {
		slog.Debug("failed to validate", "detail", err)
		return c.Status(422).JSON(fiber.Map{"message": "invalid request"})
	}

	user, refreshToken, err := service.Login(req.Email, req.Password)
	if err != nil {
		slog.Error("failed to login", "detail", err)
		return InternalError(c)
	}

	return setTokensAndRespond(c, user.ID, refreshToken)
}

// ref: https://github.com/gofiber/fiber/issues/1127
func clearCookies(c *fiber.Ctx, key ...string) {
	for i := range key {
		c.Cookie(&fiber.Cookie{
			Name:    key[i],
			Expires: time.Now().Add(-time.Hour * 24),
			Value:   "",
		})
	}
}

// Logout handler for /logout
func Logout(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		slog.Error("user_id not found in locals or is not uint")
		return InternalError(c)
	}

	err := service.Logout(userId)
	if err != nil {
		slog.Error("failed to logout", "detail", err)
		return InternalError(c)
	}

	return c.Status(200).JSON(fiber.Map{})
}

func Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token", "")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "authorization cookie is missing"})
	}

	userId, newRefreshToken, err := service.Refresh(refreshToken)
	if err != nil {
		var invalidErr *service.ErrInvalidRefreshTokenProvided
		if errors.As(err, &invalidErr) {
			slog.Warn("invalid refresh_token provided")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid refresh_token provided"})
		}

		slog.Error("failed to refresh token", "detail", err)
		return InternalError(c)
	}

	return setTokensAndRespond(c, *userId, newRefreshToken)
}
