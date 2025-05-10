package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// NotificationHandler handles HTTP requests for notifications
type NotificationHandler struct {
	notificationService ports.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler(notificationService ports.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// RegisterRoutes registers the notification routes
func (h *NotificationHandler) RegisterRoutes(router *gin.Engine) {
	notificationGroup := router.Group("/api/v1/notifications")
	{
		// User notifications
		notificationGroup.GET("", h.GetUserNotifications)
		notificationGroup.GET("/count", h.CountUnreadNotifications)
		notificationGroup.GET("/:id", h.GetNotification)
		notificationGroup.POST("", h.CreateNotification)
		notificationGroup.PUT("/:id/read", h.MarkAsRead)
		notificationGroup.PUT("/:id/archive", h.MarkAsArchived)
		notificationGroup.PUT("/read-all", h.MarkAllAsRead)
		notificationGroup.DELETE("/:id", h.DeleteNotification)

		// Notification templates (admin only)
		notificationGroup.GET("/templates", h.ListNotificationTemplates)
		notificationGroup.GET("/templates/:id", h.GetNotificationTemplate)
		notificationGroup.POST("/templates", h.CreateNotificationTemplate)
		notificationGroup.PUT("/templates/:id", h.UpdateNotificationTemplate)
		notificationGroup.DELETE("/templates/:id", h.DeleteNotificationTemplate)

		// Send notifications
		notificationGroup.POST("/send/order-status", h.SendOrderStatusNotification)
		notificationGroup.POST("/send/payment", h.SendPaymentNotification)
		notificationGroup.POST("/send/shipping", h.SendShippingNotification)
		notificationGroup.POST("/send/system", h.SendSystemNotification)
	}
}

// GetUserNotificationsRequest represents query parameters for listing notifications
type GetUserNotificationsRequest struct {
	Status string `form:"status" binding:"omitempty,oneof=unread read archived"`
	Page   int    `form:"page,default=1" binding:"min=1"`
	Limit  int    `form:"limit,default=10" binding:"min=1,max=100"`
}

// GetUserNotifications godoc
// @Summary Get user notifications
// @Description Get a list of notifications for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Param status query string false "Filter by status (unread, read, archived)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} PaginatedResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req GetUserNotificationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Convert status string to NotificationStatus
	status := entities.NotificationStatusUnread
	if req.Status != "" {
		status = entities.NotificationStatus(req.Status)
	}

	notifications, total, err := h.notificationService.ListUserNotifications(
		c.Request.Context(),
		userID.(string),
		status,
		req.Page,
		req.Limit,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       notifications,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: (total + req.Limit - 1) / req.Limit,
	})
}

// CountUnreadNotifications godoc
// @Summary Count unread notifications
// @Description Get the count of unread notifications for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/count [get]
func (h *NotificationHandler) CountUnreadNotifications(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	count, err := h.notificationService.CountUnreadNotifications(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetNotification godoc
// @Summary Get notification
// @Description Get a notification by ID
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} entities.Notification
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/{id} [get]
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	id := c.Param("id")

	notification, err := h.notificationService.GetNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Ensure the notification belongs to the user
	if notification.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "forbidden"})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// CreateNotificationRequest represents a request to create a notification
