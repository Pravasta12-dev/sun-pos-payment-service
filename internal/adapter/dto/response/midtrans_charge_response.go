package response

type MidtransChargeResponse struct {
	OrderID string `json:"order_id"`
	Actions []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"actions"`
}
