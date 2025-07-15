package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gityard-api/security"
	"gityard-api/service"
	"log/slog"
	"time"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// SignUp handler for /signup
func SignUp(c *fiber.Ctx) error {
	type Request struct {
		Email      string `json:"email" validate:"required,email"`
		Password   string `json:"password" validate:"required,min=8"`
		HandleName string `json:"handlename" validate:"required,alphanum"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		slog.Debug("failed to parse", "request body", req)
		return c.Status(422).JSON(fiber.Map{"message": "invalid request"})
	}

	// validation
	err := validate.Struct(req)
	if err != nil {
		slog.Debug("failed to validate", "request body", req)
		return c.Status(422).JSON(fiber.Map{"message": "invalid request"})
	}

	user, refreshToken, err := service.SignUp(req.Email, req.Password, req.HandleName)
	if err != nil {
		slog.Error("failed to sign up", "detail", err)
		return InternalError(c)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.RefreshToken,
		Expires:  refreshToken.ExpiresAt,
		Secure:   false,
		HTTPOnly: true,
		SameSite: "strict",
	})

	accessToken, err := security.GenerateAccessToken(user.ID)
	if err != nil {
		slog.Error("failed to generate access token", "detail", err)
		return InternalError(c)
	}

	type Response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"` // sec
	}
	res := new(Response)
	res.AccessToken = accessToken.Body
	res.TokenType = accessToken.Type
	res.ExpiresIn = int64(accessToken.ExpiresIn.Seconds())
	return c.JSON(*res)
}

// Login handler for /login
func Login(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		slog.Debug("failed to parse", "request body", req)
		return c.Status(422).JSON(fiber.Map{"message": "invalid request"})
	}
	// validation
	err := validate.Struct(req)
	if err != nil {
		slog.Debug("failed to validate", "request body", req)
		return c.Status(422).JSON(fiber.Map{"message": "invalid request"})
	}

	user, refreshToken, err := service.Login(req.Email, req.Password)
	if err != nil {
		slog.Error("failed to login", "detail", err)
		return InternalError(c)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.RefreshToken,
		Expires:  refreshToken.ExpiresAt,
		Secure:   false,
		HTTPOnly: true,
		SameSite: "strict",
	})

	accessToken, err := security.GenerateAccessToken(user.ID)
	if err != nil {
		slog.Error("failed to generate access token", "detail", err)
		return InternalError(c)
	}

	type Response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"` // sec
	}
	res := new(Response)
	res.AccessToken = accessToken.Body
	res.TokenType = accessToken.Type
	res.ExpiresIn = int64(accessToken.ExpiresIn.Seconds())
	return c.JSON(*res)
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
	userId := c.Locals("user_id").(uint)

	// 失効処理に失敗しようがしまいが、さらなる漏洩等を防ぐため消す
	clearCookies(c, "refresh_token")

	err := service.Logout(userId)
	if err != nil {
		// トークン漏洩の検出のため
		var revokedErr *service.ErrRevokedRefreshTokenProvided
		if errors.As(err, &revokedErr) {
			slog.Warn("revoked refresh_token provided", "user_id", userId)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "revoked refresh_token provided"})
		}

		slog.Error("failed to logout", "detail", err)
		return InternalError(c)
	}

	return c.Status(200).JSON(fiber.Map{})
}

func Refresh(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)
	refreshToken := c.Locals("refresh_token").(string)

	newRefreshToken, err := service.Refresh(userId, refreshToken)
	if err != nil {
		var revokedErr *service.ErrRevokedRefreshTokenProvided
		if errors.As(err, &revokedErr) {
			slog.Warn("revoked refresh_token provided", "user_id", userId)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "revoked refresh_token provided"})
		}

		slog.Error("failed to refresh token", "detail", err)
		return InternalError(c)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken.RefreshToken,
		Expires:  newRefreshToken.ExpiresAt,
		Secure:   false,
		HTTPOnly: true,
		SameSite: "strict",
	})

	accessToken, err := security.GenerateAccessToken(userId)
	if err != nil {
		slog.Error("failed to generate access token", "detail", err)
		return InternalError(c)
	}

	type Response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"` // sec
	}
	res := new(Response)
	res.AccessToken = accessToken.Body
	res.TokenType = accessToken.Type
	res.ExpiresIn = int64(accessToken.ExpiresIn.Seconds())
	return c.JSON(*res)
}
