package middleware

import (
	"log"

	"github.com/MrBista/blog-api/internal/exception"
	"github.com/gofiber/fiber/v2"
)

func HandleError(c *fiber.Ctx, err error) error {
	if customErr, ok := err.(*exception.ErrorCustom); ok {
		log.Printf("[ERROR] %s: %s", customErr.Code, customErr.Message)
		return c.Status(customErr.GetStatusHttp()).JSON(customErr)
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  fiber.StatusInternalServerError,
		"code":    exception.ERR_UNHANDLE,
		"message": "Internal Server Error",
	})

}
