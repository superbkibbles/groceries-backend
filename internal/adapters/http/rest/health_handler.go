package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
// @Summary Health check
// @Description Check if the API is healthy
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string "Healthy status"
// @Router /health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "groceries-api",
		"version": "1.0.0",
	})
}
