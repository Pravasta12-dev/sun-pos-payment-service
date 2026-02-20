package service

import (
	"sun-pos-payment-service/internal/adapter/repository"
	"sun-pos-payment-service/internal/core/domain/model"

	"github.com/labstack/gommon/log"
)

type PaymentChannelService interface {
	GetChannels() ([]*model.PaymentChannelModel, error)
}

type paymentChannelService struct {
	paymentChannelRepo repository.PaymentChannelRepository
}

// GetChannels implements [PaymentChannelService].
func (p *paymentChannelService) GetChannels() ([]*model.PaymentChannelModel, error) {
	pc, err := p.paymentChannelRepo.GetActivePaymentChannels()
	if err != nil {
		log.Errorf("[PaymentChannelService-1] failed to get active payment channels: %v", err)
		return nil, err
	}

	return pc, nil
}

func NewPaymentChannelService(paymentChannelRepo repository.PaymentChannelRepository) PaymentChannelService {
	return &paymentChannelService{paymentChannelRepo: paymentChannelRepo}
}
