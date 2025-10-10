package handler

import (
	"encoding/json"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/MrBista/blog-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Post interface {
	CreatePost(c *fiber.Ctx) error
	GetAllPosts(c *fiber.Ctx) error
	GetPostBySlug(c *fiber.Ctx) error
	UpdatePost(c *fiber.Ctx) error
	DeletePost(c *fiber.Ctx) error
}

type PostImpl struct {
	PostService services.PostService
}

func NewHandlerPost(postService services.PostService) Post {
	return &PostImpl{
		PostService: postService,
	}
}

func (h *PostImpl) CreatePost(c *fiber.Ctx) error {
	body := c.Body()

	var reqPost dto.CreatePostRequest

	if err := json.Unmarshal(body, &reqPost); err != nil {

		return err
	}

	validator := utils.GetValidator()

	if err := validator.Struct(reqPost); err != nil {
		return exception.NewValidationErr(err)
	}

	err := h.PostService.CreatePost(&reqPost)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Status:  fiber.StatusCreated,
		Message: "Successfully create post",
	})
}

func (h *PostImpl) GetAllPosts(c *fiber.Ctx) error {
	responsePost, err := h.PostService.FindAllPost()
	if err != nil {
		return err
	}

	return c.
		Status(fiber.StatusOK).
		JSON(dto.
			CommonResponseSuccess{Data: responsePost, Status: fiber.StatusOK, Message: "Successfully get all posts"},
		)
}

func (h *PostImpl) GetPostBySlug(c *fiber.Ctx) error {
	slugParam := c.Params("slug")

	postDetial, err := h.PostService.FindDetailPost(slugParam)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    postDetial,
		Status:  fiber.StatusOK,
		Message: "Success",
	})

}

func (h *PostImpl) UpdatePost(c *fiber.Ctx) error {
	body := c.Body()

	var updateBody dto.UpdatePostRequest

	if err := json.Unmarshal(body, &updateBody); err != nil {
		return err
	}

	validator := utils.GetValidator()

	if err := validator.Struct(updateBody); err != nil {
		return exception.NewValidationErr(err)
	}

	slugParam := c.Params("slug")

	updateBody.Slug = slugParam

	err := h.PostService.UpdatePost(&updateBody)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Status:  fiber.StatusOK,
		Message: "Successfully update posts",
	})

}

func (h *PostImpl) DeletePost(c *fiber.Ctx) error {
	slug := c.Params("slug")

	err := h.PostService.DeletePost(slug)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Status:  fiber.StatusOK,
		Message: "Success to delete posts",
	})
}
