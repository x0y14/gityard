package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gityard-api/crud"
	"gityard-api/model"
	"gityard-api/secutiry"
	"log/slog"
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
	dbInUser, err := crud.GetUserByEmail(req.Email)
	if err != nil {
		slog.Error("failed to get user by email", "detail", err)
		return InternalError(c)
	}
	if dbInUser != nil {
		slog.Info("signup rejected", "reason", "registered email")
		return c.Status(403).JSON(fiber.Map{"message": "registered email"})
	}

	dbInHandleName, err := crud.GetHandleNameByName(req.HandleName)
	if err != nil {
		slog.Error("failed to get handlename", "detail", err)
		return InternalError(c)
	}
	if dbInHandleName != nil {
		slog.Info("signup rejected", "reason", "registered handlename")
		return c.Status(403).JSON(fiber.Map{"message": "registered handlename"})
	}

	// 登録処理
	user, err := crud.CreateUser(req.Email)
	if err != nil {
		slog.Error("failed to create user", "detail", err)
		return InternalError(c)
	}

	handleName, err := crud.CreateHandleName(req.HandleName)
	if err != nil {
		slog.Error("failed to create handlename", "detail", err)
		return InternalError(c)
	}

	account, err := crud.CreateAccount(user.ID, handleName.ID, model.PersonalAccount)
	if err != nil {
		slog.Error("failed to create account", "detail", err)
		return InternalError(c)
	}

	_, err = crud.CreateAccountProfile(account.ID, handleName.Handlename, false)
	if err != nil {
		slog.Error("failed to create profile", "detail", err)
		return InternalError(c)
	}

	// 認証情報作成
	_, err = crud.CreateUserCredential(user.ID, req.Password)
	if err != nil {
		slog.Error("failed to create credential", "detail", err)
		return InternalError(c)
	}

	// generate & set refresh token into cookie
	refreshToken, err := crud.CreateOrUpdateUserRefreshToken(user.ID)
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

	accessToken, err := secutiry.GenerateAccessToken(user.ID)
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

	user, err := crud.GetUserByEmail(req.Email)
	if err != nil {
		slog.Error("failed to get user by email", "detail", err)
		return InternalError(c)
	}

	if user == nil {
		slog.Info("login rejected", "reason", "not registered email")
		return c.Status(401).JSON(fiber.Map{"message": "invalid credentials"})
	}

	credential, err := crud.GetUserCredentialById(user.ID)
	if err != nil {
		slog.Error("failed to get credential by userid", "detail", err)
		return InternalError(c)
	}

	if credential == nil {
		slog.Info("login rejected", "reason", "not registered credential")
		return c.Status(401).JSON(fiber.Map{"message": "invalid credentials"})
	}

	if secutiry.VerifyPassword(req.Password, credential.HashedPassword) == false {
		slog.Info("login rejected", "reason", "password does not match")
		return c.Status(401).JSON(fiber.Map{"message": "invalid credentials"})
	}

	// generate & set refresh token into cookie
	refreshToken, err := crud.CreateOrUpdateUserRefreshToken(user.ID)
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

	accessToken, err := secutiry.GenerateAccessToken(user.ID)
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

// Logout handler for /logout
func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

func Refresh(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}
