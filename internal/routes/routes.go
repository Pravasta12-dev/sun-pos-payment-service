package routes

import (
	"sun-pos-payment-service/internal/adapter/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	e *echo.Echo,
	paymentHandler handler.PaymentHandlerInterface,
	webhookHandler handler.MitransWebhookHandlerInterface,
) {
	api := e.Group("/api")
	payment := api.Group("/payment")

	payment.POST("/generate-qris", paymentHandler.GenerateQris)
	payment.POST("/webhook", webhookHandler.HandleWebhook)


	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})
}
