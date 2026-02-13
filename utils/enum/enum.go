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
