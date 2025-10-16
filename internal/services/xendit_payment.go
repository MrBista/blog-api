package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MrBista/blog-api/internal/config"
	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/repository"
	"gorm.io/gorm"
)

type XendiPaymentService struct {
	DB             *gorm.DB
	UserRepository repository.UserRepository
	Config         *config.Config
}

func NewXenditPaymentService(userRepository repository.UserRepository, db *gorm.DB) PaymentService {
	return &XendiPaymentService{
		DB:             db,
		UserRepository: userRepository,
	}
}

func (s *XendiPaymentService) CreateQrisPayment(userID uint, amount float64, durationMonths, paketPlan int) (*models.Subscription, error) {
	// Generate unique external ID
	externalID := fmt.Sprintf("sub-%d-%d-%d", userID, paketPlan, time.Now().Unix())

	// Create QRIS payment request
	reqBody := dto.CreateQRISRequest{
		ExternalID:  externalID,
		Type:        "DYNAMIC",
		CallbackURL: fmt.Sprintf("%s/webhook/xendit", s.Config.AppMain.GetBaseUrl()),
		Amount:      amount,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Call Xendit API
	req, err := http.NewRequest("POST", s.Config.Xendit.GetBaseUrl()+"/qr_codes", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(s.Config.Xendit.GetApiKey(), "")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Xendit API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("xendit API error: %s", string(body))
	}

	var qrisResp dto.QRISResponse
	if err := json.Unmarshal(body, &qrisResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	endDate := time.Now().AddDate(0, durationMonths, 0)

	subscription := &models.Subscription{
		UserID:        userID,
		PaymentID:     qrisResp.ID,
		ExternalID:    externalID,
		Amount:        amount,
		Status:        models.SubscriptionPending,
		PaymentMethod: "QRIS",
		QRString:      qrisResp.QRString,
		EndDate:       &endDate,
	}

	if err := s.DB.Create(subscription).Error; err != nil {
		return nil, fmt.Errorf("failed to save subscription: %w", err)
	}

	return subscription, nil
}
func (s *XendiPaymentService) HandleWebhook(payload []byte) error {
	var webhook dto.XenditWebhook
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return fmt.Errorf("invalid webhook payload: %w", err)
	}

	// Find subscription by external ID or payment ID
	var subscription models.Subscription
	if err := s.DB.Where("external_id = ? OR payment_id = ?", webhook.ExternalID, webhook.ID).First(&subscription).Error; err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	// Update subscription status based on payment status
	switch status := webhook.Status; status {
	case "COMPLETED":
		now := time.Now()
		subscription.Status = models.SubscriptionActive
		subscription.StartDate = &now

		// Update user subscription status
		if err := s.DB.Model(&models.User{}).Where("id = ?", subscription.UserID).Updates(map[string]interface{}{
			"is_subscribed":    true,
			"subscription_end": subscription.EndDate,
		}).Error; err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

	case "FAILED":
		subscription.Status = models.SubscriptionCancelled

	}

	// Save subscription update
	if err := s.DB.Save(&subscription).Error; err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil

}
func (s *XendiPaymentService) CheckExpiredSubscriptions() error {
	return nil
}
