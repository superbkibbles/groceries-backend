package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SettingHandler handles HTTP requests for settings
type SettingHandler struct {
	settingService ports.SettingService
}

// NewSettingHandler creates a new setting handler
func NewSettingHandler(settingService ports.SettingService) *SettingHandler {
	return &SettingHandler{
		settingService: settingService,
	}
}

// RegisterRoutes registers the routes for settings
func (h *SettingHandler) RegisterRoutes(router *gin.Engine) {
	settings := router.Group("/api/settings")
	{
		// System settings (admin only)
		settings.POST("/system", AdminRequired(), h.CreateSystemSetting)
		settings.GET("/system/:key", h.GetSystemSetting)
		settings.GET("/system", h.ListSystemSettings)
		settings.PUT("/system/:id", AdminRequired(), h.UpdateSetting)
		settings.DELETE("/system/:id", AdminRequired(), h.DeleteSetting)

		// User settings (authenticated user)
		settings.POST("/user", AuthRequired(), h.CreateUserSetting)
		settings.GET("/user/:key", AuthRequired(), h.GetUserSetting)
		settings.GET("/user", AuthRequired(), h.ListUserSettings)
		settings.PUT("/user/:id", AuthRequired(), h.UpdateUserSetting)
		settings.DELETE("/user/:id", AuthRequired(), h.DeleteUserSetting)
	}
}

// CreateSystemSetting handles the creation of a system setting
func (h *SettingHandler) CreateSystemSetting(c *gin.Context) {
	var request struct {
		Key         string      `json:"key" binding:"required"`
		Value       interface{} `json:"value" binding:"required"`
		Type        string      `json:"type" binding:"required,oneof=string number boolean json"`
		Description string      `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setting, err := h.settingService.CreateSystemSetting(
		c.Request.Context(),
		request.Key,
		request.Value,
		entities.SettingType(request.Type),
		request.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, setting)
}

// GetSystemSetting handles retrieving a system setting by key
func (h *SettingHandler) GetSystemSetting(c *gin.Context) {
	key := c.Param("key")

	setting, err := h.settingService.GetSystemSetting(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// ListSystemSettings handles listing all system settings
func (h *SettingHandler) ListSystemSettings(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Create filter from query parameters
	filter := make(map[string]interface{})
	if key := c.Query("key"); key != "" {
		filter["key"] = key
	}

	settings, total, err := h.settingService.ListSystemSettings(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  settings,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// UpdateSetting handles updating a setting
func (h *SettingHandler) UpdateSetting(c *gin.Context) {
	id := c.Param("id")

	var request struct {
		Value       interface{} `json:"value" binding:"required"`
		Description string      `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing setting
	setting, err := h.settingService.GetSystemSetting(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	// Update fields
	setting.Value = request.Value
	setting.Description = request.Description

	// Save changes
	if err := h.settingService.UpdateSetting(c.Request.Context(), setting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// DeleteSetting handles deleting a setting
func (h *SettingHandler) DeleteSetting(c *gin.Context) {
	id := c.Param("id")

	if err := h.settingService.DeleteSetting(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateUserSetting handles the creation of a user setting
func (h *SettingHandler) CreateUserSetting(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request struct {
		Key         string      `json:"key" binding:"required"`
		Value       interface{} `json:"value" binding:"required"`
		Type        string      `json:"type" binding:"required,oneof=string number boolean json"`
		Description string      `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setting, err := h.settingService.CreateUserSetting(
		c.Request.Context(),
		request.Key,
		request.Value,
		entities.SettingType(request.Type),
		request.Description,
		userID.(string),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, setting)
}

// GetUserSetting handles retrieving a user setting by key
func (h *SettingHandler) GetUserSetting(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	key := c.Param("key")

	setting, err := h.settingService.GetUserSetting(c.Request.Context(), key, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// ListUserSettings handles listing all settings for a user
func (h *SettingHandler) ListUserSettings(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Create filter from query parameters
	filter := make(map[string]interface{})
	if key := c.Query("key"); key != "" {
		filter["key"] = key
	}

	settings, total, err := h.settingService.ListUserSettings(
		c.Request.Context(),
		userID.(string),
		filter,
		page,
		limit,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  settings,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// UpdateUserSetting handles updating a user setting
func (h *SettingHandler) UpdateUserSetting(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")

	var request struct {
		Value       interface{} `json:"value" binding:"required"`
		Description string      `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing setting
	setting, err := h.settingService.GetUserSetting(c.Request.Context(), id, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	// Verify ownership
	userObjectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if setting.UserID != userObjectID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this setting"})
		return
	}

	// Update fields
	setting.Value = request.Value
	setting.Description = request.Description

	// Save changes
	if err := h.settingService.UpdateSetting(c.Request.Context(), setting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// DeleteUserSetting handles deleting a user setting
func (h *SettingHandler) DeleteUserSetting(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")

	// Get existing setting to verify ownership
	setting, err := h.settingService.GetUserSetting(c.Request.Context(), id, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	// Verify ownership
	userObjectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if setting.UserID != userObjectID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this setting"})
		return
	}

	if err := h.settingService.DeleteSetting(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
