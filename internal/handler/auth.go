package handler

import (
	"encoding/json"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/MrBista/blog-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	LoginUser(c *fiber.Ctx) error
	RegisterUser(c *fiber.Ctx) error
	ConfirmOtp(c *fiber.Ctx) error
}

type AuthHandlerImpl struct {
	AuthService services.AuthService
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &AuthHandlerImpl{
		AuthService: authService,
	}
}

func (h *AuthHandlerImpl) LoginUser(c *fiber.Ctx) error {
	var loginReq dto.LoginRequest

	body := c.Body()

	if err := json.Unmarshal(body, &loginReq); err != nil {
		return err
	}

	validator := utils.GetValidator()

	if err := validator.Struct(loginReq); err != nil {
		return exception.NewValidationErr(err)
	}

	responseLogin, err := h.AuthService.LoginUser(loginReq)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    responseLogin,
		Status:  fiber.StatusOK,
		Message: "Successfully login user",
	})
}

func (h *AuthHandlerImpl) RegisterUser(c *fiber.Ctx) error {
	var userReq dto.RegisterRequest

	body := c.Body()

	if err := json.Unmarshal(body, &userReq); err != nil {
		return err
	}

	validator := utils.GetValidator()

	if err := validator.Struct(userReq); err != nil {
		return exception.NewValidationErr(err)
	}

	err := h.AuthService.RegisterUser(userReq)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Status:  fiber.StatusCreated,
		Message: "Successfully register user",
	})
}

func (h *AuthHandlerImpl) ConfirmOtp(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}
