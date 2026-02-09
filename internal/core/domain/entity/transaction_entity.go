package entity

import "time"

type TransactionEntity struct {
	ID          int64      `db:"id"`
	MerchantID  string     `db:"merchant_id"`
	OrderID     string     `db:"order_id"`
	BillID      string     `db:"bill_id"`
	Amount      float64    `db:"amount"`
	PaymentType string     `db:"payment_type"`
	Status      string     `db:"status"`
	PaidAt      *time.Time `db:"paid_at"`
	QrURL       *string    `db:"qr_url"`
	ExpiredAt   *time.Time `db:"expired_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

// TableName overrides the table name used by GORM
func (TransactionEntity) TableName() string {
	return "transactions"
}
