package repository

import (
	"errors"
	"sun-pos-payment-service/internal/core/domain/entity"
	"sun-pos-payment-service/internal/core/domain/model"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type TransactionRepositoryInterface interface {
	CreateTransaction(
		merchantID string,
		orderID string,
		amount float64,
		paymentType string,
		expiredAt *time.Time,
	) (*model.TransactionModel, error)
	UpdateStatus(
		orderID string,
		status string,
		paidAt *time.Time,
	) error
	FindByOrderID(orderID string) (*model.TransactionModel, error)
}

type transactionRepository struct {
	db *gorm.DB
}

// CreateTransaction implements [TransactionRepositoryInterface].
func (t *transactionRepository) CreateTransaction(
	merchantID string,
	orderID string,
	amount float64,
	paymentType string,
	expiredAt *time.Time,
) (*model.TransactionModel, error) {
	query := `
		INSERT INTO transactions
			(merchant_id, order_id, amount, payment_type, status, expired_at)
		VALUES
			($1, $2, $3, $4, 'pending', $5)
		RETURNING id
	`

	var e entity.TransactionEntity

	err := t.db.Raw(
		query,
		merchantID,
		orderID,
		amount,
		paymentType,
		expiredAt,
	).Scan(&e.ID).Error

	if err != nil {
		log.Errorf("[TransactionRepository-1] failed to create transaction: %v", err)
		return nil, err
	}

	e.MerchantID = merchantID
	e.OrderID = orderID
	e.Amount = amount
	e.PaymentType = paymentType
	e.Status = model.TransactionStatusPending
	e.ExpiredAt = expiredAt

	transactionModel := &model.TransactionModel{
		ID:          e.ID,
		MerchantID:  e.MerchantID,
		OrderID:     e.OrderID,
		Amount:      e.Amount,
		PaymentType: e.PaymentType,
		Status:      e.Status,
		ExpiredAt:   e.ExpiredAt,
	}

	return transactionModel, nil
}

// FindByOrderID implements [TransactionRepositoryInterface].
func (t *transactionRepository) FindByOrderID(orderID string) (*model.TransactionModel, error) {
	query := `
		SELECT id, merchant_id, order_id, amount, payment_type, status, paid_at, expired_at
		FROM transactions
		WHERE order_id = $1
	`

	var e entity.TransactionEntity

	err := t.db.Raw(query, orderID).Scan(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[TransactionRepository-3] transaction not found: %v", err)
			return nil, err
		}

		log.Errorf("[TransactionRepository-4] failed to find transaction by order ID: %v", err)
		return nil, err
	}

	transactionModel := &model.TransactionModel{
		ID:          e.ID,
		MerchantID:  e.MerchantID,
		OrderID:     e.OrderID,
		Amount:      e.Amount,
		PaymentType: e.PaymentType,
		Status:      e.Status,
		PaidAt:      e.PaidAt,
		ExpiredAt:   e.ExpiredAt,
	}

	return transactionModel, nil
}

// UpdateStatus implements [TransactionRepositoryInterface].
func (t *transactionRepository) UpdateStatus(
	orderID string,
	status string,
	paidAt *time.Time,
) error {
	query := `
		UPDATE transactions
		SET status = $1,
			paid_at = $2,
			updated_at = NOW()
		WHERE order_id = $3
	`

	if err := t.db.Exec(query, status, paidAt, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[TransactionRepository-5] transaction not found: %v", err)
			return err
		}
		log.Errorf("[TransactionRepository-6] failed to update transaction status: %v", err)
		return err
	}

	return nil
}

func NewTransactionRepository(db *gorm.DB) TransactionRepositoryInterface {
	return &transactionRepository{
		db: db,
	}
}
