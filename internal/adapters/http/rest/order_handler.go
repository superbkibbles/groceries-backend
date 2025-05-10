package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	orderService ports.OrderService
}

// NewOrderHandler creates a new order handler and registers routes
func NewOrderHandler(router *gin.RouterGroup, orderService ports.OrderService) {
	handler := &OrderHandler{
		orderService: orderService,
	}

	orders := router.Group("/orders")
	{
		orders.POST("", handler.CreateOrder)
		orders.GET("", handler.ListOrders)
		orders.GET("/:id", handler.GetOrder)
		orders.GET("/customer/:customerId", handler.GetCustomerOrders)

		// Order items
		orders.POST("/:id/items", handler.AddItem)
		orders.PUT("/:id/items", handler.UpdateItemQuantity)
		orders.DELETE("/:id/items", handler.RemoveItem)

		// Order status and info
		orders.PUT("/:id/status", handler.UpdateOrderStatus)
		orders.PUT("/:id/payment", handler.SetPaymentInfo)
		orders.PUT("/:id/tracking", handler.SetTrackingInfo)
	}
}

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	CustomerID   string                `json:"customer_id" binding:"required"`
	ShippingInfo entities.ShippingInfo `json:"shipping_info" binding:"required"`
}

// OrderItemRequest represents the request body for adding/updating an item
type OrderItemRequest struct {
	ProductID   string `json:"product_id" binding:"required"`
	VariationID string `json:"variation_id" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,gt=0"`
}

// RemoveItemRequest represents the request body for removing an item
type RemoveItemRequest struct {
	ProductID   string `json:"product_id" binding:"required"`
	VariationID string `json:"variation_id" binding:"required"`
}

// UpdateStatusRequest represents the request body for updating order status
type UpdateStatusRequest struct {
	Status entities.OrderStatus `json:"status" binding:"required"`
}

// PaymentInfoRequest represents the request body for setting payment info
type PaymentInfoRequest struct {
	Method        string  `json:"method" binding:"required"`
	TransactionID string  `json:"transaction_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
}

// TrackingInfoRequest represents the request body for setting tracking info
type TrackingInfoRequest struct {
	Carrier     string `json:"carrier" binding:"required"`
	TrackingNum string `json:"tracking_num" binding:"required"`
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order for a customer
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "Order details"
// @Success 201 {object} entities.Order
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), req.CustomerID, req.ShippingInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get detailed information about an order by its ID
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} entities.Order
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListOrders godoc
// @Summary List orders
// @Description Get a list of orders with optional filtering and pagination
// @Tags orders
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// For now, we're not implementing complex filtering
	filter := map[string]interface{}{}

	orders, total, err := h.orderService.ListOrders(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       orders,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + limit - 1) / limit,
	})
}

// GetCustomerOrders godoc
// @Summary Get orders for a customer
// @Description Get a list of orders for a specific customer
// @Tags orders
// @Produce json
// @Param customerId path string true "Customer ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/customer/{customerId} [get]
func (h *OrderHandler) GetCustomerOrders(c *gin.Context) {
	customerID := c.Param("customerId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.orderService.GetCustomerOrders(c.Request.Context(), customerID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       orders,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + limit - 1) / limit,
	})
}

// AddItem godoc
// @Summary Add an item to an order
// @Description Add a product variation to an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param item body OrderItemRequest true "Item details"
// @Success 200 {string} string "Item added successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/items [post]
func (h *OrderHandler) AddItem(c *gin.Context) {
	orderID := c.Param("id")

	var req OrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.orderService.AddItem(
		c.Request.Context(),
		orderID,
		req.ProductID,
		req.VariationID,
		req.Quantity,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "order not found" || err.Error() == "product not found" || err.Error() == "variation not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "insufficient stock" || err.Error() == "cannot modify a non-pending order" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added successfully"})
}

// UpdateItemQuantity godoc
// @Summary Update item quantity
// @Description Update the quantity of an item in an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param item body OrderItemRequest true "Item details with new quantity"
// @Success 200 {string} string "Item quantity updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/items [put]
func (h *OrderHandler) UpdateItemQuantity(c *gin.Context) {
	orderID := c.Param("id")

	var req OrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.orderService.UpdateItemQuantity(
		c.Request.Context(),
		orderID,
		req.ProductID,
		req.VariationID,
		req.Quantity,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "order not found" || err.Error() == "item not found in order" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "insufficient stock" || err.Error() == "cannot modify a non-pending order" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item quantity updated successfully"})
}

// RemoveItem godoc
// @Summary Remove an item from an order
// @Description Remove a product variation from an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param item body RemoveItemRequest true "Item to remove"
// @Success 200 {string} string "Item removed successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/items [delete]
func (h *OrderHandler) RemoveItem(c *gin.Context) {
	orderID := c.Param("id")

	var req RemoveItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.orderService.RemoveItem(
		c.Request.Context(),
		orderID,
		req.ProductID,
		req.VariationID,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "order not found" || err.Error() == "item not found in order" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "cannot modify a non-pending order" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed successfully"})
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body UpdateStatusRequest true "New status"
// @Success 200 {string} string "Order status updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.orderService.UpdateOrderStatus(c.Request.Context(), orderID, req.Status)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "order not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "invalid status transition" || err.Error() == "cannot change status of a delivered or cancelled order" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// SetPaymentInfo godoc
// @Summary Set payment information
// @Description Set payment information for an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param payment body PaymentInfoRequest true "Payment information"
// @Success 200 {string} string "Payment information set successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/payment [put]
func (h *OrderHandler) SetPaymentInfo(c *gin.Context) {
	orderID := c.Param("id")

	var req PaymentInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.orderService.SetPaymentInfo(
		c.Request.Context(),
		orderID,
		req.Method,
		req.TransactionID,
		req.Amount,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "order not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "cannot set payment info for a non-pending order" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment information set successfully"})
}

// SetTrackingInfo godoc
// @Summary Set tracking information
// @Description Set shipping tracking information for an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param tracking body TrackingInfoRequest true "Tracking information"
// @Success 200 {string} string "Tracking information set successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/tracking [put]
func (h *OrderHandler) SetTrackingInfo(c *gin.Context) {
	orderID := c.Param("id")

	var req TrackingInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.orderService.SetTrackingInfo(
		c.Request.Context(),
		orderID,
		req.Carrier,
		req.TrackingNum,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "order not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "cannot set tracking info for an order that is not paid or shipped" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tracking information set successfully"})
}
