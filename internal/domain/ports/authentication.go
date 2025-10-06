package ports

import (
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
)

// TokenGenerator generates authentication tokens
type TokenGenerator interface {
	GenerateToken(user *entities.User) (string, error)
	ValidateToken(token string) (map[string]interface{}, error)
}

// SMSSender sends SMS messages
type SMSSender interface {
	SendOTP(phoneNumber string, otp string) error
}

// OTPGenerator generates one-time passwords
type OTPGenerator interface {
	Generate(length int) (string, error)
}
