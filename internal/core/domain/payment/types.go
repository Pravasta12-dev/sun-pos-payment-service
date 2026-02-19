package payment

import "time"

type QrisChargeInput struct {
	OrderID  string
	Amount   float64
	Acquirer string
}

type QrisChargeResult struct {
	OrderID   string
	QrURL     string
	ExpiredAt *time.Time
}

type VaChargeInput struct {
	OrderID string
	Amount  float64
	Bank    string
}

type VaChargeResult struct {
	OrderID  string
	VANumber string
	Bank     string
}
