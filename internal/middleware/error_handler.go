package middleware

import (
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func HandleError(c *fiber.Ctx, err error) error {
	if customErr, ok := err.(*exception.ErrorCustom); ok {
		// log.Printf("[ERROR] %s: %s", customErr.Code, customErr.Message)
		utils.Logger.Errorf("Error code %s with message %s", customErr.Code, customErr.Message)
		return c.Status(customErr.GetStatusHttp()).JSON(customErr)
	}

	utils.Logger.Errorf("something went wrong %v", err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  fiber.StatusInternalServerError,
		"code":    exception.ERR_UNHANDLE,
		"message": "Internal Server Error",
	})

}
