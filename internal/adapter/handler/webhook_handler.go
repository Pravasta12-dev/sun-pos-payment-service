package handler

import (
	"fmt"
	"net/http"
	"sun-pos-payment-service/internal/adapter/dto/request"
	"sun-pos-payment-service/internal/adapter/dto/response"
	"sun-pos-payment-service/internal/adapter/repository"
	"sun-pos-payment-service/internal/core/service"
	"sun-pos-payment-service/internal/security"
	"sun-pos-payment-service/utils/enum"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type MitransWebhookHandlerInterface interface {
	HandleWebhook(e echo.Context) error
}

type midtransWebhook struct {
	transactionService service.TransactionServiceInterface
	merchantRepository repository.MerchantRepositoryInterface
	ownerServerKey     string
}

// HandleWebhook implements MitransWebhookHandlerInterface.
func (m *midtransWebhook) HandleWebhook(e echo.Context) error {
	var (
		req         request.MidtransWebhookRequest
		defaultResp response.DefaultResponse
	)

	if err := e.Bind(&req); err != nil {
		log.Errorf("[Midtrans Webhook Handler] failed to bind request: %v", err)
		return e.JSON(http.StatusBadRequest, response.DefaultResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	tx, err := m.transactionService.GetByOrderID(req.OrderID)
	if err != nil {
		log.Errorf("[Midtrans Webhook Handler] failed to get transaction by order ID: %v", err)
		return e.JSON(http.StatusInternalServerError, response.DefaultResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	var serverKey string

	if tx.TransactionScope == enum.ScopeOwner {
		serverKey = m.ownerServerKey
	} else {
		merchant, err := m.merchantRepository.FindByID(tx.MerchantID)
		if err != nil {
			log.Errorf("[Midtrans Webhook Handler] failed to find merchant: %v", err)
			return e.JSON(http.StatusInternalServerError, response.DefaultResponse{
				Message: err.Error(),
				Data:    nil,
			})
		}
		serverKey = merchant.ServerKey
	}

	if serverKey == "" {
		log.Errorf("[Midtrans Webhook Handler] server key is not configured for order ID: %s", req.OrderID)
		return e.JSON(http.StatusInternalServerError, response.DefaultResponse{
			Message: "Merchant server key is not configured",
			Data:    nil,
		})
	}

	valid := security.ValidateMidtransSignature(
		req.OrderID,
		req.StatusCode,
		req.GrossAmount,
		serverKey,
		req.SignatureKey,
	)

	if !valid {
		log.Errorf("[Midtrans Webhook Handler] invalid signature for order ID: %s", req.OrderID)
		return e.JSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "Invalid signature",
			Data:    nil,
		})
	}

	fmt.Println("[Midtrans Webhook Handler] Valid signature confirmed for order ID:", req.OrderID)
	fmt.Println("[Midtrans Webhook Handler] Status Code:", req.StatusCode)
	fmt.Println("[Midtrans Webhook Handler] Transaction Status:", req.TransactionStatus)

	switch req.TransactionStatus {
	case "settlement", "capture":
		var paidAt *time.Time

		if req.SettlementTime != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", req.SettlementTime); err == nil {
				paidAt = &t
			}
		}

		if err := m.transactionService.MarkAsPaid(
			req.OrderID,
			paidAt,
		); err != nil {
			log.Errorf("[Midtrans Webhook Handler] failed to mark transaction as paid: %v", err)
			return e.JSON(http.StatusInternalServerError, response.DefaultResponse{
				Message: err.Error(),
				Data:    nil,
			})
		}

	case "expire":
		if err := m.transactionService.MarkAsExpired(
			req.OrderID,
		); err != nil {
			log.Errorf("[Midtrans Webhook Handler] failed to mark transaction as expired: %v", err)
			return e.JSON(http.StatusInternalServerError, response.DefaultResponse{
				Message: err.Error(),
				Data:    nil,
			})
		}
	case "cancel", "deny", "failure":
		if err := m.transactionService.MarkAsFailed(
			req.OrderID,
		); err != nil {
			log.Errorf("[Midtrans Webhook Handler] failed to mark transaction as cancelled: %v", err)
			return e.JSON(http.StatusInternalServerError, response.DefaultResponse{
				Message: err.Error(),
				Data:    nil,
			})
		}
	default:
		log.Warnf("[Midtrans Webhook Handler] unhandled transaction status: %s", req.TransactionStatus)
	}

	defaultResp.Message = "Webhook processed successfully"
	defaultResp.Data = nil

	return e.JSON(http.StatusOK, defaultResp)
}

func NewMidtransWebhookHandler(
	transactionService service.TransactionServiceInterface,
	merchantRepository repository.MerchantRepositoryInterface,
	ownerServerKey string,
) MitransWebhookHandlerInterface {
	return &midtransWebhook{
		transactionService: transactionService,
		merchantRepository: merchantRepository,
		ownerServerKey:     ownerServerKey,
	}
}
