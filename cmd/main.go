package main

import (
	"github.com/MrBista/blog-api/internal/config"
	"github.com/MrBista/blog-api/internal/database"
	"github.com/MrBista/blog-api/internal/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadConfig()

	database.Connect()
	defer database.Close()

	app := fiber.New()

	router.SetupAllRoutes(app)

	app.Listen(":3000")

}
