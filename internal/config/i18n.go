package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// I18nConfig holds internationalization configuration
type I18nConfig struct {
	DefaultLanguage    string
	SupportedLanguages []string
	BundlePath         string
}

// I18nManager handles internationalization
type I18nManager struct {
	bundle *i18n.Bundle
	config *I18nConfig
}

// NewI18nManager creates a new internationalization manager
func NewI18nManager() *I18nManager {
	config := &I18nConfig{
		DefaultLanguage:    "en",
		SupportedLanguages: []string{"en", "ar"},
		BundlePath:         "internal/locales",
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load translation files
	localesPath := config.BundlePath
	if _, err := os.Stat(localesPath); os.IsNotExist(err) {
		log.Printf("Locales directory not found: %s", localesPath)
		return &I18nManager{bundle: bundle, config: config}
	}

	// Load English translations
	enFile := filepath.Join(localesPath, "en.json")
	if _, err := os.Stat(enFile); err == nil {
		bundle.MustLoadMessageFile(enFile)
	}

	// Load Arabic translations
	arFile := filepath.Join(localesPath, "ar.json")
	if _, err := os.Stat(arFile); err == nil {
		bundle.MustLoadMessageFile(arFile)
	}

	return &I18nManager{
		bundle: bundle,
		config: config,
	}
}

// GetLocalizer returns a localizer for the given language
func (im *I18nManager) GetLocalizer(lang string) *i18n.Localizer {
	// Validate language
	supported := false
	for _, supportedLang := range im.config.SupportedLanguages {
		if lang == supportedLang {
			supported = true
			break
		}
	}

	if !supported {
		lang = im.config.DefaultLanguage
	}

	return i18n.NewLocalizer(im.bundle, lang)
}

// GetSupportedLanguages returns the list of supported languages
func (im *I18nManager) GetSupportedLanguages() []string {
	return im.config.SupportedLanguages
}

// GetDefaultLanguage returns the default language
func (im *I18nManager) GetDefaultLanguage() string {
	return im.config.DefaultLanguage
}

// T translates a message using the given localizer
func (im *I18nManager) T(localizer *i18n.Localizer, messageID string, data map[string]interface{}) string {
	config := &i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	}

	message, err := localizer.Localize(config)
	if err != nil {
		log.Printf("Translation error for message ID '%s': %v", messageID, err)
		return messageID // Fallback to message ID
	}

	return message
}
