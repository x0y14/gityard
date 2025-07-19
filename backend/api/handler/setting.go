package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gityard-api/service"
	"log/slog"
)

func RegisterSSHPublicKey(c *fiber.Ctx) error {
	// cookieからuseridを取り出す
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		slog.Error("user_id not found in locals or is not uint")
		return InternalError(c)
	}

	type Request struct {
		KeyName           string `json:"name" validate:"required"`
		PublicKeyFullText string `json:"full_text" validate:"required"`
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

	pk, err := service.RegisterSSHPublicKey(userId, req.KeyName, req.PublicKeyFullText)
	if err != nil {
		var invalidPkErr *service.ErrInvalidPubkeyProvided
		if errors.As(err, &invalidPkErr) {
			slog.Info("register ssh pubkey rejected", "reason", "invalid pubkey")
			return BadRequestError(c)
		}

		var duplicatesPkErr *service.ErrDuplicatesPubkeyFingerprint
		if errors.As(err, &duplicatesPkErr) {
			slog.Warn("register ssh pubkey rejected", "reason", "duplicates fingerprint in db one")
			return ConflictError(c)
		}

		slog.Error("failed to register ssh pubkey", "detail", err)
		return InternalError(c)
	}
	if pk == nil {
		slog.Error("registered ssh pubkey, but return null")
		return InternalError(c)
	}

	type Response struct {
		Fingerprint string `json:"fingerprint"`
	}
	res := Response{Fingerprint: pk.Fingerprint}
	slog.Info("user register ssh pubkey successfully", "userId", userId)
	return c.JSON(res)
}

func GetSSHPublicKeys(c *fiber.Ctx) error {
	// cookieからを取り出す
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		slog.Error("user_id not found in locals or is not uint")
		return InternalError(c)
	}

	pubkeys, err := service.GetSSHPublicKeys(userId, 0, 30)
	if err != nil {
		slog.Error("failed to get ssh pubkeys", "detail", err)
		return InternalError(c)
	}

	type Item struct {
		Name        string `json:"name"`
		Fingerprint string `json:"fingerprint"`
	}
	type Response struct {
		Pubkeys []Item `json:"keys"`
	}
	res := Response{
		Pubkeys: []Item{},
	}
	for _, pk := range pubkeys {
		res.Pubkeys = append(res.Pubkeys, Item{
			Name:        pk.Name,
			Fingerprint: pk.Fingerprint,
		})
	}
	return c.JSON(res)
}

func DeleteSSHPubkeyByFingerprint(c *fiber.Ctx) error {
	// cookieからを取り出す
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		slog.Error("user_id not found in locals or is not uint")
		return InternalError(c)
	}

	type Request struct {
		Fingerprint string `json:"fingerprint"`
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

	err = service.DeleteSSHPublicKeyByFingerprint(userId, req.Fingerprint)
	if err != nil {
		var userNotFoundErr *service.ErrUserNotFound
		if errors.As(err, &userNotFoundErr) {
			slog.Error("delete ssh pubkey rejected", "reason", "user not found")
			return InternalError(c)
		}
		slog.Error("delete ssh pubkey rejected", "detail", err)
		return InternalError(c)
	}
	return c.Status(200).JSON(fiber.Map{})
}
