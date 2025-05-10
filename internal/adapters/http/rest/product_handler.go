package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
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

		// Variation routes
		products.POST("/:id/variations", handler.AddVariation)
		products.PUT("/:id/variations/:variationId", handler.UpdateVariation)
		products.DELETE("/:id/variations/:variationId", handler.RemoveVariation)
		products.PUT("/:id/variations/:variationId/stock", handler.UpdateStock)
	}
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	BasePrice   float64  `json:"base_price" binding:"required,gt=0"`
	Categories  []string `json:"categories"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	BasePrice   float64  `json:"base_price" binding:"required,gt=0"`
	Categories  []string `json:"categories"`
}

// VariationRequest represents the request body for adding/updating a variation
type VariationRequest struct {
	Attributes    map[string]interface{} `json:"attributes" binding:"required"`
	SKU           string                 `json:"sku" binding:"required"`
	Price         float64                `json:"price" binding:"required,gt=0"`
	StockQuantity int                    `json:"stock_quantity" binding:"required,gte=0"`
	Images        []string               `json:"images"`
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
// @Success 201 {object} entities.Product
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), req.Name, req.Description, req.BasePrice, req.Categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get detailed information about a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} entities.Product
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
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

	// Get existing product
	product, err := h.productService.GetProduct(c.Request.Context(), id)
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
	product.BasePrice = req.BasePrice
	product.Categories = req.Categories

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

	products, total, err := h.productService.ListProducts(c.Request.Context(), filter, page, limit)
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

	products, total, err := h.productService.GetProductsByCategory(c.Request.Context(), category, page, limit)
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

// AddVariation godoc
// @Summary Add a variation to a product
// @Description Add a new variation to an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param variation body VariationRequest true "Variation details"
// @Success 201 {object} entities.Variation
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id}/variations [post]
func (h *ProductHandler) AddVariation(c *gin.Context) {
	productID := c.Param("id")

	var req VariationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	variation, err := h.productService.AddVariation(
		c.Request.Context(),
		productID,
		req.Attributes,
		req.SKU,
		req.Price,
		req.StockQuantity,
		req.Images,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, variation)
}

// UpdateVariation godoc
// @Summary Update a product variation
// @Description Update an existing variation of a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param variationId path string true "Variation ID"
// @Param variation body VariationRequest true "Updated variation details"
// @Success 200 {string} string "Variation updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id}/variations/{variationId} [put]
func (h *ProductHandler) UpdateVariation(c *gin.Context) {
	productID := c.Param("id")
	variationID := c.Param("variationId")

	var req VariationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.productService.UpdateVariation(
		c.Request.Context(),
		productID,
		variationID,
		req.Attributes,
		req.SKU,
		req.Price,
		req.StockQuantity,
		req.Images,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "product not found" || err.Error() == "variation not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Variation updated successfully"})
}

// RemoveVariation godoc
// @Summary Remove a product variation
// @Description Remove a variation from a product
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Param variationId path string true "Variation ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id}/variations/{variationId} [delete]
func (h *ProductHandler) RemoveVariation(c *gin.Context) {
	productID := c.Param("id")
	variationID := c.Param("variationId")

	err := h.productService.RemoveVariation(c.Request.Context(), productID, variationID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "product not found" || err.Error() == "variation not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateStock godoc
// @Summary Update variation stock
// @Description Update the stock quantity of a product variation
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param variationId path string true "Variation ID"
// @Param stock body StockUpdateRequest true "Stock update details"
// @Success 200 {string} string "Stock updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id}/variations/{variationId}/stock [put]
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	productID := c.Param("id")
	variationID := c.Param("variationId")

	var req StockUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.productService.UpdateStock(c.Request.Context(), productID, variationID, req.Quantity)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "product not found" || err.Error() == "variation not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}
