package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	Server  ServerConfig
	MongoDB MongoDBConfig
	I18n    I18nConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string
	AdminBaseURL string
}

// MongoDBConfig holds MongoDB connection configuration
type MongoDBConfig struct {
	URI      string
	Database string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	port := getEnv("SERVER_PORT", "8080")
	adminBaseURL := getEnv("ADMIN_BASE_URL", "http://localhost:3000")
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getEnv("MONGO_DB", "durra")
	defaultLang := getEnv("DEFAULT_LANGUAGE", "en")
	supportedLangs := getEnv("SUPPORTED_LANGUAGES", "en,ar")
	bundlePath := getEnv("I18N_BUNDLE_PATH", "internal/locales")

	return &Config{
		Server: ServerConfig{
			Port:         port,
			AdminBaseURL: adminBaseURL,
		},
		MongoDB: MongoDBConfig{
			URI:      mongoURI,
			Database: mongoDB,
		},
		I18n: I18nConfig{
			DefaultLanguage:    defaultLang,
			SupportedLanguages: strings.Split(supportedLangs, ","),
			BundlePath:         bundlePath,
		},
	}, nil
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get integer environment variables with defaults
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
