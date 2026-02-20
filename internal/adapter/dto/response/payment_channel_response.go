package response

import "sun-pos-payment-service/utils/enum"

type PaymentChannelResponse struct {
	ID       int64                   `json:"id"`
	Type     enum.PaymentChannelType `json:"type"`
	Code     enum.PaymentChannelCode `json:"code"`
	Label    string                  `json:"label"`
	FeeType  enum.PaymentFeeType     `json:"fee_type"`
	FeeValue float64                 `json:"fee_value"`
}
