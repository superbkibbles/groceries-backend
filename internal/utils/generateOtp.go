package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTP(length int) (string, error) {
	otp := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10)) // random digit 0–9
		if err != nil {
			return "", err
		}
		otp += n.String()
	}
	return otp, nil
}
