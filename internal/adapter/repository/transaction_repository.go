package repository

import (
	"errors"
	"fmt"
	"sun-pos-payment-service/internal/core/domain/entity"
	"sun-pos-payment-service/internal/core/domain/model"
	"sun-pos-payment-service/utils/enum"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type TransactionRepositoryInterface interface {
	CreateTransaction(
		scope enum.TransactionScope,
		merchantID *string,
		billID *string,
		orderID string,
		amount float64,
		paymentType string,
		qrURL string,
		expiredAt *time.Time,
	) (*model.TransactionModel, error)
	UpdateStatus(
		orderID string,
		status enum.TransactionStatus,
		paidAt *time.Time,
	) error
	FindByOrderID(orderID string) (*model.TransactionModel, error)
	FindByBillID(merchantID, billID string) (*model.TransactionModel, error)
	FindActivePendingTransaction(merchantID string, billID string, paymentType string) (*model.TransactionModel, error)
}

type transactionRepository struct {
	db *gorm.DB
}

// FindByBillID implements [TransactionRepositoryInterface].
func (t *transactionRepository) FindByBillID(merchantID, billID string) (*model.TransactionModel, error) {
	var e entity.TransactionEntity

	// FETCH USING WHERE().FIRST() with status 'pending' AND not expired
	err := t.db.Where("merchant_id = ? AND bill_id = ? AND (expired_at IS NULL OR expired_at > NOW())", merchantID, billID).Order("created_at DESC").Limit(1).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[TransactionRepository-1] transaction not found: %v", err)
			return nil, err
		}

		log.Errorf("[TransactionRepository-2] failed to find transaction by bill ID: %v", err)
		return nil, err
	}

	transactionModel := &model.TransactionModel{
		ID:          e.ID,
		MerchantID:  e.MerchantID,
		OrderID:     e.OrderID,
		BillID:      e.BillID,
		Amount:      e.Amount,
		PaymentType: e.PaymentType,
		Status:      e.Status,
		PaidAt:      e.PaidAt,
		ExpiredAt:   e.ExpiredAt,
	}

	return transactionModel, nil
}

// FindActivePendingTransaction implements [TransactionRepositoryInterface].
func (t *transactionRepository) FindActivePendingTransaction(merchantID string, billID string, paymentType string) (*model.TransactionModel, error) {
	var e entity.TransactionEntity

	err := t.db.Where("merchant_id = ? AND bill_id = ? AND payment_type = ? AND status = ? AND (expired_at IS NULL OR expired_at > NOW())",
		merchantID, billID, paymentType, string(enum.TransactionStatusPending)).Order("created_at DESC").Limit(1).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[TransactionRepository-2] active pending transaction not found: %v", err)
			return nil, err
		}

		log.Errorf("[TransactionRepository-2] failed to find active pending transaction: %v", err)
		return nil, err
	}

	transactionModel := &model.TransactionModel{
		ID:          e.ID,
		MerchantID:  e.MerchantID,
		OrderID:     e.OrderID,
		BillID:      e.BillID,
		Amount:      e.Amount,
		PaymentType: e.PaymentType,
		QrURL:       e.QrURL,
		Status:      e.Status,
		PaidAt:      e.PaidAt,
		ExpiredAt:   e.ExpiredAt,
	}

	return transactionModel, nil
}

// CreateTransaction implements [TransactionRepositoryInterface].
func (t *transactionRepository) CreateTransaction(
	scope enum.TransactionScope,
	merchantID *string,
	billID *string,
	orderID string,
	amount float64,
	paymentType string,
	qrURL string,
	expiredAt *time.Time,
) (*model.TransactionModel, error) {
	query := `
		INSERT INTO transactions
			(merchant_id, bill_id, order_id, amount, payment_type, status, qr_url, expired_at, transaction_scope)
		VALUES
			($1, $2, $3, $4, $5, 'pending', $6, $7, $8)
		RETURNING id
	`

	var e entity.TransactionEntity

	err := t.db.Raw(
		query,
		merchantID,
		billID,
		orderID,
		amount,
		paymentType,
		qrURL,
		expiredAt,
		scope,
	).Scan(&e.ID).Error

	if err != nil {
		log.Errorf("[TransactionRepository-1] failed to create transaction: %v", err)
		return nil, err
	}

	e.MerchantID = ""
	if merchantID != nil {
		e.MerchantID = *merchantID
	}

	e.OrderID = orderID
	e.Amount = amount
	e.PaymentType = paymentType
	e.Status = enum.TransactionStatusPending
	e.TransactionScope = scope
	e.ExpiredAt = expiredAt
	e.QrURL = &qrURL

	transactionModel := &model.TransactionModel{
		ID:               e.ID,
		MerchantID:       e.MerchantID,
		OrderID:          e.OrderID,
		Amount:           e.Amount,
		PaymentType:      e.PaymentType,
		Status:           e.Status,
		QrURL:            e.QrURL,
		ExpiredAt:        e.ExpiredAt,
		TransactionScope: e.TransactionScope,
	}

	return transactionModel, nil
}

// FindByOrderID implements [TransactionRepositoryInterface].
func (t *transactionRepository) FindByOrderID(orderID string) (*model.TransactionModel, error) {
	var e entity.TransactionEntity

	// Menggunakan Where().First() - otomatis return gorm.ErrRecordNotFound jika tidak ada data
	err := t.db.Where("order_id = ?", orderID).First(&e).Error
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
		ID:               e.ID,
		MerchantID:       e.MerchantID,
		OrderID:          e.OrderID,
		Amount:           e.Amount,
		PaymentType:      e.PaymentType,
		TransactionScope: e.TransactionScope,
		Status:           e.Status,
		PaidAt:           e.PaidAt,
		ExpiredAt:        e.ExpiredAt,
	}

	fmt.Println("[LOG] Mapped TransactionModel:", transactionModel)

	return transactionModel, nil
}

// UpdateStatus implements [TransactionRepositoryInterface].
func (t *transactionRepository) UpdateStatus(
	orderID string,
	status enum.TransactionStatus,
	paidAt *time.Time,
) error {
	query := `
		UPDATE transactions
		SET status = $1,
			paid_at = $2,
			updated_at = NOW()
		WHERE order_id = $3
		AND status != '` + string(enum.TransactionStatusPaid) + `'
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
