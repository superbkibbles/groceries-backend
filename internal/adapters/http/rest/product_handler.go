package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
	"github.com/superbkibbles/ecommerce/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	productService ports.ProductService
}

// NewProductHandler creates a new product handler and registers routes
func NewProductHandler(router *gin.RouterGroup, productService ports.ProductService) {
	handler := &ProductHandler{
		productService: productService,
	}

	products := router.Group("/products")
	{
		products.POST("", handler.CreateProduct)
		products.GET("", handler.ListProducts)
		products.GET("/:id", handler.GetProduct)
		products.PUT("/:id", handler.UpdateProduct)
		products.DELETE("/:id", handler.DeleteProduct)
		products.GET("/category/:category", handler.GetProductsByCategory)
		products.PUT("/:id/stock", handler.UpdateStock)
	}
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Categories    []string                        `json:"categories"`
	Attributes    map[string]any                  `json:"attributes"`
	SKU           string                          `json:"sku" binding:"required"`
	Price         float64                         `json:"price" binding:"required,gt=0"`
	StockQuantity int                             `json:"stock_quantity" binding:"required,gte=0"`
	Images        []string                        `json:"images"`
	Translations  map[string]entities.Translation `json:"translations" binding:"required"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name          string         `json:"name" binding:"required"`
	Description   string         `json:"description"`
	Categories    []string       `json:"categories"`
	Attributes    map[string]any `json:"attributes"`
	SKU           string         `json:"sku" binding:"required"`
	Price         float64        `json:"price" binding:"required,gt=0"`
	StockQuantity int            `json:"stock_quantity" binding:"required,gte=0"`
	Images        []string       `json:"images"`
}

// StockUpdateRequest represents the request body for updating stock
type StockUpdateRequest struct {
	Quantity int `json:"quantity" binding:"required,gte=0"`
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product details"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "validation_error")
		return
	}

	// Validate that English translation is provided
	if _, hasEnglish := req.Translations["en"]; !hasEnglish {
		BadRequest(c, "product_name_required")
		return
	}

	product, err := h.productService.CreateProduct(
		c.Request.Context(),
		req.Categories,
		req.Attributes,
		req.SKU,
		req.Price,
		req.StockQuantity,
		req.Images,
		req.Translations,
	)
	if err != nil {
		InternalServerError(c, "internal_server_error")
		return
	}

	Created(c, "product_created", product)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get detailed information about a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	language := utils.GetLanguageFromRequest(c)

	product, err := h.productService.GetProduct(c.Request.Context(), id, language)
	if err != nil {
		NotFound(c, "product_not_found")
		return
	}

	OK(c, "product_retrieved", product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product's details
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body UpdateProductRequest true "Updated product details"
// @Success 200 {object} entities.Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	// Get existing product (need to get it first to update)
	language := utils.GetLanguageFromRequest(c)
	product, err := h.productService.GetProduct(c.Request.Context(), id, language)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Parse request
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Update product fields
	product.Name = req.Name
	product.Description = req.Description

	// Convert string categories to ObjectIDs
	categoryIDs := make([]primitive.ObjectID, len(req.Categories))
	for i, catStr := range req.Categories {
		catID, err := primitive.ObjectIDFromHex(catStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid category ID: " + catStr})
			return
		}
		categoryIDs[i] = catID
	}
	product.Categories = categoryIDs
	product.Attributes = req.Attributes
	product.SKU = req.SKU
	product.Price = req.Price
	product.StockQuantity = req.StockQuantity
	product.Images = req.Images

	// Save changes
	err = h.productService.UpdateProduct(c.Request.Context(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := h.productService.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListProducts godoc
// @Summary List products
// @Description Get a list of products with optional filtering and pagination
// @Tags products
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// For now, we're not implementing complex filtering
	filter := map[string]interface{}{}

	language := utils.GetLanguageFromRequest(c)
	products, total, err := h.productService.ListProducts(c.Request.Context(), filter, page, limit, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + limit - 1) / limit,
	})
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Description Get a list of products in a specific category
// @Tags products
// @Produce json
// @Param category path string true "Category name"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/category/{category} [get]
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	language := utils.GetLanguageFromRequest(c)
	products, total, err := h.productService.GetProductsByCategory(c.Request.Context(), category, page, limit, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + limit - 1) / limit,
	})
}

// UpdateStock godoc
// @Summary Update product stock
// @Description Update the stock quantity of a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param stock body StockUpdateRequest true "Stock update details"
// @Success 200 {string} string "Stock updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id}/stock [put]
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	productID := c.Param("id")

	var req StockUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.productService.UpdateStock(c.Request.Context(), productID, req.Quantity)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}
