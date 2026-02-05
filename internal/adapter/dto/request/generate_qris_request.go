package request

type GenerateQrisRequest struct {
	MerchantID int64   `json:"merchant_id"`
	ServerKey  string  `json:"server_key"`
	OrderID    string  `json:"order_id"`
	Amount     float64 `json:"amount"`
	Acquirer   string  `json:"acquirer"`
}
