package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/adapters/http/middleware"
	"github.com/superbkibbles/ecommerce/internal/config"
)

// SetupMiddleware configures global middleware for the router
func SetupMiddleware(router *gin.Engine) {
	// Initialize i18n manager
	i18nManager := config.NewI18nManager()

	// CORS middleware
	router.Use(corsMiddleware())

	// Internationalization middleware
	router.Use(middleware.I18nMiddleware(i18nManager))

	// Request logging
	router.Use(gin.Logger())

	// Recovery middleware
	router.Use(gin.Recovery())
}

// corsMiddleware configures CORS for the API
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, user_role, Accept-Language")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Language")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// AuthRequired returns the authentication middleware
func AuthRequired() gin.HandlerFunc {
	return middleware.AuthMiddleware()
}

// AdminRequired returns the admin authorization middleware
func AdminRequired() gin.HandlerFunc {
	return middleware.AdminMiddleware()
}
