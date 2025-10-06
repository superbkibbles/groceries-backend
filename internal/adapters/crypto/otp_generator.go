package crypto

import (
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
	"github.com/superbkibbles/ecommerce/internal/utils"
)

// OTPGeneratorAdapter implements the OTPGenerator port
type OTPGeneratorAdapter struct{}

// NewOTPGenerator creates a new OTP generator adapter
func NewOTPGenerator() ports.OTPGenerator {
	return &OTPGeneratorAdapter{}
}

// Generate generates a random OTP of the specified length
func (o *OTPGeneratorAdapter) Generate(length int) (string, error) {
	return utils.GenerateOTP(length)
}
