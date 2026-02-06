package seeds

import (
	"sun-pos-payment-service/config"
	"sun-pos-payment-service/internal/adapter/repository"
	"sun-pos-payment-service/internal/core/domain/model"
	"sun-pos-payment-service/internal/security"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

func SeedMerchants(db *gorm.DB) error {
	cfg := config.NewConfig()

	hexSecret := cfg.Security.EncryptionSecret

	encryptor, err := security.NewEncryptor(hexSecret)
	if err != nil {
		log.Errorf("[Merchant Seed-1] failed to create encryptor: %v", err)
		return err
	}

	merchantRepo := repository.NewMerchantRepository(db, encryptor)

	serverKey := "SB-Mid-server-z_8NezgZ9T6wI3iB9N9etDDD"
	encryptedServerKey, err := encryptor.Encrypt(serverKey)
	if err != nil {
		log.Errorf("[Merchant Seed-2] failed to encrypt server key: %v", err)
		return err
	}

	merchants := []model.MerchantModel{
		{
			Name:           "Demo Merchant",
			ServerKey:      encryptedServerKey,
			KeyEnvironment: "sandbox",
		},
	}

	for _, m := range merchants {
		_, err := merchantRepo.Create( "", m.Name, m.KeyEnvironment)
		if err != nil {
			return err
		}
	}

	return nil
}
