package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/utils"
)

// LanguageHandler handles language-related requests
type LanguageHandler struct{}

// NewLanguageHandler creates a new language handler and registers routes
func NewLanguageHandler(router *gin.RouterGroup) {
	handler := &LanguageHandler{}

	languages := router.Group("/languages")
	{
		languages.GET("", handler.GetSupportedLanguages)
		languages.POST("/switch", handler.SwitchLanguage)
	}
}

// GetSupportedLanguagesResponse represents the response for supported languages
type GetSupportedLanguagesResponse struct {
	Languages       []string `json:"languages"`
	DefaultLanguage string   `json:"default_language"`
}

// SwitchLanguageRequest represents the request for switching language
type SwitchLanguageRequest struct {
	Language string `json:"language" binding:"required"`
}

// GetSupportedLanguages returns the list of supported languages
// @Summary Get supported languages
// @Description Get the list of supported languages and default language
// @Tags Language
// @Accept json
// @Produce json
// @Success 200 {object} GetSupportedLanguagesResponse
// @Router /api/v1/languages [get]
func (h *LanguageHandler) GetSupportedLanguages(c *gin.Context) {
	supportedLanguages := utils.GetSupportedLanguages()
	defaultLanguage := utils.GetDefaultLanguage()

	OK(c, "success", GetSupportedLanguagesResponse{
		Languages:       supportedLanguages,
		DefaultLanguage: defaultLanguage,
	})
}

// SwitchLanguage switches the user's language preference
// @Summary Switch language
// @Description Switch the user's language preference by setting a cookie
// @Tags Language
// @Accept json
// @Produce json
// @Param request body SwitchLanguageRequest true "Language switch request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/languages/switch [post]
func (h *LanguageHandler) SwitchLanguage(c *gin.Context) {
	var req SwitchLanguageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "validation_error")
		return
	}

	// Check if the requested language is supported
	supportedLanguages := utils.GetSupportedLanguages()
	isSupported := false
	for _, lang := range supportedLanguages {
		if lang == req.Language {
			isSupported = true
			break
		}
	}

	if !isSupported {
		BadRequest(c, "validation_error")
		return
	}

	// Set the language cookie
	utils.SetLanguage(c, req.Language)

	// Return success response with the new language
	templateData := map[string]interface{}{
		"language": req.Language,
	}
	SuccessWithData(c, http.StatusOK, "success", templateData, map[string]string{
		"language": req.Language,
	})
}
