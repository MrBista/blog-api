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

type ReadingListHandler interface {
	CreateReadingList(c *fiber.Ctx) error
	GetReadingLists(c *fiber.Ctx) error
	GetReadingListByID(c *fiber.Ctx) error
	UpdateReadingList(c *fiber.Ctx) error
	DeleteReadingList(c *fiber.Ctx) error
	CreateSavedPost(c *fiber.Ctx) error
	GetSavedPosts(c *fiber.Ctx) error
	UpdateSavedPost(c *fiber.Ctx) error
	DeleteSavedPost(c *fiber.Ctx) error
	DeleteSavedPostByPostAndList(c *fiber.Ctx) error
	MarkAllAsRead(c *fiber.Ctx) error
}

type ReadingListHandlerImpl struct {
	PostService services.PostService
}

func NewReadingListHandler(postService services.PostService) ReadingListHandler {
	return &ReadingListHandlerImpl{
		PostService: postService,
	}
}

// !TODO refactoring di pindah di handler terpisah
func (h *ReadingListHandlerImpl) CreateReadingList(c *fiber.Ctx) error {
	var readingListDto dto.CreateReadingListRequest
	body := c.Body()
	if err := json.Unmarshal(body, &readingListDto); err != nil {
		return err
	}

	validator := utils.GetValidator()
	if err := validator.Struct(&readingListDto); err != nil {
		return exception.NewValidationErr(err)
	}

	detailUser, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	if err := h.PostService.CreateReadingList(readingListDto, detailUser); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "successfully create new reading list",
		Status:  fiber.StatusCreated,
	})
}

func (h *ReadingListHandlerImpl) GetReadingLists(c *fiber.Ctx) error {
	detailUser, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	data, err := h.PostService.GetReadingLists(detailUser)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    data,
		Message: "successfully get reading lists",
		Status:  fiber.StatusOK,
	})
}

func (h *ReadingListHandlerImpl) GetReadingListByID(c *fiber.Ctx) error {
	listID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return exception.NewBadRequestErr("Invalid reading list ID")
	}

	detailUser, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	data, err := h.PostService.GetReadingListByID(listID, detailUser)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    data,
		Message: "successfully get reading list detail",
		Status:  fiber.StatusOK,
	})
}

func (h *ReadingListHandlerImpl) UpdateReadingList(c *fiber.Ctx) error {
	listID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return exception.NewBadRequestErr("Invalid reading list ID")
	}

	var updateDto dto.UpdateReadingListRequest
	body := c.Body()
	if err := json.Unmarshal(body, &updateDto); err != nil {
		return err
	}

	validator := utils.GetValidator()
	if err := validator.Struct(&updateDto); err != nil {
		return exception.NewValidationErr(err)
	}

	detailUser, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	if err := h.PostService.UpdateReadingList(listID, updateDto, detailUser); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "successfully update reading list",
		Status:  fiber.StatusOK,
	})
}

func (h *ReadingListHandlerImpl) DeleteReadingList(c *fiber.Ctx) error {
	listID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return exception.NewBadRequestErr("Invalid reading list ID")
	}

	detailUser, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	if err := h.PostService.DeleteReadingList(listID, detailUser); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "successfully delete reading list",
		Status:  fiber.StatusOK,
	})
}

func (h *ReadingListHandlerImpl) CreateSavedPost(c *fiber.Ctx) error {
	var savedPostDto dto.CreateSavedPostRequest
	body := c.Body()
	if err := json.Unmarshal(body, &savedPostDto); err != nil {
		return err
	}

	validator := utils.GetValidator()
	if err := validator.Struct(&savedPostDto); err != nil {
		return exception.NewValidationErr(err)
	}

	detailUser, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	if err := h.PostService.CreateSavedPost(savedPostDto, detailUser); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "successfully save post to reading list",
		Status:  fiber.StatusCreated,
	})
}

func (h *ReadingListHandlerImpl) GetSavedPosts(c *fiber.Ctx) error {
	readingListID, err := strconv.ParseInt(c.Params("listId"), 10, 64)
	if err != nil {
		return exception.NewBadRequestErr("Invalid reading list ID")
	}

	detailUser, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	data, err := h.PostService.GetSavedPosts(readingListID, detailUser)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    data,
		Message: "successfully get saved posts",
		Status:  fiber.StatusOK,
	})
}

func (h *ReadingListHandlerImpl) UpdateSavedPost(c *fiber.Ctx) error {
	savedPostID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return exception.NewBadRequestErr("Invalid saved post ID")
	}

	var updateDto dto.UpdateSavedPostRequest
	body := c.Body()
	if err := json.Unmarshal(body, &updateDto); err != nil {
		return err
	}

	validator := utils.GetValidator()
	if err := validator.Struct(&updateDto); err != nil {
		return exception.NewValidationErr(err)
	}

	detailUser, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	if err := h.PostService.UpdateSavedPost(savedPostID, updateDto, detailUser); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "successfully update saved post",
		Status:  fiber.StatusOK,
	})
}

func (h *ReadingListHandlerImpl) DeleteSavedPost(c *fiber.Ctx) error {
	savedPostID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return exception.NewBadRequestErr("Invalid saved post ID")
	}

	detailUser, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	if err := h.PostService.DeleteSavedPost(savedPostID, detailUser); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "successfully delete saved post",
		Status:  fiber.StatusOK,
	})
}

func (h *ReadingListHandlerImpl) DeleteSavedPostByPostAndList(c *fiber.Ctx) error {
	postID, err := strconv.ParseInt(c.Query("postId"), 10, 64)
	if err != nil {
		return exception.NewBadRequestErr("Invalid post ID")
	}

	readingListID, err := strconv.ParseInt(c.Query("readingListId"), 10, 64)
	if err != nil {
		return exception.NewBadRequestErr("Invalid reading list ID")
	}

	detailUser, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	if err := h.PostService.DeleteSavedPostByPostAndList(postID, readingListID, detailUser); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "successfully delete saved post",
		Status:  fiber.StatusOK,
	})
}

func (h *ReadingListHandlerImpl) MarkAllAsRead(c *fiber.Ctx) error {
	readingListID, err := strconv.ParseInt(c.Params("listId"), 10, 64)
	if err != nil {
		return exception.NewBadRequestErr("Invalid reading list ID")
	}

	detailUser, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	if err := h.PostService.MarkAllAsRead(readingListID, detailUser); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "successfully mark all posts as read",
		Status:  fiber.StatusOK,
	})
}
