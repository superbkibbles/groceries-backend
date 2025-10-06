package sms

import (
	"context"
	"fmt"

	"github.com/superbkibbles/ecommerce/internal/config"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
	"github.com/superbkibbles/ecommerce/internal/utils"
)

// SMSSenderAdapter implements the SMSSender port
type SMSSenderAdapter struct {
	config *config.SMSConfig
}

// NewSMSSender creates a new SMS sender adapter
func NewSMSSender(cfg *config.SMSConfig) ports.SMSSender {
	return &SMSSenderAdapter{
		config: cfg,
	}
}

// SendOTP sends an OTP via SMS
func (s *SMSSenderAdapter) SendOTP(phoneNumber string, otp string) error {
	// Create OTP message
	otpMessage := fmt.Sprintf("Your National Taxi Password is: %s", otp)

	// Send the SMS with OTP
	ctx := context.Background()
	if err := utils.SendSMSRequest(
		ctx,
		s.config.APIURL,
		s.config.Token,
		s.config.SenderID,
		phoneNumber,
		otpMessage,
		s.config.ProductType,
	); err != nil {
		utils.Logger.Error().
			Err(err).
			Str("phone", phoneNumber).
			Msg("Failed to send OTP")
		return fmt.Errorf("failed to send OTP: %w", err)
	}

	utils.Logger.Info().
		Str("phone", phoneNumber).
		Str("link", s.config.APIURL).
		Msg("OTP sent successfully")

	return nil
}
