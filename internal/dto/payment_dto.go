package dto

import "time"

type CreateQRISRequest struct {
	ExternalID  string  `json:"externalId"`
	Type        string  `json:"type"`
	CallbackURL string  `json:"callbackUrl"`
	Amount      float64 `json:"amount"`
}

type QRISResponse struct {
	ID         string    `json:"id"`
	ExternalID string    `json:"externalId"`
	Amount     float64   `json:"amount"`
	QRString   string    `json:"qrString"`
	Status     string    `json:"status"`
	Created    time.Time `json:"created"`
}

type XenditWebhook struct {
	ID         string  `json:"id"`
	ExternalID string  `json:"external_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	QRCode     string  `json:"qr_code"`
}
