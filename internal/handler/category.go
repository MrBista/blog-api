package handler

import (
	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler interface {
	FindAllCategory(c *fiber.Ctx) error
	FindCategoryById(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
}

type CategoryHandlerImpl struct {
	CategoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) CategoryHandler {
	return &CategoryHandlerImpl{
		CategoryService: categoryService,
	}
}

func (h *CategoryHandlerImpl) FindAllCategory(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

func (h *CategoryHandlerImpl) FindCategoryById(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

func (h *CategoryHandlerImpl) CreateCategory(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

func (h *CategoryHandlerImpl) UpdateCategory(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}
