package middleware

import (
	"strconv"

	"github.com/MrBista/blog-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func ChatMessageMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		roomIdStr := c.Params("roomId")
		roomId, err := strconv.ParseUint(roomIdStr, 10, 32)
		if err != nil {
			return err
		}

		userDetail, err := utils.GetUserClaims(c)

		if err != nil {
			return err
		}

		c.Locals("room_id", roomId)
		c.Locals("user_id", userDetail.ID)

		return c.Next()
	}
}
