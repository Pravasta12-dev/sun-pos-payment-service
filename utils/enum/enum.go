package enum

type TransactionScope string

const (
	ScopeMerchant TransactionScope = "merchant"
	ScopeOwner    TransactionScope = "owner"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusPaid    TransactionStatus = "paid"
	TransactionStatusExpired TransactionStatus = "expired"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type PaymentChannelType string

const (
	PaymentChannelTypeVa      PaymentChannelType = "va"
	PaymentChannelTypeEwallet PaymentChannelType = "ewallet"
	PaymentChannelTypeQris    PaymentChannelType = "qris"
)

type PaymentChannelCode string

const (
	PaymentChannelCodeBca  PaymentChannelCode = "bca"
	PaymentChannelCodeBri  PaymentChannelCode = "bri"
	PaymentChannelCodeQris PaymentChannelCode = "qris"
	PaymentChannelCodeBni  PaymentChannelCode = "bni"
)

type PaymentFeeType string

const (
	PaymentFeeTypeFlat    PaymentFeeType = "flat"
	PaymentFeeTypePercent PaymentFeeType = "percentage"
)
