package middleware

import (
	"strings"

	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddlware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return exception.NewBadRequestErr("Missing authorization header")
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			return exception.NewBadRequestErr("Invalid authorization header")
		}

		tokenString := parts[1]

		jwtService := utils.GetJwtService()
		claim, err := jwtService.VerifyToken(tokenString)

		if err != nil {
			return exception.NewBadRequestErr("invalid or expired token")
		}

		c.Locals("user", claim)
		c.Locals("userId", claim.UserId)
		c.Locals("role", claim.Role)

		return c.Next()
	}
}
