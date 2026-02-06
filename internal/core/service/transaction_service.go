package service

import (
	"sun-pos-payment-service/internal/adapter/repository"
	"sun-pos-payment-service/internal/core/domain/model"
	"time"

	"github.com/labstack/gommon/log"
)

type TransactionServiceInterface interface {
	CreateTransaction(
		merchantID string,
		orderID string,
		amount float64,
		paymentType string,
		expiredAt *time.Time,
	) (*model.TransactionModel, error)
	MarkAsPaid(
		orderID string,
		paidAt *time.Time,
	) error
	MarkAsExpired(orderID string) error
	MarkAsFailed(orderID string) error
	GetByOrderID(orderID string) (*model.TransactionModel, error)
}

type transactionService struct {
	transactionRepo repository.TransactionRepositoryInterface
}

// CreateTransaction implements [TransactionServiceInterface].
func (t *transactionService) CreateTransaction(merchantID string, orderID string, amount float64, paymentType string, expiredAt *time.Time) (*model.TransactionModel, error) {
	return t.transactionRepo.CreateTransaction(
		merchantID,
		orderID,
		amount,
		paymentType,
		expiredAt,
	)
}

// GetByOrderID implements [TransactionServiceInterface].
func (t *transactionService) GetByOrderID(orderID string) (*model.TransactionModel, error) {
	tx, err := t.transactionRepo.FindByOrderID(orderID)

	if err != nil {
		log.Errorf("[TransactionService-4] failed to get transaction by order ID: %v", err)
		return nil, err
	}

	return tx, nil
}

// MarkAsExpired implements [TransactionServiceInterface].
func (t *transactionService) MarkAsExpired(orderID string) error {
	if err := t.transactionRepo.UpdateStatus(
		orderID,
		model.TransactionStatusExpired,
		nil,
	); err != nil {
		log.Errorf("[TransactionService-2] failed to mark transaction as expired: %v", err)
		return err
	}

	return nil
}

// MarkAsFailed implements [TransactionServiceInterface].
func (t *transactionService) MarkAsFailed(orderID string) error {
	if err := t.transactionRepo.UpdateStatus(
		orderID,
		model.TransactionStatusFailed,
		nil,
	); err != nil {
		log.Errorf("[TransactionService-3] failed to mark transaction as failed: %v", err)
		return err
	}

	return nil
}

// MarkAsPaid implements [TransactionServiceInterface].
func (t *transactionService) MarkAsPaid(orderID string, paidAt *time.Time) error {
	if paidAt == nil {
		now := time.Now()
		paidAt = &now
	}

	if err := t.transactionRepo.UpdateStatus(
		orderID,
		model.TransactionStatusPaid,
		paidAt,
	); err != nil {
		log.Errorf("[TransactionService-1] failed to mark transaction as paid: %v", err)
		return err
	}

	return nil

}

func NewTransactionService(
	transactionRepo repository.TransactionRepositoryInterface,
) TransactionServiceInterface {
	return &transactionService{
		transactionRepo: transactionRepo,
	}
}
