package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sun-pos-payment-service/config"
	"sun-pos-payment-service/internal/adapter/handler"
	"sun-pos-payment-service/internal/adapter/payment"
	"sun-pos-payment-service/internal/adapter/repository"
	"sun-pos-payment-service/internal/core/service"
	"sun-pos-payment-service/internal/routes"
	"sun-pos-payment-service/internal/security"
	"sun-pos-payment-service/utils/validator"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
)

func RunServer() {
	cfg := config.NewConfig()

	db, err := cfg.ConnectionPostgres()

	if err != nil {
		log.Fatalf("[Run Server - 1] Failed to connect to database %v", err)
		return
	}

	hexSecret := cfg.Security.EncryptionSecret

	encryptor, err := security.NewEncryptor(hexSecret)
	if err != nil {
		log.Fatalf("[Run Server - 1] Failed to create encryptor %v", err)
		return
	}

	transactionRepository := repository.NewTransactionRepository(db.DB)
	merchantRepository := repository.NewMerchantRepository(db.DB, encryptor)

	transactionService := service.NewTransactionService(transactionRepository)
	midtransClient := payment.NewMidtransClient(cfg.Midtrans.BaseURL)
	paymentService := service.NewPaymentService(
		midtransClient,
		transactionRepository,
		merchantRepository,
		cfg.Midtrans.ServerKey,
	)

	e := echo.New()

	// Register custom validator
	customValidator := validator.NewValidator()
	// Register English translations
	en.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator)
	// Set the custom validator to Echo instance
	e.Validator = customValidator

	paymentHandler := handler.NewPaymentHandler(paymentService)
	webhookHandler := handler.NewMidtransWebhookHandler(transactionService, merchantRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// // seed
	// seeds.SeedMerchants(db.DB)

	routes.RegisterRoutes(
		e,
		paymentHandler,
		webhookHandler,
		transactionHandler,
	)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := e.Start(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatalf("[Run Server - 2] Failed to start server %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	e.Shutdown(ctx)
}
