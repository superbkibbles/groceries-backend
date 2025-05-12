package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}

// SendSMSRequest sends an SMS request to the SMS service
func SendSMSRequest(ctx context.Context, apiURL, token, from, phoneNumber, message string) error {
	// Prepare form data
	formData := url.Values{}
	formData.Add("from", from)
	formData.Add("phone_number", phoneNumber)
	formData.Add("message", message)

	// Create request with form data
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiURL,
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	Logger.Debug().
		Str("url", req.URL.String()).
		Str("method", req.Method).
		Str("from", from).
		Str("phone", phoneNumber).
		Msg("Sending SMS request")

	// Send request
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	Logger.Debug().
		Int("status_code", resp.StatusCode).
		Msg("SMS request successful")

	return nil
}
