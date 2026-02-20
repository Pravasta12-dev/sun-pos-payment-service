package response

import (
	"sun-pos-payment-service/utils/enum"
	"time"
)

type GenerateOwnerVaResponse struct {
	OrderID    string                 `json:"order_id"`
	VaNumber   string                 `json:"va_number"`
	Bank       string                 `json:"bank"`
	BillID     string                 `json:"bill_id"`
	Status     enum.TransactionStatus `json:"status"`
	ExipiredAt *time.Time             `json:"expired_at,omitempty"`
}
