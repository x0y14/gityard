package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gityard-api/model"
	"gityard-api/repository"
	"gityard-api/security"
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

	// 重複確認
	dbInUser, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		slog.Error("failed to get user by email", "detail", err)
		return InternalError(c)
	}
	if dbInUser != nil {
		slog.Info("signup rejected", "reason", "registered email")
		return c.Status(403).JSON(fiber.Map{"message": "registered email"})
	}

	dbInHandleName, err := repository.GetHandleNameByName(req.HandleName)
	if err != nil {
		slog.Error("failed to get handlename", "detail", err)
		return InternalError(c)
	}
	if dbInHandleName != nil {
		slog.Info("signup rejected", "reason", "registered handlename")
		return c.Status(403).JSON(fiber.Map{"message": "registered handlename"})
	}

	// 登録処理
	user, err := repository.CreateUser(req.Email)
	if err != nil {
		slog.Error("failed to create user", "detail", err)
		return InternalError(c)
	}

	handleName, err := repository.CreateHandleName(req.HandleName)
	if err != nil {
		slog.Error("failed to create handlename", "detail", err)
		return InternalError(c)
	}

	account, err := repository.CreateAccount(user.ID, handleName.ID, model.PersonalAccount)
	if err != nil {
		slog.Error("failed to create account", "detail", err)
		return InternalError(c)
	}

	_, err = repository.CreateAccountProfile(account.ID, handleName.Handlename, false)
	if err != nil {
		slog.Error("failed to create profile", "detail", err)
		return InternalError(c)
	}

	// 認証情報作成
	_, err = repository.CreateUserCredential(user.ID, req.Password)
	if err != nil {
		slog.Error("failed to create credential", "detail", err)
		return InternalError(c)
	}

	// generate & set refresh token into cookie
	refreshToken, err := repository.CreateOrUpdateUserRefreshToken(user.ID)
	if err != nil {
		slog.Error("failed to create or update refresh token", "detail", err)
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

	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		slog.Error("failed to get user by email", "detail", err)
		return InternalError(c)
	}

	if user == nil {
		slog.Info("login rejected", "reason", "not registered email")
		return c.Status(401).JSON(fiber.Map{"message": "invalid credentials"})
	}

	credential, err := repository.GetUserCredentialById(user.ID)
	if err != nil {
		slog.Error("failed to get credential by userid", "detail", err)
		return InternalError(c)
	}

	if credential == nil {
		slog.Info("login rejected", "reason", "not registered credential")
		return c.Status(401).JSON(fiber.Map{"message": "invalid credentials"})
	}

	if security.VerifyPassword(req.Password, credential.HashedPassword) == false {
		slog.Info("login rejected", "reason", "password does not match")
		return c.Status(401).JSON(fiber.Map{"message": "invalid credentials"})
	}

	// generate & set refresh token into cookie
	refreshToken, err := repository.CreateOrUpdateUserRefreshToken(user.ID)
	if err != nil {
		slog.Error("failed to create/update refresh token", "detail", err)
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

func ClearCookies(c *fiber.Ctx, key ...string) {
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
	//c.ClearCookie("refresh_token")
	ClearCookies(c, "refresh_token") // ref: https://github.com/gofiber/fiber/issues/1127
	return c.Status(200).JSON(fiber.Map{})
}

func Refresh(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	// generate & set refresh token into cookie
	refreshToken, err := repository.CreateOrUpdateUserRefreshToken(userId)
	if err != nil {
		slog.Error("failed to create/update refresh token", "detail", err)
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
