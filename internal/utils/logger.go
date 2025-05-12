package utils

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger is a global logger instance
var Logger zerolog.Logger

func init() {
	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set up pretty logging for development
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}

	// Create logger
	Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	// Set log level based on environment
	if os.Getenv("ENV") == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
