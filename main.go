package main

import (
	"sun-pos-payment-service/cmd"
	. "sun-pos-payment-service/internal/security"
)

func main() {
	// JUST FOR DEMO PURPOSES
	TestEncryptDecrypt()

	cmd.Execute()
}
