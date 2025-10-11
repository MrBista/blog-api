package router

import (
	"github.com/MrBista/blog-api/internal/database"
	"github.com/gofiber/fiber/v2"
)

func SetupAllRoutes(app *fiber.App) {
	router := app.Group("/api")

	SetupPostRoute(router, database.DB)
	SetAuthRoute(router, database.DB)
	SetupCategoryRouter(router, database.DB)
	// SetUserRoute(router)

}
