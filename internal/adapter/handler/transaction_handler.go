package handler

import (
	"net/http"
	"sun-pos-payment-service/internal/adapter/dto/response"
	"sun-pos-payment-service/internal/core/service"

	"github.com/labstack/echo/v4"
)

type TransactionHandlerInterface interface {
	GetByOrderID(c echo.Context) error
}

type transactionHandler struct {
	transactionService service.TransactionServiceInterface
}

// GetByOrderID implements [TransactionHandlerInterface].
func (t *transactionHandler) GetByOrderID(c echo.Context) error {
	var (
		req         string = c.Param("order_id")
		defaultResp response.DefaultResponse
		dataResp    response.TransactionResponse
	)

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.DefaultResponse{
			Message: "invalid request",
			Data:    nil,
		})
	}

	merchantID := c.Request().Header.Get("X-Merchant-ID")
	if merchantID == "" {
		return c.JSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "Missing Merchant ID in header",
			Data:    nil,
		})
	}

	tx, err := t.transactionService.GetByBillID(merchantID, req)

	if err != nil {
		if err.Error() == "404" {
			return c.JSON(http.StatusNotFound, response.DefaultResponse{
				Message: "transaction not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, response.DefaultResponse{
			Message: "failed to get transaction",
			Data:    nil,
		})
	}

	dataResp = response.TransactionResponse{
		ID:          tx.ID,
		MerchantID:  tx.MerchantID,
		OrderID:     tx.OrderID,
		Amount:      tx.Amount,
		Status:      tx.Status,
		PaymentType: tx.PaymentType,
		PaidAt:      tx.PaidAt,
		ExpiredAt:   tx.ExpiredAt,
	}

	msg := "transaction retrieved successfully"
	defaultResp = response.DefaultResponse{
		Message: msg,
		Data:    dataResp,
	}

	return c.JSON(http.StatusOK, defaultResp)
}

func NewTransactionHandler(
	transactionService service.TransactionServiceInterface,
) TransactionHandlerInterface {
	return &transactionHandler{
		transactionService: transactionService,
	}
}
