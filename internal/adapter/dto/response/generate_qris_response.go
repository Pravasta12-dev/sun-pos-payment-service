package response

import "time"

type GenerateQrisResponse struct {
	OrderID   string     `json:"order_id"`
	QrUrl     string     `json:"qr_url"`
	ExpiredAt *time.Time `json:"expired_at"`
	Status    string     `json:"status"`
}
