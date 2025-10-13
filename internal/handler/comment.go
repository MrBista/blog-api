package handler

import (
	"encoding/json"
	"strconv"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/MrBista/blog-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type CommentHandler interface {
	FindAllComment(c *fiber.Ctx) error
	CreateComment(c *fiber.Ctx) error
}

type CommentHandlerImpl struct {
	CommentService services.CommentService
}

func NewCommentHandler(commentService services.CommentService) CommentHandler {
	return &CommentHandlerImpl{
		CommentService: commentService,
	}
}

func (h *CommentHandlerImpl) FindAllComment(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	sort := c.Query("sort", "created_at desc")
	postId, _ := strconv.Atoi(c.Params("postId"))

	filter := dto.CommentFilterRequest{
		PostId: postId,
		PaginationParams: dto.PaginationParams{
			Page:     page,
			PageSize: pageSize,
			Sort:     sort,
		},
	}

	userDetail, err := utils.GetUserClaims(c)

	if err != nil {
		return err
	}

	data, err := h.CommentService.FindAllCommentByPostId(filter, *userDetail)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    data,
		Status:  fiber.StatusOK,
		Message: "Successfully get list comment by post id",
	})
}

func (h *CommentHandlerImpl) CreateComment(c *fiber.Ctx) error {

	var commentBody dto.CommentRequest

	body := c.Body()

	if err := json.Unmarshal(body, &commentBody); err != nil {
		return err
	}

	userDetail, err := utils.GetUserClaims(c)

	if err != nil {
		return err
	}

	data, err := h.CommentService.CreateComment(commentBody, *userDetail)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CommonResponseSuccess{
		Data:    data,
		Status:  fiber.StatusCreated,
		Message: "Successfully create comment",
	})
}
