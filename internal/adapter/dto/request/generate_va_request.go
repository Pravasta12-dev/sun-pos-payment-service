package request

type GenerateOwnerVARequest struct {
	OrderID string  `json:"order_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required,gt=0"`
	Bank    string  `json:"bank" validate:"required"`
}
