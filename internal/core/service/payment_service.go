package service

import (
	"errors"
	"sun-pos-payment-service/internal/adapter/repository"
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
	merchantRepository repository.MerchantRepositoryInterface
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

	expMinutes := input.ExpireMinutes
	if expMinutes <= 0 {
		expMinutes = 15
	}

	mtRes, err := p.midtransClient.ChargeQris(
		merchant.ServerKey,
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
		merchant.ID,
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
	merchantRepository repository.MerchantRepositoryInterface,
) PaymentServiceInterface {
	return &paymentService{
		midtransClient:     midtransClient,
		transactionService: transactionService,
		merchantRepository: merchantRepository,
	}
}
