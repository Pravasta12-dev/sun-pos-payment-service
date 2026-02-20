package request

type GenerateQrisRequest struct {
	MerchantID string  `json:"merchant_id"`
	ServerKey  string  `json:"server_key"`
	OrderID    string  `json:"order_id"`
	Amount     float64 `json:"amount"`
	Acquirer   string  `json:"acquirer"`
}

type GenerateOwnerQrisRequest struct {
	BillID        string  `json:"bill_id"`
	Amount        float64 `json:"amount"`
	Acquirer      string  `json:"acquirer"`
	ExpireMinutes int     `json:"expire_minutes"`
}
