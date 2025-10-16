package handler

import (
	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/MrBista/blog-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type SubscriptionHandler interface {
	CreateSubscription(c *fiber.Ctx) error
	WebhookPayment(c *fiber.Ctx) error
}

type SubscriptionHandlerImpl struct {
	xenditService services.PaymentService
}

func NewSubscriptionHandler(xenditService services.PaymentService) SubscriptionHandler {
	return &SubscriptionHandlerImpl{
		xenditService: xenditService,
	}
}

type CreateSubscriptionRequest struct {
	Plan string `json:"plan"` // monthly, yearly
}

// CreateSubscription membuat subscription baru
func (h *SubscriptionHandlerImpl) CreateSubscription(c *fiber.Ctx) error {
	userDetail, err := utils.GetUserClaims(c)

	if err != nil {
		return err
	}

	if userDetail == nil {
		return exception.NewUnAuthorizationErr("unauthorized")
	}

	userID := uint(userDetail.UserId)

	var req CreateSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.NewBadRequestErr("Invalid request body")
	}

	var amount float64
	var durationMonths int

	var paketPlan int

	switch req.Plan {
	case "monthly":
		amount = 50000 // Rp 50.000
		durationMonths = 1
		paketPlan = 1 // monthly
	case "yearly":
		amount = 500000 // Rp 500.000
		durationMonths = 12
		paketPlan = 2 // yearloy
	default:
		return exception.NewBadRequestErr("Invalid plan. Choose 'monthly' or 'yearly'")
	}

	subscription, err := h.xenditService.CreateQrisPayment(userID, amount, durationMonths, paketPlan)
	if err != nil {
		return exception.NewBadRequestErr("Failed to create payment: " + err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    subscription,
		Message: "Subscription created. Please scan the QR code to complete payment.",
		Status:  fiber.StatusOK,
	})
}

// XenditWebhook handles webhook dari Xendit
func (h *SubscriptionHandlerImpl) WebhookPayment(c *fiber.Ctx) error {
	// Read webhook payload

	// Verify webhook signature (recommended for production)
	// webhookToken := c.Get("x-callback-token")
	// if webhookToken != h.xenditService.config.WebhookKey {
	//     return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	//         "error": "Invalid webhook token",
	//     })
	// }

	// Handle webhook
	if err := h.xenditService.HandleWebhook(c.Body()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process webhook: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.CommonResponseSuccess{
		Data:    true,
		Status:  fiber.StatusOK,
		Message: "Successfully bayar",
	})
}
