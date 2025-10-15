package main

import (
	"github.com/MrBista/blog-api/internal/config"
	"github.com/MrBista/blog-api/internal/database"
	"github.com/MrBista/blog-api/internal/middleware"
	"github.com/MrBista/blog-api/internal/router"
	"github.com/MrBista/blog-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadConfig()

	database.Connect()
	defer database.Close()

	utils.InitJwtService()

	utils.GetValidator()

	utils.InitLogger()

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.HandleError,
		Prefork:      true,
	})

	app.Static("/public", "./public")

	router.SetupAllRoutes(app)

	app.Listen(":3000")

}
