package request

type MidtransWebhookRequest struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	PaymentType       string `json:"payment_type"`
	SettlementTime    string `json:"settlement_time"`
}
