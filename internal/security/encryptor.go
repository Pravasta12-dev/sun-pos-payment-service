package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"github.com/labstack/gommon/log"
)

type EncryptorInterface interface {
	Encrypt(plainText string) (string, error)
	Decrypt(cipherText string) (string, error)
}

type encryptor struct {
	key []byte
}

// Decrypt implements [EncryptorInterface].
func (e *encryptor) Decrypt(cipherText string) (string, error) {
	data, err := hex.DecodeString(cipherText)

	if err != nil {
		log.Errorf("[Encryptor-6] Failed to decode hex string: %v", err)
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		log.Errorf("[Encryptor-7] Failed to create cipher block: %v", err)
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Errorf("[Encryptor-8] Failed to create GCM: %v", err)
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		log.Errorf("[Encryptor-9] Ciphertext too short")
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherTextBytes := data[:nonceSize], data[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		log.Errorf("[Encryptor-10] Failed to decrypt data: %v", err)
		return "", err
	}

	return string(plainText), nil
}

// Encrypt implements [EncryptorInterface].
func (e *encryptor) Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		log.Errorf("[Encryptor-1] Failed to create cipher block: %v", err)
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Errorf("[Encryptor-2] Failed to create GCM: %v", err)
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Errorf("[Encryptor-3] Failed to generate nonce: %v", err)
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	return hex.EncodeToString(cipherText), nil
}

func NewEncryptor(hexSecret string) (EncryptorInterface, error) {
	key, err := hex.DecodeString(hexSecret)
	if err != nil {
		log.Errorf("[Encryptor-4] Encryptor must be hex encoded")
		return nil, errors.New("Encryption Secret Must be hex encoded")
	}

	if len(key) != 32 {
		log.Errorf("[Encryptor-5] Encryptor key length must be 32 bytes")
		return nil, errors.New("Encryption Secret length must be 32 bytes")
	}

	return &encryptor{
		key: key,
	}, nil
}
