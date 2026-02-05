package model

import "time"

type TransactionModel struct {
	ID          int64
	MerchantID  int64
	OrderID     string
	Amount      float64
	PaymentType string
	Status      string
	PaidAt      *time.Time
	ExpiredAt   *time.Time
}
