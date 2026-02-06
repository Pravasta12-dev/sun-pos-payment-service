package security

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
)

func ValidateMidtransSignature(
	orderID string,
	statusCode string,
	grossAmount string,
	serverKey string,
	signatureKey string,
) bool {
	payload := orderID + statusCode + grossAmount + serverKey

	hash := sha512.Sum512([]byte(payload))
	expectedSignature := hex.EncodeToString(hash[:])

	return strings.EqualFold(expectedSignature, signatureKey)
}
