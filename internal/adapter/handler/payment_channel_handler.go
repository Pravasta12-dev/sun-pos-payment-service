package handler

import (
	"sun-pos-payment-service/internal/adapter/dto/response"
	"sun-pos-payment-service/internal/core/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type PaymentChannelHandlerInterface interface {
	GetChannels(c echo.Context) error
}

type paymentChannelHandler struct {
	paymentChannelService service.PaymentChannelService
}

// GetChannels implements [PaymentChannelHandlerInterface].
func (p *paymentChannelHandler) GetChannels(c echo.Context) error {
	var (
		defaultResponse response.DefaultResponse
		resultResponse  []response.PaymentChannelResponse
	)

	result, err := p.paymentChannelService.GetChannels()
	if err != nil {
		defaultResponse.Message = "failed to get payment channels: " + err.Error()
		defaultResponse.Data = nil

		if err.Error() == "404" {
			log.Infof("[Payment Channel Handler-1] no active payment channels found: %v", err)
			return c.JSON(404, defaultResponse)
		}

		log.Errorf("[Payment Channel Handler-2] failed to get payment channels: %v", err)
		return c.JSON(500, defaultResponse)
	}

	for _, channel := range result {
		resultResponse = append(resultResponse, response.PaymentChannelResponse{
			ID:       channel.ID,
			Type:     channel.Type,
			Code:     channel.Code,
			Label:    channel.Label,
			FeeType:  channel.FeeType,
			FeeValue: channel.FeeValue,
		})
	}

	defaultResponse.Message = "success"
	defaultResponse.Data = resultResponse

	return c.JSON(200, defaultResponse)
}

func NewPaymentChannelHandler(paymentChannelService service.PaymentChannelService) PaymentChannelHandlerInterface {
	return &paymentChannelHandler{paymentChannelService: paymentChannelService}
}
