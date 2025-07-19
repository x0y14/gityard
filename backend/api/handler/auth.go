package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gityard-api/model"
	"gityard-api/security"
	"gityard-api/service"
	"log/slog"
	"time"
)

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
		var registeredEmailErr *service.ErrRegisteredEmail
		if errors.As(err, &registeredEmailErr) {
			slog.Info("sign up rejected", "reason", "registered email")
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "registered email"})
		}

		var registeredHandleNameErr *service.ErrRegisteredHandleName
		if errors.As(err, &registeredHandleNameErr) {
			slog.Info("sign up rejected", "reason", "registered handlename")
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "registered handlename"})
		}

		slog.Error("failed to sign up", "detail", err)
		return InternalError(c)
	}

	slog.Info("user signed up successfully", "userId", user.ID, "handleName", req.HandleName)
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
		var userNotFoundErr *service.ErrUserNotFound
		if errors.As(err, &userNotFoundErr) {
			slog.Warn("user login rejected", "reason", "email not registered")
			return UnauthorizedError(c)
		}

		var credentialNotFoundErr *service.ErrCredentialNotFound
		if errors.As(err, &credentialNotFoundErr) {
			// ユーザは存在しているのにクレデンシャルデータが存在しない
			slog.Error("user login rejected", "reason", "credential not found")
			return InternalError(c)
		}

		var passwordMissMatchErr *service.ErrPasswordMissMatch
		if errors.As(err, &passwordMissMatchErr) {
			slog.Warn("user login rejected", "reason", "password miss match", "email", req.Email)
			return UnauthorizedError(c)
		}

		slog.Error("user failed to login", "detail", err)
		return InternalError(c)
	}

	slog.Info("user logged in successfully", "userId", user.ID)
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

	slog.Info("user logged out successfully", "userId", userId)
	return c.Status(200).JSON(fiber.Map{})
}

func Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token", "")
	if refreshToken == "" {
		return UnauthorizedError(c)
	}

	userId, newRefreshToken, err := service.Refresh(refreshToken)
	if err != nil {
		var invalidErr *service.ErrInvalidRefreshTokenProvided
		if errors.As(err, &invalidErr) {
			slog.Warn("invalid refresh_token provided")
			return UnauthorizedError(c)
		}
		var expiredErr *service.ErrExpiredRefreshTokenProvided
		if errors.As(err, &expiredErr) {
			slog.Warn("expired refresh_token provided")
			return UnauthorizedError(c)
		}

		slog.Error("failed to refresh token", "detail", err)
		return InternalError(c)
	}

	slog.Info("token refreshed successfully", "userId", *userId)
	return setTokensAndRespond(c, *userId, newRefreshToken)
}
