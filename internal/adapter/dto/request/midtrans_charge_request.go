package request

type MidtransChargeRequest struct {
	PaymentType        string               `json:"payment_type"`
	TransactionDetails TransactionDetails   `json:"transaction_details"`
	Qris               *Qris                `json:"qris,omitempty"`
	BankTransfer       *BankTransferRequest `json:"bank_transfer,omitempty"`
}

type TransactionDetails struct {
	OrderID     string  `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}

type Qris struct {
	Acquirer string `json:"acquirer"`
}

type BankTransferRequest struct {
	Bank string `json:"bank"`
}
