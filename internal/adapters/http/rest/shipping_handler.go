package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// ShippingHandler handles HTTP requests for shipping settings
type ShippingHandler struct {
	shippingService ports.ShippingService
}

// NewShippingHandler creates a new ShippingHandler
func NewShippingHandler(shippingService ports.ShippingService) *ShippingHandler {
	return &ShippingHandler{
		shippingService: shippingService,
	}
}

// RegisterRoutes registers the shipping routes
func (h *ShippingHandler) RegisterRoutes(router *gin.Engine) {
	shippingGroup := router.Group("/api/v1/shipping")
	{
		// Shipping methods
		shippingGroup.GET("/methods", h.ListShippingMethods)
		shippingGroup.GET("/methods/:id", h.GetShippingMethod)
		shippingGroup.POST("/methods", h.CreateShippingMethod)
		shippingGroup.PUT("/methods/:id", h.UpdateShippingMethod)
		shippingGroup.DELETE("/methods/:id", h.DeleteShippingMethod)

		// Shipping zones
		shippingGroup.GET("/zones", h.ListShippingZones)
		shippingGroup.GET("/zones/:id", h.GetShippingZone)
		shippingGroup.POST("/zones", h.CreateShippingZone)
		shippingGroup.PUT("/zones/:id", h.UpdateShippingZone)
		shippingGroup.DELETE("/zones/:id", h.DeleteShippingZone)

		// Shipping rates
		shippingGroup.GET("/rates", h.ListShippingRates)
		shippingGroup.GET("/rates/:id", h.GetShippingRate)
		shippingGroup.POST("/rates", h.CreateShippingRate)
		shippingGroup.PUT("/rates/:id", h.UpdateShippingRate)
		shippingGroup.DELETE("/rates/:id", h.DeleteShippingRate)

		// Calculate shipping
		shippingGroup.POST("/calculate", h.CalculateShipping)
	}
}

// CreateShippingMethodRequest represents a request to create a shipping method
type CreateShippingMethodRequest struct {
	Name                  string  `json:"name" binding:"required"`
	Description           string  `json:"description" binding:"required"`
	BasePrice             float64 `json:"base_price" binding:"required,min=0"`
	EstimatedDeliveryDays int     `json:"estimated_delivery_days" binding:"required,min=1"`
}

// ListShippingMethods godoc
// @Summary List shipping methods
// @Description Get a list of available shipping methods
// @Tags shipping
// @Accept json
// @Produce json
// @Param active query bool false "Filter by active status"
// @Success 200 {array} entities.ShippingMethod
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/methods [get]
func (h *ShippingHandler) ListShippingMethods(c *gin.Context) {
	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active", "true"))

	methods, err := h.shippingService.ListShippingMethods(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, methods)
}

// GetShippingMethod godoc
// @Summary Get shipping method
// @Description Get a shipping method by ID
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "Shipping Method ID"
// @Success 200 {object} entities.ShippingMethod
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/methods/{id} [get]
func (h *ShippingHandler) GetShippingMethod(c *gin.Context) {
	id := c.Param("id")

	method, err := h.shippingService.GetShippingMethod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, method)
}

// CreateShippingMethod godoc
// @Summary Create shipping method
// @Description Create a new shipping method
// @Tags shipping
// @Accept json
// @Produce json
// @Param method body CreateShippingMethodRequest true "Shipping Method Details"
// @Success 201 {object} entities.ShippingMethod
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/methods [post]
func (h *ShippingHandler) CreateShippingMethod(c *gin.Context) {
	var req CreateShippingMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	method, err := h.shippingService.CreateShippingMethod(
		c.Request.Context(),
		req.Name,
		req.Description,
		req.BasePrice,
		req.EstimatedDeliveryDays,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, method)
}

