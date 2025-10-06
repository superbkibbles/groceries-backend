package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/superbkibbles/ecommerce/internal/adapters/http/middleware"
	"github.com/superbkibbles/ecommerce/internal/config"
)

// T translates a message using the localizer from the Gin context
func T(c *gin.Context, messageID string, data map[string]interface{}) string {
	localizer := middleware.GetLocalizerFromContext(c)
	i18nManager := config.NewI18nManager()
	return i18nManager.T(localizer, messageID, data)
}

// TSimple translates a message without template data
func TSimple(c *gin.Context, messageID string) string {
	return T(c, messageID, nil)
}

// TWithData translates a message with template data
func TWithData(c *gin.Context, messageID string, data map[string]interface{}) string {
	return T(c, messageID, data)
}

// GetSupportedLanguages returns the list of supported languages
func GetSupportedLanguages() []string {
	i18nManager := config.NewI18nManager()
	return i18nManager.GetSupportedLanguages()
}

// GetDefaultLanguage returns the default language
func GetDefaultLanguage() string {
	i18nManager := config.NewI18nManager()
	return i18nManager.GetDefaultLanguage()
}

// SetLanguage sets the language cookie for the user
func SetLanguage(c *gin.Context, lang string) {
	i18nManager := config.NewI18nManager()
	middleware.SetLanguageCookie(c, lang, i18nManager)
}

// GetLocalizerFromContext is a convenience function to get the localizer
func GetLocalizerFromContext(c *gin.Context) *i18n.Localizer {
	return middleware.GetLocalizerFromContext(c)
}

// GetLanguageFromRequest extracts the language from the request
func GetLanguageFromRequest(c *gin.Context) string {
	// Try query parameter first
	if lang := c.Query("lang"); lang != "" {
		return lang
	}

	// Try Accept-Language header
	if lang := c.GetHeader("Accept-Language"); lang != "" {
		// Simple parsing - if contains "ar" then Arabic, otherwise English
		if len(lang) >= 2 && lang[:2] == "ar" {
			return "ar"
		}
	}

	// Default to English
	return "en"
}
