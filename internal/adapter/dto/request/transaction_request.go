package request

type TransactionRequest struct {
	MerchantID string `json:"merchant_id" validate:"required"`
	BillID     string `json:"bill_id" validate:"required"`
}
