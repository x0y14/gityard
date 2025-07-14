package router

import (
	"gityard-api/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")

	v1.Get("/healthcheck", handler.HealthCheck)

	auth := v1.Group("/auth")
	auth.Post("/signup", handler.SignUp)
	auth.Post("/login", handler.Login)
}
