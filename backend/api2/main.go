package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gityard-api/database"
	"gityard-api/router"
	"log"
)

func main() {
	//logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := fiber.New()
	//app.Use(slogfiber.New(logger))
	app.Use(cors.New())

	database.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":8000"))
}