// UpdateShippingMethod godoc
// @Summary Update shipping method
// @Description Update an existing shipping method
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "Shipping Method ID"
// @Param method body CreateShippingMethodRequest true "Shipping Method Details"
// @Success 200 {object} entities.ShippingMethod
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/methods/{id} [put]
func (h *ShippingHandler) UpdateShippingMethod(c *gin.Context) {
	id := c.Param("id")

	var req CreateShippingMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get existing method
	method, err := h.shippingService.GetShippingMethod(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Update fields
	method.Name = req.Name
	method.Description = req.Description
	method.BasePrice = req.BasePrice
	method.EstimatedDeliveryDays = req.EstimatedDeliveryDays

	// Save changes
	err = h.shippingService.UpdateShippingMethod(c.Request.Context(), method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, method)
}

// DeleteShippingMethod godoc
// @Summary Delete shipping method
// @Description Delete a shipping method
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "Shipping Method ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/methods/{id} [delete]
func (h *ShippingHandler) DeleteShippingMethod(c *gin.Context) {
	id := c.Param("id")

	err := h.shippingService.DeleteShippingMethod(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "shipping method not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping method deleted successfully"})
}

// CreateShippingZoneRequest represents a request to create a shipping zone
type CreateShippingZoneRequest struct {
	Name      string   `json:"name" binding:"required"`
	Countries []string `json:"countries" binding:"required"`
}

// ListShippingZones godoc
// @Summary List shipping zones
// @Description Get a list of shipping zones
// @Tags shipping
// @Accept json
// @Produce json
// @Success 200 {array} entities.ShippingZone
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/zones [get]
func (h *ShippingHandler) ListShippingZones(c *gin.Context) {
	zones, err := h.shippingService.ListShippingZones(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, zones)
}

// GetShippingZone godoc
// @Summary Get shipping zone
// @Description Get a shipping zone by ID
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "Shipping Zone ID"
// @Success 200 {object} entities.ShippingZone
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/zones/{id} [get]
func (h *ShippingHandler) GetShippingZone(c *gin.Context) {
	id := c.Param("id")

	zone, err := h.shippingService.GetShippingZone(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, zone)
}

// CreateShippingZone godoc
// @Summary Create shipping zone
// @Description Create a new shipping zone
// @Tags shipping
// @Accept json
// @Produce json
// @Param zone body CreateShippingZoneRequest true "Shipping Zone Details"
// @Success 201 {object} entities.ShippingZone
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/zones [post]
func (h *ShippingHandler) CreateShippingZone(c *gin.Context) {
	var req CreateShippingZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	zone, err := h.shippingService.CreateShippingZone(
		c.Request.Context(),
		req.Name,
		req.Countries,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, zone)
}

// UpdateShippingZone godoc
// @Summary Update shipping zone
// @Description Update an existing shipping zone
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "Shipping Zone ID"
// @Param zone body CreateShippingZoneRequest true "Shipping Zone Details"
// @Success 200 {object} entities.ShippingZone
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/zones/{id} [put]
func (h *ShippingHandler) UpdateShippingZone(c *gin.Context) {
	id := c.Param("id")

	var req CreateShippingZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get existing zone
	zone, err := h.shippingService.GetShippingZone(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Update fields
	zone.Name = req.Name
	zone.Countries = req.Countries

	// Save changes
	err = h.shippingService.UpdateShippingZone(c.Request.Context(), zone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, zone)
}

// DeleteShippingZone godoc
// @Summary Delete shipping zone
// @Description Delete a shipping zone
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "Shipping Zone ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/zones/{id} [delete]
func (h *ShippingHandler) DeleteShippingZone(c *gin.Context) {
	id := c.Param("id")

	err := h.shippingService.DeleteShippingZone(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "shipping zone not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping zone deleted successfully"})
}

// CreateShippingRateRequest represents a request to create a shipping rate
type CreateShippingRateRequest struct {
	ShippingZoneID   string  `json:"shipping_zone_id" binding:"required"`
	ShippingMethodID string  `json:"shipping_method_id" binding:"required"`
	Price            float64 `json:"price" binding:"required,min=0"`
}

