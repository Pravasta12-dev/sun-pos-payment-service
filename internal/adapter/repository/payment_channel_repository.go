package repository

import (
	"errors"
	"sun-pos-payment-service/internal/core/domain/entity"
	"sun-pos-payment-service/internal/core/domain/model"

	"gorm.io/gorm"

	"github.com/labstack/gommon/log"
)

type PaymentChannelRepository interface {
	GetActivePaymentChannels() ([]*model.PaymentChannelModel, error)
}

type paymentChannelRepository struct {
	db *gorm.DB
}

// GetActivePaymentChannels implements [PaymentChannelRepository].
func (p *paymentChannelRepository) GetActivePaymentChannels() ([]*model.PaymentChannelModel, error) {
	var e []entity.PaymentChannelEntity

	err := p.db.Where("is_active = ?", true).Order("created_at ASC").Find(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[PaymentChannelRepository-1] no active payment channels found: %v", err)
			return nil, err
		}

		log.Errorf("[PaymentChannelRepository-2] failed to get active payment channels: %v", err)
		return nil, err
	}

	var m []*model.PaymentChannelModel
	for _, v := range e {
		m = append(m, &model.PaymentChannelModel{
			ID:       v.ID,
			Type:     v.Type,
			Code:     v.Code,
			Label:    v.Label,
			FeeType:  v.FeeType,
			FeeValue: v.FeeValue,
		})
	}
	return m, nil
}

func NewPaymentChannelRepository(db *gorm.DB) PaymentChannelRepository {
	return &paymentChannelRepository{db: db}
}
