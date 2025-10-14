package middleware

import (
	"strings"

	"github.com/MrBista/blog-api/internal/enum"
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

func RoleMiddleare(allowedRoles ...enum.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {

		roleValue := c.Locals("role")

		if roleValue == nil {
			return exception.NewForbiddenErr("You dont have permission")
		}

		role, ok := roleValue.(enum.UserRole)
		if !ok {
			if roleInt, ok := roleValue.(int); ok {
				role = enum.UserRole(roleInt)
			} else {
				return exception.NewBadRequestErr("Invalid role type")
			}
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				return c.Next()
			}
		}

		return exception.NewForbiddenErr("You do not have permission to access this resource")

	}
}
