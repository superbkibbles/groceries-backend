package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// WishlistHandler handles HTTP requests related to wishlists
type WishlistHandler struct {
	wishlistService ports.WishlistService
}

// NewWishlistHandler creates a new WishlistHandler
func NewWishlistHandler(wishlistService ports.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}

// RegisterRoutes registers the wishlist routes
func (h *WishlistHandler) RegisterRoutes(router *gin.Engine) {
	wishlistGroup := router.Group("/api/v1/wishlist")
	{
		wishlistGroup.GET("", h.GetWishlist)
		wishlistGroup.POST("/items", h.AddItem)
		wishlistGroup.DELETE("/items/:itemId", h.RemoveItem)
		wishlistGroup.DELETE("/items", h.ClearWishlist)
		wishlistGroup.GET("/check/:productId", h.CheckProduct)
	}
}

// GetWishlist godoc
// @Summary Get user's wishlist
// @Description Retrieves the current user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Success 200 {object} entities.Wishlist
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/wishlist [get]
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Get user's wishlist
	wishlist, err := h.wishlistService.GetUserWishlist(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, wishlist)
}

// WishlistAddItemRequest represents a request to add an item to the wishlist
type WishlistAddItemRequest struct {
	ProductID   string `json:"product_id" binding:"required"`
	VariationID string `json:"variation_id,omitempty"`
}

// AddItem godoc
// @Summary Add item to wishlist
// @Description Adds a product to the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Param item body AddItemRequest true "Item to add"
// @Success 200 {object} entities.WishlistItem
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/wishlist/items [post]
func (h *WishlistHandler) AddItem(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Parse request body
	var req WishlistAddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Add item to wishlist
	item, err := h.wishlistService.AddItem(c.Request.Context(), userID.(string), req.ProductID, req.VariationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// RemoveItem godoc
// @Summary Remove item from wishlist
// @Description Removes an item from the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Param itemId path string true "Item ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/wishlist/items/{itemId} [delete]
func (h *WishlistHandler) RemoveItem(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Get item ID from path
	itemID := c.Param("itemId")

	// Remove item from wishlist
	err := h.wishlistService.RemoveItem(c.Request.Context(), userID.(string), itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from wishlist"})
}

// ClearWishlist godoc
// @Summary Clear wishlist
// @Description Removes all items from the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/wishlist/items [delete]
func (h *WishlistHandler) ClearWishlist(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Clear wishlist
	err := h.wishlistService.ClearWishlist(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wishlist cleared"})
}

// CheckProductRequest represents a request to check if a product is in the wishlist
type CheckProductRequest struct {
	ProductID   string `json:"product_id" binding:"required"`
	VariationID string `json:"variation_id,omitempty"`
}

// CheckProduct godoc
// @Summary Check if product is in wishlist
// @Description Checks if a product is in the user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Param variationId query string false "Variation ID"
// @Success 200 {object} map[string]bool
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/wishlist/check/{productId} [get]
func (h *WishlistHandler) CheckProduct(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Get product ID from path
	productID := c.Param("productId")
	// Get variation ID from query
	variationID := c.Query("variationId")

	// Check if product is in wishlist
	isInWishlist, err := h.wishlistService.IsProductInWishlist(c.Request.Context(), userID.(string), productID, variationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"in_wishlist": isInWishlist})
}
