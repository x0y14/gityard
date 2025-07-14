package router

import (
	"gityard-api/handler"
	"gityard-api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")

	v1.Get("/healthcheck", handler.HealthCheck)

	auth := v1.Group("/auth")
	auth.Post("/signup", middleware.WithoutAuthInfoProtection, handler.SignUp)
	auth.Post("/login", middleware.WithoutAuthInfoProtection, handler.Login)
	auth.Post("/logout", middleware.AuthHeaderProtection, handler.Logout)
	auth.Post("/refresh", middleware.AuthCookieProtection, handler.Refresh)
}
