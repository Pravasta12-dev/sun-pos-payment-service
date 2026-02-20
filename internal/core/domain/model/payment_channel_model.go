package model

import "sun-pos-payment-service/utils/enum"

type PaymentChannelModel struct {
	ID       int64
	Type     enum.PaymentChannelType
	Code     enum.PaymentChannelCode
	Label    string
	FeeType  enum.PaymentFeeType
	FeeValue float64
}
