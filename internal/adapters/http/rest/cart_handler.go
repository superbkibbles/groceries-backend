package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// CartHandler handles HTTP requests related to shopping carts
type CartHandler struct {
	cartService ports.CartService
}

// NewCartHandler creates a new CartHandler
func NewCartHandler(cartService ports.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// RegisterRoutes registers the cart routes
func (h *CartHandler) RegisterRoutes(router *gin.Engine) {
	cartGroup := router.Group("/api/v1/cart")
	{
		cartGroup.GET("", h.GetCart)
		cartGroup.POST("/items", h.AddItem)
		cartGroup.PUT("/items/:itemId", h.UpdateItemQuantity)
		cartGroup.DELETE("/items/:itemId", h.RemoveItem)
		cartGroup.DELETE("/items", h.ClearCart)
		cartGroup.POST("/checkout", h.Checkout)
	}
}

// GetCart godoc
// @Summary Get user's cart
// @Description Retrieves the current user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} entities.Cart
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Get user's cart
	cart, err := h.cartService.GetUserCart(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddItemRequest represents a request to add an item to the cart
type AddItemRequest struct {
	ProductID   string `json:"product_id" binding:"required"`
	VariationID string `json:"variation_id" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
}

// AddItem godoc
// @Summary Add item to cart
// @Description Adds a product to the user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Param item body AddItemRequest true "Item to add"
// @Success 200 {object} entities.CartItem
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/cart/items [post]
func (h *CartHandler) AddItem(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Parse request body
	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Add item to cart
	item, err := h.cartService.AddItem(c.Request.Context(), userID.(string), req.ProductID, req.VariationID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateQuantityRequest represents a request to update an item's quantity
type UpdateQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// UpdateItemQuantity godoc
// @Summary Update item quantity
// @Description Updates the quantity of an item in the user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Param itemId path string true "Item ID"
// @Param quantity body UpdateQuantityRequest true "New quantity"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/cart/items/{itemId} [put]
func (h *CartHandler) UpdateItemQuantity(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Get item ID from path
	itemID := c.Param("itemId")

	// Parse request body
	var req UpdateQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Update item quantity
	err := h.cartService.UpdateItemQuantity(c.Request.Context(), userID.(string), itemID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item quantity updated"})
}

// RemoveItem godoc
// @Summary Remove item from cart
// @Description Removes an item from the user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Param itemId path string true "Item ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/cart/items/{itemId} [delete]
func (h *CartHandler) RemoveItem(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Get item ID from path
	itemID := c.Param("itemId")

	// Remove item from cart
	err := h.cartService.RemoveItem(c.Request.Context(), userID.(string), itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

// ClearCart godoc
// @Summary Clear cart
// @Description Removes all items from the user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/cart/items [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Clear cart
	err := h.cartService.ClearCart(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared"})
}

// CheckoutRequest represents a request to checkout a cart
type CheckoutRequest struct {
	ShippingInfo entities.ShippingInfo `json:"shipping_info" binding:"required"`
}

// Checkout godoc
// @Summary Checkout cart
// @Description Converts the user's cart to an order
// @Tags cart
// @Accept json
// @Produce json
// @Param checkout body CheckoutRequest true "Checkout information"
// @Success 200 {object} entities.Order
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/cart/checkout [post]
func (h *CartHandler) Checkout(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Parse request body
	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Convert cart to order
	order, err := h.cartService.ConvertToOrder(c.Request.Context(), userID.(string), req.ShippingInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
