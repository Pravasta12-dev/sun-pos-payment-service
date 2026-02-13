package service

import (
	"sun-pos-payment-service/utils/enum"
	"time"
)

type GenerateQRISInput struct {
	MerchantID    string
	ServerKey     string
	BillID        string
	Amount        float64
	Acquirer      string
	ExpireMinutes int
}

type GenerateQRISResult struct {
	OrderID   string
	QrURL     string
	BillID    string
	ExpiredAt *time.Time
	Status    enum.TransactionStatus
}

type GenerateOwnerQRISInput struct {
	OrderID       string
	Amount        float64
	ExpireMinutes int
}
