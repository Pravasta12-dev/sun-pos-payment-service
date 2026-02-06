package service

import "time"

type GenerateQRISInput struct {
	MerchantID    string
	ServerKey     string
	OrderID       string
	Amount        float64
	Acquirer      string
	ExpireMinutes int
}

type GenerateQRISResult struct {
	OrderID   string
	QrURL     string
	ExpiredAt *time.Time
	Status    string
}
