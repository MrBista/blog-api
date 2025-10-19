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

type UserHandler interface {
	GetAllUser(c *fiber.Ctx) error
	DeactiveUser(c *fiber.Ctx) error
	GetDetailUser(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	FollowUser(c *fiber.Ctx) error
	UnfollowUser(c *fiber.Ctx) error
	GetListFollower(c *fiber.Ctx) error
	GetListFollowing(c *fiber.Ctx) error
	GetFollowerCount(c *fiber.Ctx) error
	GetFollowingCount(c *fiber.Ctx) error
	CheckFollowStatus(c *fiber.Ctx) error
	GetMyFollowers(c *fiber.Ctx) error
	GetMyFollowing(c *fiber.Ctx) error
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
	userToFollowParam := c.Params("id")

	userId, err := strconv.Atoi(userToFollowParam)
	if err != nil {
		return exception.NewBadRequestErr("Invalid user ID")
	}

	userResponseDetail, err := h.UserService.DetailUser(userId)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    userResponseDetail,
		Status:  fiber.StatusOK,
		Message: "Successfully get detail users",
	})
}

func (h *UserHandlerImpl) CreateUser(c *fiber.Ctx) error {

	var userBody dto.RegisterRequest

	body := c.Body()

	if err := json.Unmarshal(body, &userBody); err != nil {
		return err
	}

	validate := utils.GetValidator()

	if err := validate.Struct(&userBody); err != nil {
		return exception.NewValidationErr(err)
	}

	result, err := h.UserService.CreateUser(userBody)

	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(dto.CommonResponseSuccess{
		Data:    result,
		Status:  fiber.StatusCreated,
		Message: "successfully created user",
	})
}

func (h *UserHandlerImpl) FollowUser(c *fiber.Ctx) error {
	userToFollowParam := c.Params("id")

	userToFollow, err := strconv.Atoi(userToFollowParam)
	if err != nil {
		return exception.NewBadRequestErr("Invalid user ID")
	}

	userDetail, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	if err := h.UserService.FollowUser(userToFollow, userDetail); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "Successfully followed user",
		Status:  fiber.StatusCreated,
	})
}

func (h *UserHandlerImpl) UnfollowUser(c *fiber.Ctx) error {
	userToUnfollowParam := c.Params("id")

	userToUnfollow, err := strconv.Atoi(userToUnfollowParam)
	if err != nil {
		return exception.NewBadRequestErr("Invalid user ID")
	}

	userDetail, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	if err := h.UserService.UnFollowUser(userToUnfollow, userDetail); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Message: "Successfully unfollowed user",
		Status:  fiber.StatusOK,
	})
}

func (h *UserHandlerImpl) GetListFollower(c *fiber.Ctx) error {
	userIdParam := c.Params("id")

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		return exception.NewBadRequestErr("Invalid user ID")
	}

	followers, err := h.UserService.GetListFollower(userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    followers,
		Message: "Successfully retrieved followers",
		Status:  fiber.StatusOK,
	})
}

func (h *UserHandlerImpl) GetListFollowing(c *fiber.Ctx) error {
	userIdParam := c.Params("id")

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		return exception.NewBadRequestErr("Invalid user ID")
	}

	following, err := h.UserService.GetListFollowing(userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    following,
		Message: "Successfully retrieved following",
		Status:  fiber.StatusOK,
	})
}

func (h *UserHandlerImpl) GetFollowerCount(c *fiber.Ctx) error {
	userIdParam := c.Params("id")

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		return exception.NewBadRequestErr("Invalid user ID")
	}

	count, err := h.UserService.CountFollower(userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data: map[string]int64{
			"follower_count": count,
		},
		Message: "Successfully retrieved follower count",
		Status:  fiber.StatusOK,
	})
}

func (h *UserHandlerImpl) GetFollowingCount(c *fiber.Ctx) error {
	userIdParam := c.Params("id")

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		return exception.NewBadRequestErr("Invalid user ID")
	}

	count, err := h.UserService.CountFollowing(userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data: map[string]int64{
			"following_count": count,
		},
		Message: "Successfully retrieved following count",
		Status:  fiber.StatusOK,
	})
}

func (h *UserHandlerImpl) CheckFollowStatus(c *fiber.Ctx) error {
	targetUserIdParam := c.Params("id")

	targetUserId, err := strconv.Atoi(targetUserIdParam)
	if err != nil {
		return exception.NewBadRequestErr("Invalid user ID")
	}

	userDetail, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	isFollowing, err := h.UserService.CheckIsFollowing(targetUserId, userDetail.UserId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data: map[string]bool{
			"is_following": isFollowing,
		},
		Message: "Successfully checked follow status",
		Status:  fiber.StatusOK,
	})
}

func (h *UserHandlerImpl) GetMyFollowers(c *fiber.Ctx) error {
	userDetail, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	followers, err := h.UserService.GetListFollower(userDetail.UserId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    followers,
		Message: "Successfully retrieved my followers",
		Status:  fiber.StatusOK,
	})
}

func (h *UserHandlerImpl) GetMyFollowing(c *fiber.Ctx) error {
	userDetail, err := utils.GetUserClaims(c)
	if err != nil {
		return err
	}

	following, err := h.UserService.GetListFollowing(userDetail.UserId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    following,
		Message: "Successfully retrieved my following",
		Status:  fiber.StatusOK,
	})
}
