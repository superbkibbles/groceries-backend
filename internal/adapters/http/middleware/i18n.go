package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/superbkibbles/ecommerce/internal/config"
)

const (
	LanguageHeader     = "Accept-Language"
	LanguageQuery      = "lang"
	LanguageCookie     = "language"
	LanguageContextKey = "localizer"
)

// I18nMiddleware creates middleware for internationalization
func I18nMiddleware(i18nManager *config.I18nManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get language from multiple sources with priority
		lang := getLanguageFromRequest(c, i18nManager)

		// Get localizer for the detected language
		localizer := i18nManager.GetLocalizer(lang)

		// Store localizer in context
		c.Set(LanguageContextKey, localizer)

		// Continue to next handler
		c.Next()
	}
}

// getLanguageFromRequest extracts language from request with priority order:
// 1. Query parameter (?lang=ar)
// 2. Header (Accept-Language: ar)
// 3. Cookie (language=ar)
// 4. Default language
func getLanguageFromRequest(c *gin.Context, i18nManager *config.I18nManager) string {
	// 1. Check query parameter
	if lang := c.Query(LanguageQuery); lang != "" {
		if isValidLanguage(lang, i18nManager) {
			return lang
		}
	}

	// 2. Check Accept-Language header
	if lang := c.GetHeader(LanguageHeader); lang != "" {
		// Parse Accept-Language header (e.g., "ar,en;q=0.9,en-US;q=0.8")
		languages := parseAcceptLanguage(lang)
		for _, l := range languages {
			if isValidLanguage(l, i18nManager) {
				return l
			}
		}
	}

	// 3. Check cookie
	if lang, err := c.Cookie(LanguageCookie); err == nil && lang != "" {
		if isValidLanguage(lang, i18nManager) {
			return lang
		}
	}

	// 4. Return default language
	return i18nManager.GetDefaultLanguage()
}

// parseAcceptLanguage parses Accept-Language header and returns languages in order of preference
func parseAcceptLanguage(header string) []string {
	var languages []string

	// Split by comma and process each language
	parts := strings.Split(header, ",")
	for _, part := range parts {
		// Remove quality values (q=0.9) and whitespace
		lang := strings.TrimSpace(strings.Split(part, ";")[0])

		// Handle language tags like "en-US" -> "en"
		if strings.Contains(lang, "-") {
			lang = strings.Split(lang, "-")[0]
		}

		if lang != "" {
			languages = append(languages, lang)
		}
	}

	return languages
}

// isValidLanguage checks if the language is supported
func isValidLanguage(lang string, i18nManager *config.I18nManager) bool {
	supportedLanguages := i18nManager.GetSupportedLanguages()
	for _, supported := range supportedLanguages {
		if lang == supported {
			return true
		}
	}
	return false
}

// GetLocalizerFromContext retrieves the localizer from the Gin context
func GetLocalizerFromContext(c *gin.Context) *i18n.Localizer {
	if localizer, exists := c.Get(LanguageContextKey); exists {
		if l, ok := localizer.(*i18n.Localizer); ok {
			return l
		}
	}

	// Fallback: create a default localizer
	// This should not happen if middleware is properly set up
	i18nManager := config.NewI18nManager()
	return i18nManager.GetLocalizer(i18nManager.GetDefaultLanguage())
}

// SetLanguageCookie sets a cookie with the selected language
func SetLanguageCookie(c *gin.Context, lang string, i18nManager *config.I18nManager) {
	if isValidLanguage(lang, i18nManager) {
		c.SetCookie(LanguageCookie, lang, 86400*30, "/", "", false, true) // 30 days
	}
}
