package model

import (
	"sun-pos-payment-service/utils/enum"
	"time"
)

type TransactionModel struct {
	ID               int64
	MerchantID       string
	OrderID          string
	BillID           string
	Amount           float64
	PaymentType      string
	Status           enum.TransactionStatus
	TransactionScope enum.TransactionScope
	QrURL            *string
	PaidAt           *time.Time
	ExpiredAt        *time.Time
}
