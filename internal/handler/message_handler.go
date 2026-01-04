package handler

import (
	"strconv"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/MrBista/blog-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type MessageHandler interface {
	FindOrCreateRoom(c *fiber.Ctx) error
}

type MessageHandlerImpl struct {
	MessageService services.MessageService
}

func NewMessageServiceImpl(messageService services.MessageService) MessageHandler {

	return &MessageHandlerImpl{
		MessageService: messageService,
	}
}

func (h *MessageHandlerImpl) FindOrCreateRoom(c *fiber.Ctx) error {

	userDetail, err := utils.GetUserClaims(c)

	if err != nil {
		return err
	}

	blogIdToParse := c.Params("blogId")

	var blogId int

	if blogIdToParse != "" {
		blogIdNew, err := strconv.Atoi(blogIdToParse)
		if err != nil {
			return err
		}
		blogId = blogIdNew
	}

	detailRoom, err := h.MessageService.FindOrCreateRoom(blogId, userDetail.UserId, 0)

	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    detailRoom,
		Status:  fiber.StatusOK,
		Message: "Successfully get detail room",
	})
}
