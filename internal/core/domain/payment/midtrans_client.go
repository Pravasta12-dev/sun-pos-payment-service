package payment

type MidtransClientInterface interface {
	ChargeQris(serverKey string, input QrisChargeInput) (*QrisChargeResult, error)
	ChargeVa(serverKey string, input VaChargeInput) (*VaChargeResult, error)
}
