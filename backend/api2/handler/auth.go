package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gityard-api/crud"
	"gityard-api/model"
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
		return c.Status(422).JSON(fiber.Map{"message": "invalid request"})
	}

	// validation
	err := validate.Struct(req)
	if err != nil {
		return c.Status(422).JSON(fiber.Map{"message": "invalid request"})
	}

	// email check
	dbInUser, err := crud.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "internal error"})
	}
	if dbInUser != nil {
		return c.Status(403).JSON(fiber.Map{"message": "registered email"})
	}

	// handle name check
	dbInHandleName, err := crud.GetHandleNameByName(req.HandleName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "internal error"})
	}
	if dbInHandleName != nil {
		return c.Status(403).JSON(fiber.Map{"message": "registered handlename"})
	}

	user, err := crud.CreateUser(req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "internal error"})
	}

	accountHandleName, err := crud.CreateHandleName(req.HandleName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "internal error"})
	}

	_, err = crud.CreateAccount(user.ID, accountHandleName.ID, model.PersonalAccount)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "internal error"})
	}

	type Response struct {
		Email      string `json:"email"`
		HandleName string `json:"handlename"`
	}
	res := new(Response)
	res.Email = *user.Email
	res.HandleName = accountHandleName.Handlename

	return c.JSON(res)
}

// Login handler for /login
func Login(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

// Logout handler for /logout
func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

func Refresh(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}
