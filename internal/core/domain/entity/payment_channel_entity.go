package entity

import (
	"sun-pos-payment-service/utils/enum"
	"time"
)

type PaymentChannelEntity struct {
	ID        int64                   `db:"id"`
	Type      enum.PaymentChannelType `db:"type"`
	Code      enum.PaymentChannelCode `db:"code"`
	Label     string                  `db:"label"`
	FeeType   enum.PaymentFeeType     `db:"fee_type"`
	FeeValue  float64                 `db:"fee_value"`
	IsActive  bool                    `db:"is_active"`
	CreatedAt time.Time               `db:"created_at"`
}

func (p *PaymentChannelEntity) TableName() string {
	return "payment_channels"
}
