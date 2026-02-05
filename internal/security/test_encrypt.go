package security

import (
	"log"
	"sun-pos-payment-service/config"

	"github.com/spf13/viper"
)

func TestEncryptDecrypt() {
	// Load config file first
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	// encryptor interface
	cfg := config.NewConfig()

	encryptor, err := NewEncryptor(
		cfg.Security.EncryptionSecret,
	)

	if err != nil {
		log.Fatalf("failed to create encryptor: %v", err)
	}

	originalText := "Hello, World!"
	encryptedText, err := encryptor.Encrypt(originalText)
	if err != nil {
		log.Fatalf("encryption failed: %v", err)
	}

	decryptedText, err := encryptor.Decrypt(encryptedText)
	if err != nil {
		log.Fatalf("decryption failed: %v", err)
	}

	if decryptedText != originalText {
		log.Fatalf("decrypted text does not match original. got: %s, want: %s", decryptedText, originalText)
	}

	log.Println("Encryption and decryption successful")
	log.Println("ENCRYPTED", encryptedText)
	log.Println("DECRYPTED", decryptedText)
}
