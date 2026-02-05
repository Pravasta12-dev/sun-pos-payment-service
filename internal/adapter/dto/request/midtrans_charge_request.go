package request

type MidtransChargeRequest struct {
	PaymentType string `json:"payment_type"`

	TransactionDetails struct {
		OrderID     string  `json:"order_id"`
		GrossAmount float64 `json:"gross_amount"`
	} `json:"transaction_details"`

	Qris struct {
		Acquirer string `json:"acquirer"`
	} `json:"qris"`
}
