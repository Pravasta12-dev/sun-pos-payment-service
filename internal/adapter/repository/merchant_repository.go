package repository

import (
	"errors"
	"sun-pos-payment-service/internal/core/domain/entity"
	"sun-pos-payment-service/internal/core/domain/model"
	"sun-pos-payment-service/internal/security"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type MerchantRepositoryInterface interface {
	Create(name, env string) (*model.MerchantModel, error)
	SaveServerKey(merchantID int64, serverKey string) error
	FindByID(merchantID int64) (*model.MerchantModel, error)
}

type merchantRepository struct {
	db        *gorm.DB
	encryptor security.EncryptorInterface
}

// Create implements [MerchantRepositoryInterface].
func (m *merchantRepository) Create(name, env string) (*model.MerchantModel, error) {
	query := `
		INSERT INTO merchants (name, key_environment)
		VALUES ($1, $2)
		RETURNING id
	`

	var e entity.MerchantEntity

	err := m.db.Raw(query, name, env).Scan(&e).Error
	if err != nil {
		log.Errorf("[Merchant Repository-1] failed to create merchant: %v", err)
		return nil, err
	}

	e.Name = name
	e.KeyEnvirontment = env

	merchantModel := &model.MerchantModel{
		ID:             e.ID,
		Name:           e.Name,
		KeyEnvironment: e.KeyEnvirontment,
	}

	return merchantModel, nil
}

// FindByID implements [MerchantRepositoryInterface].
func (m *merchantRepository) FindByID(merchantID int64) (*model.MerchantModel, error) {

	var e entity.MerchantEntity

	if err := m.db.Where("id = ?", merchantID).First(&e).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[Merchant Repository-2] merchant not found: %v", err)
			return nil, err
		}

		log.Errorf("[Merchant Repository-2] failed to find merchant by ID: %v", err)
		return nil, err
	}

	merchantModel := model.MerchantModel{
		ID:             e.ID,
		Name:           e.Name,
		KeyEnvironment: e.KeyEnvirontment,
	}

	// Decrypt server key if exists
	if e.ServerKey != nil && *e.ServerKey != "" {
		decryptedKey, err := m.encryptor.Decrypt(*e.ServerKey)
		if err != nil {
			log.Errorf("[Merchant Repository-3] failed to decrypt server key: %v", err)
			return nil, err
		}
		merchantModel.ServerKey = decryptedKey
	}

	return &merchantModel, nil
}

// SaveServerKey implements [MerchantRepositoryInterface].
func (m *merchantRepository) SaveServerKey(merchantID int64, serverKey string) error {
	encryptedKey, err := m.encryptor.Encrypt(serverKey)
	if err != nil {
		log.Errorf("[Merchant Repository-4] failed to encrypt server key: %v", err)
		return err
	}

	query := `
		UPDATE merchants
		SET server_key = $1, updated_at = NOW()
		WHERE id = $2
	`

	if err := m.db.Exec(query, encryptedKey, merchantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[Merchant Repository-5] merchant not found when saving server key: %v", err)
			return err
		}
		log.Errorf("[Merchant Repository-6] failed to save server key: %v", err)
		return err
	}

	return nil
}

func NewMerchantRepository(
	db *gorm.DB,
	encryptor security.EncryptorInterface,
) MerchantRepositoryInterface {
	return &merchantRepository{
		db:        db,
		encryptor: encryptor,
	}
}
