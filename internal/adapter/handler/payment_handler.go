package handler

import (
	"net/http"
	"sun-pos-payment-service/internal/adapter/dto/request"
	"sun-pos-payment-service/internal/adapter/dto/response"
	"sun-pos-payment-service/internal/core/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type PaymentHandlerInterface interface {
	GenerateQris(c echo.Context) error
	GenerateOwnerQris(c echo.Context) error
}

type paymentHandler struct {
	paymentService service.PaymentServiceInterface
}

// GenerateOwnerQris implements [PaymentHandlerInterface].
func (p *paymentHandler) GenerateOwnerQris(c echo.Context) error {
	var (
		request         request.GenerateOwnerQrisRequest
		defaultResponse response.DefaultResponse
	)

	if err := c.Bind(&request); err != nil {
		log.Errorf("[Payment Handler-Owner-1] failed to bind request: %v", err)
		defaultResponse.Message = "invalid request payload"
		defaultResponse.Data = nil
		return c.JSON(http.StatusBadRequest, defaultResponse)
	}

	if err := c.Validate(&request); err != nil {
		log.Errorf("[Payment Handler-Owner-2] validation failed: %v", err)
		defaultResponse.Message = "validation failed: " + err.Error()
		defaultResponse.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, defaultResponse)
	}

	result, err := p.paymentService.GenerateOwnerQRIS(
		service.GenerateOwnerQRISInput{
			OrderID:       request.OrderID,
			Amount:        request.Amount,
			Acquirer:      request.Acquirer,
			ExpireMinutes: request.ExpireMinutes,
		},
	)

	if err != nil {
		defaultResponse.Message = "failed to generate QRIS: " + err.Error()
		defaultResponse.Data = nil

		if err.Error() == "404" {
			log.Infof("[Payment Handler-Owner-3] merchant not found: %v", err)
			return c.JSON(http.StatusNotFound, defaultResponse)
		}

		log.Errorf("[Payment Handler-Owner-4] failed to generate QRIS: %v", err)
		return c.JSON(http.StatusInternalServerError, defaultResponse)
	}

	defaultResponse.Message = "QRIS generated successfully"
	defaultResponse.Data = response.GenerateQrisResponse{
		OrderID:   result.OrderID,
		QrUrl:     result.QrURL,
		ExpiredAt: result.ExpiredAt,
		Status:    result.Status,
	}

	return c.JSON(http.StatusOK, defaultResponse)
}

// GenerateQris implements [PaymentHandlerInterface].
func (p *paymentHandler) GenerateQris(c echo.Context) error {
	var (
		req             request.GenerateQrisRequest
		defaultResponse response.DefaultResponse
		dataResponse    response.GenerateQrisResponse
	)

	if err := c.Bind(&req); err != nil {
		log.Errorf("[Payment Handler-1] failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "invalid request payload",
			Data:    nil,
		})
	}

	if err := c.Validate(&req); err != nil {
		log.Errorf("[Payment Handler-4] validation failed: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.DefaultResponse{
			Message: "validation failed: " + err.Error(),
			Data:    nil,
		})
	}

	if req.MerchantID == "" {
		log.Infof("[Payment Handler-5] Merchant ID is required")
		return c.JSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "merchant id is required",
			Data:    nil,
		})
	}

	result, err := p.paymentService.GenerateQRIS(
		service.GenerateQRISInput{
			MerchantID: req.MerchantID,
			ServerKey:  req.ServerKey,
			BillID:     req.OrderID,
			Amount:     req.Amount,
			Acquirer:   req.Acquirer,
		},
	)

	if err != nil {
		if err.Error() == "404" {
			log.Infof("[Payment Handler-3] merchant not found: %v", err)
			return c.JSON(http.StatusNotFound, response.DefaultResponse{
				Message: "merchant not found",
				Data:    nil,
			})
		}

		log.Errorf("[Payment Handler-2] failed to generate QRIS: %v", err)
		return c.JSON(http.StatusInternalServerError, response.DefaultResponse{
			Message: "failed to generate QRIS",
			Data:    nil,
		})
	}

	dataResponse = response.GenerateQrisResponse{
		OrderID:   result.OrderID,
		QrUrl:     result.QrURL,
		BillID:    result.BillID,
		ExpiredAt: result.ExpiredAt,
		Status:    result.Status,
	}
	msg := "QRIS generated successfully"
	defaultResponse = response.DefaultResponse{
		Message: msg,
		Data:    dataResponse,
	}

	return c.JSON(http.StatusOK, defaultResponse)
}

func NewPaymentHandler(
	paymentService service.PaymentServiceInterface,
) PaymentHandlerInterface {
	return &paymentHandler{
		paymentService: paymentService,
	}
}
