package service

import (
	"errors"
	"sun-pos-payment-service/internal/core/domain/model"
	"sun-pos-payment-service/internal/core/domain/payment"
	"time"

	"github.com/labstack/gommon/log"
)

type PaymentServiceInterface interface {
	GenerateQRIS(input GenerateQRISInput) (*GenerateQRISResult, error)
}

type paymentService struct {
	midtransClient     payment.MidtransClientInterface
	transactionService TransactionServiceInterface
}

// GenerateQRIS implements [PaymentServiceInterface].
func (p *paymentService) GenerateQRIS(
	input GenerateQRISInput,
) (*GenerateQRISResult, error) {
	if input.ServerKey == "" {
		log.Errorf("[Payment Service-1] server key is empty")
		return nil, errors.New("server key is required")
	}

	if input.OrderID == "" || input.Amount <= 0 {
		log.Errorf("[Payment Service-2] invalid order ID or amount")
		return nil, errors.New("invalid order ID or amount")
	}

	expMinutes := input.ExpireMinutes
	if expMinutes <= 0 {
		expMinutes = 15
	}

	mtRes, err := p.midtransClient.ChargeQris(
		input.ServerKey,
		payment.QrisChargeInput{
			OrderID:  input.OrderID,
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

	_, err = p.transactionService.CreateTransaction(
		input.MerchantID,
		input.OrderID,
		input.Amount,
		model.PaymentTypeQRIS,
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
		Status:    model.TransactionStatusPending,
	}

	return result, nil
}

func NewPaymentService(
	midtransClient payment.MidtransClientInterface,
	transactionService TransactionServiceInterface,
) PaymentServiceInterface {
	return &paymentService{
		midtransClient:     midtransClient,
		transactionService: transactionService,
	}
}
