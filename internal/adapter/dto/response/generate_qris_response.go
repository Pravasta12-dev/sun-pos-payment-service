package response

import (
	"sun-pos-payment-service/utils/enum"
	"time"
)

type GenerateQrisResponse struct {
	OrderID   string                 `json:"order_id"`
	QrUrl     string                 `json:"qr_url"`
	BillID    string                 `json:"bill_id"`
	ExpiredAt *time.Time             `json:"expired_at"`
	Status    enum.TransactionStatus `json:"status"`
}
