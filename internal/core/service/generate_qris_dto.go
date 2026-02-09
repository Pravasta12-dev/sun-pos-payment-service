package service

import "time"

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
	Status    string
}
