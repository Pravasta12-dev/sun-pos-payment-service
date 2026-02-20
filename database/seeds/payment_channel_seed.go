package seeds

import (
	"sun-pos-payment-service/internal/core/domain/entity"
	"sun-pos-payment-service/utils/enum"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

func SeedPaymentChannels(db *gorm.DB) error {
	var count int64
	err := db.Model(&entity.PaymentChannelEntity{}).Count(&count).Error
	if err != nil {
		log.Errorf("[Payment Channel Seed-1] failed to count payment channels: %v", err)
		return err
	}

	if count > 0 {
		log.Infof("[Payment Channel Seed-1] payment channels already seeded, skipping seeding")
		return nil
	}

	paymentChannels := []entity.PaymentChannelEntity{
		{
			Type:     enum.PaymentChannelTypeQris,
			Code:     enum.PaymentChannelCodeQris,
			Label:    "QRIS",
			FeeType:  "percentage",
			FeeValue: 0.8,
			IsActive: true,
		},
		{
			Type:     enum.PaymentChannelTypeVa,
			Code:     enum.PaymentChannelCodeBca,
			Label:    "BCA Virtual Account",
			FeeType:  "percentage",
			FeeValue: 0.5,
			IsActive: true,
		},
		{
			Type:     enum.PaymentChannelTypeVa,
			Code:     enum.PaymentChannelCodeBri,
			Label:    "BRI Virtual Account",
			FeeType:  "percentage",
			FeeValue: 0.5,
			IsActive: true,
		},
		{
			Type:     enum.PaymentChannelTypeVa,
			Code:     enum.PaymentChannelCodeBni,
			Label:    "BNI Virtual Account",
			FeeType:  "percentage",
			FeeValue: 0.5,
			IsActive: true,
		},
	}

	for _, pc := range paymentChannels {
		err := db.Create(&pc).Error
		if err != nil {
			log.Errorf("[Payment Channel Seed-2] failed to create payment channel: %v", err)
			return err
		}
	}

	log.Infof("[Payment Channel Seed-3] successfully seeded payment channels")
	return nil
}
