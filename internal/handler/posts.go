package handler

import (
	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/services"
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
	panic("not implemented") // TODO: Implement
}

func (h *PostImpl) GetAllPosts(c *fiber.Ctx) error {
	// panic("not implemented") // TODO: Implement
	responsePost, err := h.PostService.FindAllPost()
	if err != nil {
		dataResponse := map[string]interface{}{
			"data":    false,
			"message": err,
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dataResponse)
	}

	return c.
		Status(fiber.StatusOK).
		JSON(dto.
			CommonResponseSuccess{Data: responsePost, Status: fiber.StatusOK, Message: "Successfully get all posts"},
		)
}

func (h *PostImpl) GetPostBySlug(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

func (h *PostImpl) UpdatePost(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

func (h *PostImpl) DeletePost(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}
