package routes

import (
	"sun-pos-payment-service/internal/adapter/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	e *echo.Echo,
	paymentHandler handler.PaymentHandlerInterface,
	webhookHandler handler.MitransWebhookHandlerInterface,
	transactionHandler handler.TransactionHandlerInterface,
) {
	api := e.Group("/api")
	payment := api.Group("/payment")
	payment.POST("/webhook", webhookHandler.HandleWebhook)

	merchant := payment.Group("/merchant")
	merchant.POST("/generate-qris", paymentHandler.GenerateQris)
	merchant.GET("/transaction/:order_id", transactionHandler.GetByOrderID)

	owner := payment.Group("/owner")
	owner.POST("/generate-qris", paymentHandler.GenerateOwnerQris)
	owner.POST("/generate-va", paymentHandler.GenerateOwnerVA)
	owner.GET("/transaction/:bill_id", transactionHandler.GetByOwnerOrderID)

	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})
}
