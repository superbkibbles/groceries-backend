package config

import (
	"os"
)

// SMSConfig contains configuration for SMS service
type SMSConfig struct {
	APIURL      string
	Token       string
	SenderID    string
	ProductType string
}

// NewSMSConfig creates a new SMS configuration from environment variables
func NewSMSConfig() *SMSConfig {
	return &SMSConfig{
		APIURL:      getEnvOrDefault("SMS_API_URL", "https://smscloud.alsardfiber.com/api/v1/campaign/infinite"),
		Token:       getEnvOrDefault("SMS_API_TOKEN", "warqVMC0OwKcVJby7HYcAvplV9wEMbnQBX8bArBQ6cZDnPp0WsMDjE1Wdklr"),
		SenderID:    getEnvOrDefault("SMS_SENDER_ID", "Watany-taxi"),
		ProductType: getEnvOrDefault("SMS_PRODUCT_TYPE", "nationalcmp2"),
	}
}

// Helper function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
