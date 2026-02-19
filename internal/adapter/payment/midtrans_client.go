package payment

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sun-pos-payment-service/internal/adapter/dto/request"
	"sun-pos-payment-service/internal/adapter/dto/response"
	domain "sun-pos-payment-service/internal/core/domain/payment"
	"time"

	"github.com/labstack/gommon/log"
)

type midtransClient struct {
	baseURL string
	client  *http.Client
}

// ChargeVa implements [payment.MidtransClientInterface].
func (m *midtransClient) ChargeVa(serverKey string, input domain.VaChargeInput) (*domain.VaChargeResult, error) {
	reqBody := request.MidtransChargeRequest{
		PaymentType: "bank_transfer",
	}

	reqBody.TransactionDetails.OrderID = input.OrderID
	reqBody.TransactionDetails.GrossAmount = input.Amount
	reqBody.BankTransfer = &request.BankTransferRequest{
		Bank: input.Bank,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		log.Errorf("[MidtransClient-1] failed to marshal request body: %v", err)
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))

	req, err := http.NewRequest(
		http.MethodPost,
		m.baseURL+"/v2/charge",
		bytes.NewBuffer(payload),
	)

	fmt.Println("[MidtransClient] Request URL:", m.baseURL+"/v2/charge")

	if err != nil {
		log.Errorf("[MidtransClient-2] failed to create new request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		log.Errorf("[MidtransClient-3] failed to perform request: %v", err)
		return nil, err
	}

	fmt.Println("[MidtransClient] Request Data : ", string(payload))

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Errorf("[MidtransClient-4] received non-2xx response: %d", resp.StatusCode)
		return nil, fmt.Errorf("midtrans returned status code: %d", resp.StatusCode)
	}

	fmt.Println("[MidtransClient] resp data : ", resp)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[MidtransClient-5] failed to read response body: %v", err)
		return nil, err
	}

	fmt.Println("[MidtransClient] Raw Response Body:", string(bodyBytes))

	var result response.MidtransChargeResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		log.Errorf("[MidtransClient-6] failed to unmarshal response body: %v", err)
		return nil, err
	}

	fmt.Printf("[MidtransClient] Charge VA Response: %+v\n", result)

	if len(result.VaNumbers) == 0 {
		log.Errorf("[MidtransClient-6] va number not found in response")
		return nil, errors.New("va number not found in response")
	}

	return &domain.VaChargeResult{
		OrderID:  result.OrderID,
		VANumber: result.VaNumbers[0].VANumber,
		Bank:     result.VaNumbers[0].Bank,
	}, nil
}

// ChargeQris implements [payment.MidtransClientInterface].
func (m *midtransClient) ChargeQris(
	serverKey string,
	input domain.QrisChargeInput,
) (*domain.QrisChargeResult, error) {
	reqBody := request.MidtransChargeRequest{
		PaymentType: "qris",
	}

	if serverKey == "" {
		log.Errorf("[MidtransClient] server key is empty")
		return nil, errors.New("server key is required")
	}

	reqBody.TransactionDetails.OrderID = input.OrderID
	reqBody.TransactionDetails.GrossAmount = input.Amount
	reqBody.Qris.Acquirer = input.Acquirer

	payload, err := json.Marshal(reqBody)
	if err != nil {
		log.Errorf("[MidtransClient-1] failed to marshal request body: %v", err)
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))

	req, err := http.NewRequest(
		http.MethodPost,
		m.baseURL+"/v2/charge",
		bytes.NewBuffer(payload),
	)

	fmt.Println("[MidtransClient] Request URL:", m.baseURL+"/v2/charge")

	if err != nil {
		log.Errorf("[MidtransClient-2] failed to create new request: %v", err)
		return nil, err
	}

	fmt.Println("[MidtransClient] Request Created Successfully")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		log.Errorf("[MidtransClient-3] failed to perform request: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Errorf("[MidtransClient-4] received non-2xx response: %d", resp.StatusCode)
		return nil, fmt.Errorf("midtrans returned status code: %d", resp.StatusCode)
	}

	var result response.MidtransChargeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Errorf("[MidtransClient-5] failed to decode response body: %v", err)
		return nil, err
	}

	var qrUrl string
	for _, a := range result.Actions {
		if a.Name == "generate-qr-code" {
			qrUrl = a.URL
			break
		}
	}

	if qrUrl == "" {
		log.Errorf("[MidtransClient-6] qr code url not found in response")
		for _, a := range result.Actions {
			if a.Name == "generate-qr-code-v2" {
				qrUrl = a.URL
				break
			}
		}
	}

	if qrUrl == "" {
		log.Errorf("[MidtransClient-7] qr code url v2 not found in response")
		return nil, errors.New("qr code url not found in response")
	}

	fmt.Println("[MidtransClient-Result] QR URL:", qrUrl)

	return &domain.QrisChargeResult{
		OrderID: result.OrderID,
		QrURL:   qrUrl,
	}, nil
}

func NewMidtransClient(baseURL string) domain.MidtransClientInterface {
	return &midtransClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}
