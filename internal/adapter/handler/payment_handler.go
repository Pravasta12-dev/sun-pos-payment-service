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
}

type paymentHandler struct {
	paymentService service.PaymentServiceInterface
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

	result, err := p.paymentService.GenerateQRIS(
		service.GenerateQRISInput{
			MerchantID: req.MerchantID,
			ServerKey:  req.ServerKey,
			OrderID:    req.OrderID,
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
