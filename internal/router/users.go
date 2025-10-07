package router

import "github.com/gofiber/fiber/v2"

func SetUserRoute(router fiber.Router) {
	userRoute := router.Group("/users")

	userRoute.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello users")

	})
}
