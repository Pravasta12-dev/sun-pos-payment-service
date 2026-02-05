package payment

type MidtransClientInterface interface {
	ChargeQris(serverKey string, input QrisChargeInput) (*QrisChargeResult, error)
}
