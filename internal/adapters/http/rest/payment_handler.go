package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// PaymentHandler handles HTTP requests for payment settings
type PaymentHandler struct {
	paymentService ports.PaymentService
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(paymentService ports.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// RegisterRoutes registers the payment routes
func (h *PaymentHandler) RegisterRoutes(router *gin.Engine) {
	paymentGroup := router.Group("/api/v1/payments")
	{
		// Payment methods
		paymentGroup.GET("/methods", h.ListPaymentMethods)
		paymentGroup.GET("/methods/:id", h.GetPaymentMethod)
		paymentGroup.POST("/methods", h.CreatePaymentMethod)
		paymentGroup.PUT("/methods/:id", h.UpdatePaymentMethod)
		paymentGroup.DELETE("/methods/:id", h.DeletePaymentMethod)

		// Payment gateways
		paymentGroup.GET("/gateways", h.ListPaymentGateways)
		paymentGroup.GET("/gateways/:id", h.GetPaymentGateway)
		paymentGroup.POST("/gateways", h.CreatePaymentGateway)
		paymentGroup.PUT("/gateways/:id", h.UpdatePaymentGateway)
		paymentGroup.DELETE("/gateways/:id", h.DeletePaymentGateway)

		// Customer payment methods
		paymentGroup.GET("/customer/:customerId/methods", h.ListCustomerPaymentMethods)
		paymentGroup.GET("/customer/methods/:id", h.GetCustomerPaymentMethod)
		paymentGroup.POST("/customer/methods", h.CreateCustomerPaymentMethod)
		paymentGroup.PUT("/customer/methods/:id", h.UpdateCustomerPaymentMethod)
		paymentGroup.DELETE("/customer/methods/:id", h.DeleteCustomerPaymentMethod)
		paymentGroup.PUT("/customer/methods/:id/default", h.SetDefaultPaymentMethod)

		// Payment processing
		paymentGroup.POST("/process", h.ProcessPayment)
		paymentGroup.GET("/verify/:transactionId", h.VerifyPayment)
		paymentGroup.POST("/refund", h.RefundPayment)
	}
}

// CreatePaymentMethodRequest represents a request to create a payment method
type CreatePaymentMethodRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description" binding:"required"`
	Type        string                 `json:"type" binding:"required"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// ListPaymentMethods godoc
// @Summary List payment methods
// @Description Get a list of available payment methods
// @Tags payments
// @Accept json
// @Produce json
// @Param active query bool false "Filter by active status"
// @Success 200 {array} entities.PaymentMethod
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/methods [get]
func (h *PaymentHandler) ListPaymentMethods(c *gin.Context) {
	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active", "true"))

	methods, err := h.paymentService.ListPaymentMethods(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, methods)
}

// GetPaymentMethod godoc
// @Summary Get payment method
// @Description Get a payment method by ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment Method ID"
// @Success 200 {object} entities.PaymentMethod
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/methods/{id} [get]
func (h *PaymentHandler) GetPaymentMethod(c *gin.Context) {
	id := c.Param("id")

	method, err := h.paymentService.GetPaymentMethod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, method)
}

// CreatePaymentMethod godoc
// @Summary Create payment method
// @Description Create a new payment method
// @Tags payments
// @Accept json
// @Produce json
// @Param method body CreatePaymentMethodRequest true "Payment Method Details"
// @Success 201 {object} entities.PaymentMethod
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/methods [post]
func (h *PaymentHandler) CreatePaymentMethod(c *gin.Context) {
	var req CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Convert string type to PaymentMethodType
	methodType := entities.PaymentMethodType(req.Type)

	method, err := h.paymentService.CreatePaymentMethod(
		c.Request.Context(),
		req.Name,
		req.Description,
		methodType,
		req.Config,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, method)
}

// UpdatePaymentMethod godoc
// @Summary Update payment method
// @Description Update an existing payment method
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment Method ID"
// @Param method body CreatePaymentMethodRequest true "Payment Method Details"
// @Success 200 {object} entities.PaymentMethod
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/methods/{id} [put]
func (h *PaymentHandler) UpdatePaymentMethod(c *gin.Context) {
	id := c.Param("id")

	var req CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get existing method
	method, err := h.paymentService.GetPaymentMethod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Update fields
	method.Name = req.Name
	method.Description = req.Description
	method.Type = entities.PaymentMethodType(req.Type)
	method.Config = req.Config

	// Save changes
	err = h.paymentService.UpdatePaymentMethod(c.Request.Context(), method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, method)
}

// DeletePaymentMethod godoc
// @Summary Delete payment method
// @Description Delete a payment method
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment Method ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/methods/{id} [delete]
func (h *PaymentHandler) DeletePaymentMethod(c *gin.Context) {
	id := c.Param("id")

	err := h.paymentService.DeletePaymentMethod(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "payment method not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method deleted successfully"})
}

// CreatePaymentGatewayRequest represents a request to create a payment gateway
type CreatePaymentGatewayRequest struct {
	Name     string                 `json:"name" binding:"required"`
	Provider string                 `json:"provider" binding:"required"`
	Config   map[string]interface{} `json:"config,omitempty"`
}

// ListPaymentGateways godoc
// @Summary List payment gateways
// @Description Get a list of payment gateways
// @Tags payments
// @Accept json
// @Produce json
// @Param active query bool false "Filter by active status"
// @Success 200 {array} entities.PaymentGateway
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/gateways [get]
func (h *PaymentHandler) ListPaymentGateways(c *gin.Context) {
	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active", "true"))

	gateways, err := h.paymentService.ListPaymentGateways(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gateways)
}

// GetPaymentGateway godoc
// @Summary Get payment gateway
// @Description Get a payment gateway by ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment Gateway ID"
// @Success 200 {object} entities.PaymentGateway
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/gateways/{id} [get]
func (h *PaymentHandler) GetPaymentGateway(c *gin.Context) {
	id := c.Param("id")

	gateway, err := h.paymentService.GetPaymentGateway(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gateway)
}

// CreatePaymentGateway godoc
// @Summary Create payment gateway
// @Description Create a new payment gateway
// @Tags payments
// @Accept json
// @Produce json
// @Param gateway body CreatePaymentGatewayRequest true "Payment Gateway Details"
// @Success 201 {object} entities.PaymentGateway
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/gateways [post]
func (h *PaymentHandler) CreatePaymentGateway(c *gin.Context) {
	var req CreatePaymentGatewayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	gateway, err := h.paymentService.CreatePaymentGateway(
		c.Request.Context(),
		req.Name,
		req.Provider,
		req.Config,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gateway)
}

// UpdatePaymentGateway godoc
// @Summary Update payment gateway
// @Description Update an existing payment gateway
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment Gateway ID"
// @Param gateway body CreatePaymentGatewayRequest true "Payment Gateway Details"
// @Success 200 {object} entities.PaymentGateway
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/gateways/{id} [put]
func (h *PaymentHandler) UpdatePaymentGateway(c *gin.Context) {
	id := c.Param("id")

	var req CreatePaymentGatewayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get existing gateway
	gateway, err := h.paymentService.GetPaymentGateway(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Update fields
	gateway.Name = req.Name
	gateway.Provider = req.Provider
	gateway.Config = req.Config

	// Save changes
	err = h.paymentService.UpdatePaymentGateway(c.Request.Context(), gateway)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gateway)
}

// DeletePaymentGateway godoc
// @Summary Delete payment gateway
// @Description Delete a payment gateway
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment Gateway ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/gateways/{id} [delete]
func (h *PaymentHandler) DeletePaymentGateway(c *gin.Context) {
	id := c.Param("id")

	err := h.paymentService.DeletePaymentGateway(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "payment gateway not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment gateway deleted successfully"})
}

// CreateCustomerPaymentMethodRequest represents a request to create a customer payment method
type CreateCustomerPaymentMethodRequest struct {
	CustomerID      string `json:"customer_id" binding:"required"`
	PaymentMethodID string `json:"payment_method_id" binding:"required"`
	Token           string `json:"token" binding:"required"`
	Last4           string `json:"last4,omitempty"`
	ExpiryMonth     int    `json:"expiry_month,omitempty"`
	ExpiryYear      int    `json:"expiry_year,omitempty"`
	IsDefault       bool   `json:"is_default"`
}

// ListCustomerPaymentMethods godoc
// @Summary List customer payment methods
// @Description Get a list of payment methods for a customer
// @Tags payments
// @Accept json
// @Produce json
// @Param customerId path string true "Customer ID"
// @Success 200 {array} entities.CustomerPaymentMethod
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/customer/{customerId}/methods [get]
func (h *PaymentHandler) ListCustomerPaymentMethods(c *gin.Context) {
	customerID := c.Param("customerId")

	methods, err := h.paymentService.ListCustomerPaymentMethods(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, methods)
}

// GetCustomerPaymentMethod godoc
// @Summary Get customer payment method
// @Description Get a customer payment method by ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Customer Payment Method ID"
// @Success 200 {object} entities.CustomerPaymentMethod
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/customer/methods/{id} [get]
func (h *PaymentHandler) GetCustomerPaymentMethod(c *gin.Context) {
	id := c.Param("id")

	method, err := h.paymentService.GetCustomerPaymentMethod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, method)
}

// CreateCustomerPaymentMethod godoc
// @Summary Create customer payment method
// @Description Create a new payment method for a customer
// @Tags payments
// @Accept json
// @Produce json
// @Param method body CreateCustomerPaymentMethodRequest true "Customer Payment Method Details"
// @Success 201 {object} entities.CustomerPaymentMethod
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/customer/methods [post]
func (h *PaymentHandler) CreateCustomerPaymentMethod(c *gin.Context) {
	var req CreateCustomerPaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	method, err := h.paymentService.CreateCustomerPaymentMethod(
		c.Request.Context(),
		req.CustomerID,
		req.PaymentMethodID,
		req.Token,
		req.Last4,
		req.ExpiryMonth,
		req.ExpiryYear,
		req.IsDefault,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, method)
}

// UpdateCustomerPaymentMethod godoc
// @Summary Update customer payment method
// @Description Update an existing customer payment method
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Customer Payment Method ID"
// @Param method body CreateCustomerPaymentMethodRequest true "Customer Payment Method Details"
// @Success 200 {object} entities.CustomerPaymentMethod
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/customer/methods/{id} [put]
func (h *PaymentHandler) UpdateCustomerPaymentMethod(c *gin.Context) {
	id := c.Param("id")

	var req CreateCustomerPaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get existing method
	method, err := h.paymentService.GetCustomerPaymentMethod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Update fields
	method.CustomerID = req.CustomerID
	method.PaymentMethodID = req.PaymentMethodID
	method.Token = req.Token
	method.Last4 = req.Last4
	method.ExpiryMonth = req.ExpiryMonth
	method.ExpiryYear = req.ExpiryYear
	method.IsDefault = req.IsDefault

	// Save changes
	err = h.paymentService.UpdateCustomerPaymentMethod(c.Request.Context(), method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, method)
}

// DeleteCustomerPaymentMethod godoc
// @Summary Delete customer payment method
// @Description Delete a customer payment method
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Customer Payment Method ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/customer/methods/{id} [delete]
func (h *PaymentHandler) DeleteCustomerPaymentMethod(c *gin.Context) {
	id := c.Param("id")

	err := h.paymentService.DeleteCustomerPaymentMethod(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "customer payment method not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer payment method deleted successfully"})
}

// SetDefaultPaymentMethod godoc
// @Summary Set default payment method
// @Description Set a payment method as the default for a customer
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Customer Payment Method ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/customer/methods/{id}/default [put]
func (h *PaymentHandler) SetDefaultPaymentMethod(c *gin.Context) {
	id := c.Param("id")

	// Get the method to get the customer ID
	method, err := h.paymentService.GetCustomerPaymentMethod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	err = h.paymentService.SetDefaultCustomerPaymentMethod(c.Request.Context(), method.CustomerID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default payment method set successfully"})
}

// ProcessPaymentRequest represents a request to process a payment
type ProcessPaymentRequest struct {
	OrderID         string  `json:"order_id" binding:"required"`
	PaymentMethodID string  `json:"payment_method_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
}

// ProcessPayment godoc
// @Summary Process payment
// @Description Process a payment for an order
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body ProcessPaymentRequest true "Payment Details"
// @Success 200 {object} entities.PaymentInfo
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/process [post]
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	paymentInfo, err := h.paymentService.ProcessPayment(
		c.Request.Context(),
		req.OrderID,
		req.PaymentMethodID,
		req.Amount,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentInfo)
}

// VerifyPayment godoc
// @Summary Verify payment
// @Description Verify a payment transaction
// @Tags payments
// @Accept json
// @Produce json
// @Param transactionId path string true "Transaction ID"
// @Success 200 {object} map[string]bool
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/verify/{transactionId} [get]
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	transactionID := c.Param("transactionId")

	verified, err := h.paymentService.VerifyPayment(c.Request.Context(), transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"verified": verified})
}

// RefundPaymentRequest represents a request to refund a payment
type RefundPaymentRequest struct {
	OrderID string  `json:"order_id" binding:"required"`
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	Reason  string  `json:"reason" binding:"required"`
}

// RefundPayment godoc
// @Summary Refund payment
// @Description Refund a payment for an order
// @Tags payments
// @Accept json
// @Produce json
// @Param refund body RefundPaymentRequest true "Refund Details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/refund [post]
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	var req RefundPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.paymentService.RefundPayment(
		c.Request.Context(),
		req.OrderID,
		req.Amount,
		req.Reason,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment refunded successfully"})
}
