package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/superbkibbles/ecommerce/docs"
	"github.com/superbkibbles/ecommerce/internal/adapters/http/rest"
	"github.com/superbkibbles/ecommerce/internal/adapters/repository/mongodb"
	"github.com/superbkibbles/ecommerce/internal/adapters/repository/redisdb"
	"github.com/superbkibbles/ecommerce/internal/application/services"
	"github.com/superbkibbles/ecommerce/internal/config"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title E-Commerce API
// @version 1.0
// @description A Hexagonal Architecture E-Commerce API with product variations support
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost
// @BasePath /api/v1
// @schemes http
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup database connection
	db, err := mongodb.NewMongoDBConnection(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient, err := redisdb.NewRedisConnection(cfg.Redis.URI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	productRepo := mongodb.NewProductRepository(db)
	orderRepo := mongodb.NewOrderRepository(db)
	categoryRepo := mongodb.NewCategoryRepository(db, productRepo)
	userRepo := mongodb.NewUserRepository(db, redisClient)
	settingRepo := mongodb.NewSettingRepository(db)
	homeSectionRepo := mongodb.NewHomeSectionRepository(db)

	// Initialize services
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(orderRepo, productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	userService := services.NewUserService(userRepo)
	settingService := services.NewSettingService(settingRepo)
	homeSectionService := services.NewHomeSectionService(homeSectionRepo)

	// Ensure superuser admin exists if configured via environment variables
	if cfg.Server.AdminEmail != "" && cfg.Server.AdminPassword != "" {
		ctx := context.Background()
		if _, err := userRepo.GetByEmail(ctx, cfg.Server.AdminEmail); err != nil {
			adminUser, errNew := entities.NewUser(
				cfg.Server.AdminEmail,
				cfg.Server.AdminPassword,
				cfg.Server.AdminFirstName,
				cfg.Server.AdminLastName,
				entities.UserRoleAdmin,
			)
			if errNew != nil {
				log.Printf("Failed to construct admin user: %v", errNew)
			} else {
				if errCreate := userRepo.Create(ctx, adminUser); errCreate != nil {
					log.Printf("Failed to create admin user: %v", errCreate)
				} else {
					log.Printf("Admin user ensured: %s", cfg.Server.AdminEmail)
				}
			}
		} else {
			// Admin exists; ensure it's active and has admin role
			existing, _ := userRepo.GetByEmail(ctx, cfg.Server.AdminEmail)
			if existing != nil {
				existing.Active = true
				existing.Role = entities.UserRoleAdmin
				existing.UpdatedAt = time.Now()
				if errUpdate := userRepo.Update(ctx, existing); errUpdate != nil {
					log.Printf("Failed to update existing admin user: %v", errUpdate)
				}
			}
		}
	}

	// Setup Gin router
	router := gin.Default()

	// Setup middleware including CORS
	rest.SetupMiddleware(router)

	// Setup reverse proxy for admin panel
	adminURL, err := url.Parse(cfg.Server.AdminBaseURL)
	if err != nil {
		log.Fatalf("Failed to parse admin base URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(adminURL)

	// Custom director to handle path rewriting
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = adminURL.Host
		// Remove /admin prefix as Next.js already has basePath configured
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/admin")
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}
	}

	// Admin panel routes - must be registered before API routes
	router.Any("/admin", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/admin/")
	})
	router.Any("/admin/*proxyPath", func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// Setup API routes
	api := router.Group("/api/v1")
	rest.NewProductHandler(api, productService)
	rest.NewOrderHandler(api, orderService)
	rest.NewCategoryHandler(api, categoryService)
	rest.NewUserHandler(api, userService, orderService)
	rest.NewLanguageHandler(api)

	// Setup home section handler (public + admin endpoints)
	rest.NewHomeSectionHandler(router, homeSectionService)

	// Setup settings handler
	settingHandler := rest.NewSettingHandler(settingService)
	settingHandler.RegisterRoutes(router)

	// Health check endpoint
	router.GET("/api/v1/health", rest.HealthHandler)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
