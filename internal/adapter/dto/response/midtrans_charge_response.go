package response

type MidtransChargeResponse struct {
	OrderID   string     `json:"order_id"`
	Actions   []Actions  `json:"actions,omitempty"`
	VaNumbers []VaNumber `json:"va_numbers"`
}

type Actions struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VaNumber struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}
