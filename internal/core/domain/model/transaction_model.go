package model

import "time"

type TransactionModel struct {
	ID          int64
	MerchantID  string
	OrderID     string
	BillID      string
	Amount      float64
	PaymentType string
	Status      string
	QrURL       *string
	PaidAt      *time.Time
	ExpiredAt   *time.Time
}