type CreateNotificationRequest struct {
	UserID  string                 `json:"user_id" binding:"required"`
	Type    string                 `json:"type" binding:"required"`
	Title   string                 `json:"title" binding:"required"`
	Message string                 `json:"message" binding:"required"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// CreateNotification godoc
// @Summary Create notification
// @Description Create a new notification (admin only)
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body CreateNotificationRequest true "Notification Details"
// @Success 201 {object} entities.Notification
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications [post]
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	// This endpoint should be admin-only
	// Check if user is admin (implementation depends on your auth system)

	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	notification, err := h.notificationService.CreateNotification(
		c.Request.Context(),
		req.UserID,
		entities.NotificationType(req.Type),
		req.Title,
		req.Message,
		req.Data,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// MarkAsRead godoc
// @Summary Mark notification as read
// @Description Mark a notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/{id}/read [put]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	id := c.Param("id")

	// Get notification to check ownership
	notification, err := h.notificationService.GetNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Ensure the notification belongs to the user
	if notification.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "forbidden"})
		return
	}

	err = h.notificationService.MarkNotificationAsRead(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// MarkAsArchived godoc
// @Summary Mark notification as archived
// @Description Mark a notification as archived
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/{id}/archive [put]
func (h *NotificationHandler) MarkAsArchived(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	id := c.Param("id")

	// Get notification to check ownership
	notification, err := h.notificationService.GetNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Ensure the notification belongs to the user
	if notification.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "forbidden"})
		return
	}

	err = h.notificationService.MarkNotificationAsArchived(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification archived"})
}

// MarkAllAsRead godoc
// @Summary Mark all notifications as read
// @Description Mark all notifications as read for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/read-all [put]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	err := h.notificationService.MarkAllNotificationsAsRead(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

// DeleteNotification godoc
// @Summary Delete notification
// @Description Delete a notification
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/{id} [delete]
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	id := c.Param("id")

	// Get notification to check ownership
	notification, err := h.notificationService.GetNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Ensure the notification belongs to the user
	if notification.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "forbidden"})
		return
	}

	err = h.notificationService.DeleteNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted"})
}

// CreateNotificationTemplateRequest represents a request to create a notification template
type CreateNotificationTemplateRequest struct {
	Name            string `json:"name" binding:"required"`
	Type            string `json:"type" binding:"required"`
	TitleTemplate   string `json:"title_template" binding:"required"`
	MessageTemplate string `json:"message_template" binding:"required"`
}

// ListNotificationTemplates godoc
// @Summary List notification templates
// @Description Get a list of notification templates (admin only)
// @Tags notifications
// @Accept json
// @Produce json
// @Param active query bool false "Filter by active status"
// @Success 200 {array} entities.NotificationTemplate
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/templates [get]
func (h *NotificationHandler) ListNotificationTemplates(c *gin.Context) {
	// This endpoint should be admin-only
	// Check if user is admin (implementation depends on your auth system)

	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active", "true"))

	templates, err := h.notificationService.ListNotificationTemplates(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// GetNotificationTemplate godoc
// @Summary Get notification template
// @Description Get a notification template by ID (admin only)
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} entities.NotificationTemplate
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/templates/{id} [get]
func (h *NotificationHandler) GetNotificationTemplate(c *gin.Context) {
	// This endpoint should be admin-only
	// Check if user is admin (implementation depends on your auth system)

	id := c.Param("id")

	template, err := h.notificationService.GetNotificationTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// CreateNotificationTemplate godoc
// @Summary Create notification template
// @Description Create a new notification template (admin only)
// @Tags notifications
// @Accept json
// @Produce json
// @Param template body CreateNotificationTemplateRequest true "Template Details"
// @Success 201 {object} entities.NotificationTemplate
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/templates [post]
func (h *NotificationHandler) CreateNotificationTemplate(c *gin.Context) {
	// This endpoint should be admin-only
	// Check if user is admin (implementation depends on your auth system)

	var req CreateNotificationTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	template, err := h.notificationService.CreateNotificationTemplate(
		c.Request.Context(),
		req.Name,
		entities.NotificationType(req.Type),
		req.TitleTemplate,
		req.MessageTemplate,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// UpdateNotificationTemplate godoc
// @Summary Update notification template
// @Description Update an existing notification template (admin only)
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Param template body CreateNotificationTemplateRequest true "Template Details"
// @Success 200 {object} entities.NotificationTemplate
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/templates/{id} [put]
func (h *NotificationHandler) UpdateNotificationTemplate(c *gin.Context) {
	// This endpoint should be admin-only
	// Check if user is admin (implementation depends on your auth system)

	id := c.Param("id")

	var req CreateNotificationTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get existing template
	template, err := h.notificationService.GetNotificationTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Update fields
	template.Name = req.Name
	template.Type = entities.NotificationType(req.Type)
	template.TitleTemplate = req.TitleTemplate
	template.MessageTemplate = req.MessageTemplate

	// Save changes
	err = h.notificationService.UpdateNotificationTemplate(c.Request.Context(), template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteNotificationTemplate godoc
// @Summary Delete notification template
// @Description Delete a notification template (admin only)
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/templates/{id} [delete]
func (h *NotificationHandler) DeleteNotificationTemplate(c *gin.Context) {
	// This endpoint should be admin-only
	// Check if user is admin (implementation depends on your auth system)

	id := c.Param("id")

	err := h.notificationService.DeleteNotificationTemplate(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "notification template not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification template deleted successfully"})
}

// SendOrderStatusNotificationRequest represents a request to send an order status notification
type SendOrderStatusNotificationRequest struct {
	OrderID string               `json:"order_id" binding:"required"`
	Status  entities.OrderStatus `json:"status" binding:"required"`
}

// SendOrderStatusNotification godoc
// @Summary Send order status notification
// @Description Send a notification about an order status change
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body SendOrderStatusNotificationRequest true "Notification Details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/send/order-status [post]
func (h *NotificationHandler) SendOrderStatusNotification(c *gin.Context) {
	var req SendOrderStatusNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.notificationService.SendOrderStatusNotification(
		c.Request.Context(),
		req.OrderID,
		req.Status,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status notification sent"})
}

// SendPaymentNotificationRequest represents a request to send a payment notification
type SendPaymentNotificationRequest struct {
	OrderID     string               `json:"order_id" binding:"required"`
	PaymentInfo entities.PaymentInfo `json:"payment_info" binding:"required"`
}

// SendPaymentNotification godoc
// @Summary Send payment notification
// @Description Send a notification about a payment
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body SendPaymentNotificationRequest true "Notification Details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/send/payment [post]
func (h *NotificationHandler) SendPaymentNotification(c *gin.Context) {
	var req SendPaymentNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.notificationService.SendPaymentNotification(
		c.Request.Context(),
		req.OrderID,
		req.PaymentInfo,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment notification sent"})
}

// SendShippingNotificationRequest represents a request to send a shipping notification
type SendShippingNotificationRequest struct {
	OrderID      string                `json:"order_id" binding:"required"`
	TrackingInfo entities.ShippingInfo `json:"tracking_info" binding:"required"`
}

// SendShippingNotification godoc
// @Summary Send shipping notification
// @Description Send a notification about shipping
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body SendShippingNotificationRequest true "Notification Details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/send/shipping [post]
func (h *NotificationHandler) SendShippingNotification(c *gin.Context) {
	var req SendShippingNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.notificationService.SendShippingNotification(
		c.Request.Context(),
		req.OrderID,
		req.TrackingInfo,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping notification sent"})
}

// SendSystemNotificationRequest represents a request to send a system notification
type SendSystemNotificationRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// SendSystemNotification godoc
// @Summary Send system notification
// @Description Send a system notification to a user
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body SendSystemNotificationRequest true "Notification Details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/notifications/send/system [post]
func (h *NotificationHandler) SendSystemNotification(c *gin.Context) {
	var req SendSystemNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.notificationService.SendSystemNotification(
		c.Request.Context(),
		req.UserID,
		req.Title,
		req.Message,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "System notification sent"})
}