// ListShippingRates godoc
// @Summary List shipping rates
// @Description Get a list of shipping rates
// @Tags shipping
// @Accept json
// @Produce json
// @Param zone_id query string false "Filter by Zone ID"
// @Param method_id query string false "Filter by Method ID"
// @Success 200 {array} entities.ShippingRate
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/rates [get]
func (h *ShippingHandler) ListShippingRates(c *gin.Context) {
	zoneID := c.Query("zone_id")
	methodID := c.Query("method_id")

	var rates []*entities.ShippingRate
	var err error

	if zoneID != "" {
		rates, err = h.shippingService.GetShippingRatesByZone(c.Request.Context(), zoneID)
	} else if methodID != "" {
		rates, err = h.shippingService.GetShippingRatesByMethod(c.Request.Context(), methodID)
	} else {
		// Get all rates (not implemented in the service interface)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "must provide zone_id or method_id"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}

// GetShippingRate godoc
// @Summary Get shipping rate
// @Description Get a shipping rate by ID
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "Shipping Rate ID"
// @Success 200 {object} entities.ShippingRate
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/rates/{id} [get]
func (h *ShippingHandler) GetShippingRate(c *gin.Context) {
	id := c.Param("id")

	rate, err := h.shippingService.GetShippingRate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, rate)
}

// CreateShippingRate godoc
// @Summary Create shipping rate
// @Description Create a new shipping rate
// @Tags shipping
// @Accept json
// @Produce json
// @Param rate body CreateShippingRateRequest true "Shipping Rate Details"
// @Success 201 {object} entities.ShippingRate
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/rates [post]
func (h *ShippingHandler) CreateShippingRate(c *gin.Context) {
	var req CreateShippingRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	rate, err := h.shippingService.CreateShippingRate(
		c.Request.Context(),
		req.ShippingZoneID,
		req.ShippingMethodID,
		req.Price,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rate)
}

// UpdateShippingRate godoc
// @Summary Update shipping rate
// @Description Update an existing shipping rate
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "Shipping Rate ID"
// @Param rate body CreateShippingRateRequest true "Shipping Rate Details"
// @Success 200 {object} entities.ShippingRate
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/rates/{id} [put]
func (h *ShippingHandler) UpdateShippingRate(c *gin.Context) {
	id := c.Param("id")

	var req CreateShippingRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get existing rate
	rate, err := h.shippingService.GetShippingRate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Update fields
	rate.ShippingZoneID = req.ShippingZoneID
	rate.ShippingMethodID = req.ShippingMethodID
	rate.Price = req.Price

	// Save changes
	err = h.shippingService.UpdateShippingRate(c.Request.Context(), rate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, rate)
}

// DeleteShippingRate godoc
// @Summary Delete shipping rate
// @Description Delete a shipping rate
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "Shipping Rate ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/rates/{id} [delete]
func (h *ShippingHandler) DeleteShippingRate(c *gin.Context) {
	id := c.Param("id")

	err := h.shippingService.DeleteShippingRate(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "shipping rate not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping rate deleted successfully"})
}

// CalculateShippingRequest represents a request to calculate shipping costs
type CalculateShippingRequest struct {
	CountryCode string               `json:"country_code" binding:"required"`
	Items       []*entities.CartItem `json:"items" binding:"required"`
}

// CalculateShipping godoc
// @Summary Calculate shipping costs
// @Description Calculate shipping costs for a set of items to a specific country
// @Tags shipping
// @Accept json
// @Produce json
// @Param request body CalculateShippingRequest true "Calculation Request"
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shipping/calculate [post]
func (h *ShippingHandler) CalculateShipping(c *gin.Context) {
	var req CalculateShippingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	shippingOptions, err := h.shippingService.CalculateShippingCost(
		c.Request.Context(),
		req.CountryCode,
		req.Items,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, shippingOptions)
}
