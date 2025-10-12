package handler

import (
	"strconv"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	GetAllUser(c *fiber.Ctx) error
	DeactiveUser(c *fiber.Ctx) error
	GetDetailUser(c *fiber.Ctx) error
}

type UserHandlerImpl struct {
	UserService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &UserHandlerImpl{
		UserService: userService,
	}
}

func (h *UserHandlerImpl) GetAllUser(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	sort := c.Query("sort", "created_at desc")

	// Parse filter parameters
	email := c.Query("email")
	role, _ := strconv.Atoi(c.Query("role"))
	username := c.Query("author_id")
	// status := c.Query("status")

	// Create filter params
	filter := dto.UserFilterRequest{
		Email:    email,
		Role:     role,
		Username: username,
		PaginationParams: dto.PaginationParams{
			Page:     page,
			PageSize: pageSize,
			Sort:     sort,
		},
	}
	datas, err := h.UserService.FindAllUserWithPaginatin(filter)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    datas,
		Status:  fiber.StatusOK,
		Message: "Success get list users",
	})
}

func (h *UserHandlerImpl) DeactiveUser(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

func (h *UserHandlerImpl) GetDetailUser(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}
