package service

import (
	"errors"
	"fmt"
	"sun-pos-payment-service/internal/adapter/repository"
	"sun-pos-payment-service/internal/core/domain/model"
	"sun-pos-payment-service/internal/core/domain/payment"
	"sun-pos-payment-service/utils/enum"
	"time"

	"github.com/labstack/gommon/log"
)

type PaymentServiceInterface interface {
	GenerateQRIS(input GenerateQRISInput) (*GenerateQRISResult, error)
	GenerateOwnerQRIS(input GenerateOwnerQRISInput) (*GenerateQRISResult, error)
}

type paymentService struct {
	midtransClient     payment.MidtransClientInterface
	transactionRepo    repository.TransactionRepositoryInterface
	merchantRepository repository.MerchantRepositoryInterface
	ownerServerKey     string
}

// GenerateOwnerQRIS implements [PaymentServiceInterface].
func (p *paymentService) GenerateOwnerQRIS(input GenerateOwnerQRISInput) (*GenerateQRISResult, error) {
	if input.OrderID == "" || input.Amount <= 0 {
		log.Errorf("[Payment Service-Owner-1] invalid order ID or amount")
		return nil, errors.New("invalid order ID or amount")
	}

	if p.ownerServerKey == "" {
		log.Errorf("[Payment Service-Owner-2] owner server key is not configured")
		return nil, errors.New("owner server key is not configured")
	}

	expMinutes := input.ExpireMinutes
	if expMinutes <= 0 {
		expMinutes = 15
	}

	mtRes, err := p.midtransClient.ChargeQris(
		p.ownerServerKey,
		payment.QrisChargeInput{
			OrderID: input.OrderID,
			Amount:  input.Amount,
		},
	)

	if err != nil {
		log.Errorf("[Payment Service-Owner-3] failed to charge QRIS: %v", err)
		return nil, err
	}

	expiredAt := mtRes.ExpiredAt
	if expiredAt == nil {
		t := time.Now().Add(time.Duration(expMinutes) * time.Minute)
		expiredAt = &t
	}

	_, err = p.transactionRepo.CreateTransaction(
		enum.ScopeOwner,
		nil,
		nil,
		input.OrderID,
		input.Amount,
		model.PaymentTypeQRIS,
		mtRes.QrURL,
		expiredAt,
	)

	if err != nil {
		log.Errorf("[Payment Service-Owner-4] failed to create transaction: %v", err)
		return nil, err
	}

	result := &GenerateQRISResult{
		OrderID:   mtRes.OrderID,
		QrURL:     mtRes.QrURL,
		ExpiredAt: expiredAt,
		Status:    enum.TransactionStatusPending,
	}

	return result, nil
}

// GenerateQRIS implements [PaymentServiceInterface].
func (p *paymentService) GenerateQRIS(
	input GenerateQRISInput,
) (*GenerateQRISResult, error) {
	if input.ServerKey == "" {
		log.Errorf("[Payment Service-1] server key is empty")
		return nil, errors.New("server key is required")
	}

	if input.BillID == "" || input.Amount <= 0 {
		log.Errorf("[Payment Service-2] invalid bill ID or amount")
		return nil, errors.New("invalid bill ID or amount")
	}

	merchant, err := p.merchantRepository.FindByID(input.MerchantID)

	if err != nil {
		if err.Error() == "404" {
			log.Errorf("[Payment Service-2] merchant not found: %v", err)
			merchant, err = p.merchantRepository.Create(
				input.MerchantID,
				input.Acquirer,
				"sandbox",
			)

			if err != nil {
				log.Errorf("[Payment Service-2] failed to create merchant: %v", err)
				return nil, err
			}

			if merchant.ServerKey == "" && input.ServerKey != "" {
				err = p.merchantRepository.SaveServerKey(merchant.ID, input.ServerKey)
				if err != nil {
					log.Errorf("[Payment Service-2] failed to save server key: %v", err)
					return nil, err
				}
				// Update merchant model with the server key
				merchant.ServerKey = input.ServerKey
			}
		} else {
			log.Errorf("[Payment Service-2] failed to find merchant: %v", err)
			return nil, err
		}
	}

	existingTransaction, err := p.transactionRepo.FindActivePendingTransaction(
		merchant.ID,
		input.BillID,
		model.PaymentTypeQRIS,
	)

	if err == nil && existingTransaction != nil {
		if existingTransaction.Amount == input.Amount {
			log.Infof("[Payment Service-3] reusing existing pending transaction for bill ID: %s", existingTransaction.BillID)

			return &GenerateQRISResult{
				OrderID:   existingTransaction.OrderID,
				QrURL:     *existingTransaction.QrURL,
				ExpiredAt: existingTransaction.ExpiredAt,
				Status:    existingTransaction.Status,
				BillID:    existingTransaction.BillID,
			}, nil
		}

		log.Infof("[Payment Service-3] existing pending transaction found with different amount, creating new transaction")
		_ = p.transactionRepo.UpdateStatus(existingTransaction.OrderID, enum.TransactionStatusFailed, nil)
	}

	paymentOrderID := fmt.Sprintf("QRIS-%s-%d", input.BillID, time.Now().UnixNano())

	expMinutes := input.ExpireMinutes
	if expMinutes <= 0 {
		expMinutes = 15
	}

	mtRes, err := p.midtransClient.ChargeQris(
		merchant.ServerKey,
		payment.QrisChargeInput{
			OrderID:  paymentOrderID,
			Amount:   input.Amount,
			Acquirer: input.Acquirer,
		},
	)

	if err != nil {
		log.Errorf("[Payment Service-3] failed to charge QRIS: %v", err)
		return nil, err
	}

	expiredAt := mtRes.ExpiredAt
	if expiredAt == nil {
		t := time.Now().Add(time.Duration(expMinutes) * time.Minute)
		expiredAt = &t
	}

	_, err = p.transactionRepo.CreateTransaction(
		enum.ScopeMerchant,
		&merchant.ID,
		&input.BillID,
		paymentOrderID,
		input.Amount,
		model.PaymentTypeQRIS,
		mtRes.QrURL,
		expiredAt,
	)

	if err != nil {
		log.Errorf("[Payment Service-4] failed to create transaction: %v", err)
		return nil, err
	}

	result := &GenerateQRISResult{
		OrderID:   mtRes.OrderID,
		QrURL:     mtRes.QrURL,
		ExpiredAt: expiredAt,
		BillID:    input.BillID,
		Status:    enum.TransactionStatusPending,
	}

	return result, nil
}

func NewPaymentService(
	midtransClient payment.MidtransClientInterface,
	transactionRepo repository.TransactionRepositoryInterface,
	merchantRepository repository.MerchantRepositoryInterface,
	ownerServerKey string,
) PaymentServiceInterface {
	return &paymentService{
		midtransClient:     midtransClient,
		transactionRepo:    transactionRepo,
		merchantRepository: merchantRepository,
		ownerServerKey:     ownerServerKey,
	}
}
