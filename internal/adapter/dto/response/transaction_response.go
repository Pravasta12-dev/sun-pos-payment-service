package response

import (
	"sun-pos-payment-service/utils/enum"
	"time"
)

type TransactionResponse struct {
	ID          int64                  `json:"id"`
	MerchantID  string                 `json:"merchant_id"`
	OrderID     string                 `json:"order_id"`
	Amount      float64                `json:"amount"`
	PaymentType string                 `json:"payment_type"`
	Status      enum.TransactionStatus `json:"status"`
	PaidAt      *time.Time             `json:"paid_at,omitempty"`
	ExpiredAt   *time.Time             `json:"expired_at,omitempty"`
}
