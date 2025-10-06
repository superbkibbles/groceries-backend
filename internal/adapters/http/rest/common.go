package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/utils"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	Message    string      `json:"message,omitempty"`
}

// Helper functions for common responses

// Success sends a success response with localized message
func Success(c *gin.Context, statusCode int, messageKey string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Message: utils.TSimple(c, messageKey),
		Data:    data,
	})
}

// SuccessWithData sends a success response with localized message and template data
func SuccessWithData(c *gin.Context, statusCode int, messageKey string, templateData map[string]interface{}, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Message: utils.TWithData(c, messageKey, templateData),
		Data:    data,
	})
}

// Error sends an error response with localized message
func Error(c *gin.Context, statusCode int, messageKey string) {
	c.JSON(statusCode, ErrorResponse{
		Error:   utils.TSimple(c, messageKey),
		Message: utils.TSimple(c, messageKey),
		Code:    statusCode,
	})
}

// ErrorWithData sends an error response with localized message and template data
func ErrorWithData(c *gin.Context, statusCode int, messageKey string, templateData map[string]interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Error:   utils.TWithData(c, messageKey, templateData),
		Message: utils.TWithData(c, messageKey, templateData),
		Code:    statusCode,
	})
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, messageKey string) {
	Error(c, http.StatusBadRequest, messageKey)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, messageKey string) {
	Error(c, http.StatusUnauthorized, messageKey)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, messageKey string) {
	Error(c, http.StatusForbidden, messageKey)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, messageKey string) {
	Error(c, http.StatusNotFound, messageKey)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, messageKey string) {
	Error(c, http.StatusInternalServerError, messageKey)
}

// Created sends a 201 Created response
func Created(c *gin.Context, messageKey string, data interface{}) {
	Success(c, http.StatusCreated, messageKey, data)
}

// OK sends a 200 OK response
func OK(c *gin.Context, messageKey string, data interface{}) {
	Success(c, http.StatusOK, messageKey, data)
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, data interface{}, total, page, limit, totalPages int, messageKey string) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Message:    utils.TSimple(c, messageKey),
	})
}
