package router

import "github.com/gofiber/fiber/v2"

func SetupAllRoutes(app *fiber.App) {
	router := app.Group("/api")

	SetupPostRoute(router)
	// SetAuthRoute(router)
	// SetUserRoute(router)

}
