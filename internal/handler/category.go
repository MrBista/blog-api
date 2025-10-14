package handler

import (
	"encoding/json"
	"strconv"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/MrBista/blog-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler interface {
	FindAllCategory(c *fiber.Ctx) error
	FindCategoryById(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
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

	categoriesRes, err := h.CategoryService.FindAllCategory()

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    categoriesRes,
		Status:  fiber.StatusOK,
		Message: "Successfully get all category",
	})
}

func (h *CategoryHandlerImpl) FindCategoryById(c *fiber.Ctx) error {
	// panic("not implemented") // TODO: Implement

	paramId := c.Params("id")

	val, err := strconv.Atoi(paramId)

	if err != nil {
		return exception.NewBadRequestErr(err.Error())
	}

	categoryDetail, err := h.CategoryService.FindById(val)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    categoryDetail,
		Message: "Successfully get detail category",
		Status:  fiber.StatusOK,
	})
}

func (h *CategoryHandlerImpl) CreateCategory(c *fiber.Ctx) error {
	var categoryReqbody dto.CategoryRequst

	body := c.Body()

	if err := json.Unmarshal(body, &categoryReqbody); err != nil {
		return err
	}

	validator := utils.GetValidator()

	if err := validator.Struct(&categoryReqbody); err != nil {
		return exception.NewValidationErr(err)
	}

	err := h.CategoryService.CreateCategory(categoryReqbody)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "Successfully create category",
		Status:  fiber.StatusCreated,
	})
}

func (h *CategoryHandlerImpl) UpdateCategory(c *fiber.Ctx) error {
	var categoryReqbody dto.CategoryRequst

	body := c.Body()

	if err := json.Unmarshal(body, &categoryReqbody); err != nil {
		return err
	}

	paramId := c.Params("id")
	valId, err := strconv.Atoi(paramId)

	if err != nil {
		return exception.NewBadRequestErr(err.Error())
	}

	err = h.CategoryService.UpdateCategory(valId, categoryReqbody)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "Successfully update category",
		Status:  fiber.StatusOK,
	})
}

func (h *CategoryHandlerImpl) DeleteCategory(c *fiber.Ctx) error {

	paramId := c.Params("id")
	valId, err := strconv.Atoi(paramId)

	if err != nil {
		return exception.NewBadRequestErr(err.Error())
	}

	err = h.CategoryService.DeleteById(valId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "Successfully delete category",
		Status:  fiber.StatusOK,
	})
}
