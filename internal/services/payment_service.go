package services

import "github.com/MrBista/blog-api/internal/models"

type PaymentService interface {
	CreateQrisPayment(userID uint, amount float64, durationMonths, paketPlan int) (*models.Subscription, error)
	HandleWebhook(payload []byte) error
	CheckExpiredSubscriptions() error
}
