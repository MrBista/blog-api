package main

import (
	"github.com/MrBista/blog-api/internal/config"
	"github.com/MrBista/blog-api/internal/database"
	"github.com/MrBista/blog-api/internal/middleware"
	"github.com/MrBista/blog-api/internal/router"
	"github.com/MrBista/blog-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.LoadConfig()

	database.Connect()
	defer database.Close()

	utils.InitJwtService()

	utils.GetValidator()

	utils.InitLogger()

	utils.InitGoogleOAuth()

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.HandleError,
		Prefork:      true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*", // Ganti "*" dengan domain spesifik jika perlu
		AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length",
		// AllowCredentials: true,
	}))

	app.Static("/public", "./public")

	router.SetupAllRoutes(app)

	app.Listen(":3000")

}
