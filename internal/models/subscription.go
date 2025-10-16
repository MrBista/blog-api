package models

import "time"

type SubscriptionStatus int

const (
	SubscriptionPending SubscriptionStatus = iota + 1
	SubscriptionActive
	SubscriptionExpired
	SubscriptionCancelled
)

type Subscription struct {
	ID            uint               `gorm:"column:id;primarykey" json:"id"`
	UserID        uint               `gorm:"column:user_id;not null" json:"userId"`
	PaymentID     string             `gorm:"column:payment_id;uniqueIndex" json:"paymentId"` // Xendit payment ID
	ExternalID    string             `gorm:"column:external_id;uniqueIndex" json:"externalId"`
	Amount        float64            `gorm:"column:amount" json:"amount"`
	Status        SubscriptionStatus `gorm:"column:status;default:'1'" json:"status"`
	PaymentMethod string             `gorm:"column:payment_method" json:"paymentMethod"` // QRIS
	QRString      string             `gorm:"column:qr_string;type:text" json:"qrString,omitempty"`
	QRImageURL    string             `gorm:"column:qr_image_url" json:"qrImageUrl,omitempty"`
	StartDate     *time.Time         `gorm:"column:start_date" json:"startDate"`
	EndDate       *time.Time         `gorm:"column:end_date" json:"endDate"`
	CreatedAt     time.Time          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time          `gorm:"column:updated_at" json:"updatedAt"`
}

func (s *Subscription) TableName() string {
	return "subscriptions"
}
