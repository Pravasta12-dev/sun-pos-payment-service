package message

import "errors"

var (
	ErrMerchantNotFound = errors.New("Merchant Not Found")
	ErrMerchantConflict = errors.New("Merchant Already Exists")
	ErrDatabaseFailure  = errors.New("Database Failure")
)
