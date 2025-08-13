package rest

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService  ports.UserService
	orderService ports.OrderService
}

// NewUserHandler creates a new user handler and registers routes
func NewUserHandler(router *gin.RouterGroup, userService ports.UserService, orderService ports.OrderService) {
	handler := &UserHandler{
		userService:  userService,
		orderService: orderService,
	}

	users := router.Group("/users")
	{
		// Authentication routes
		users.POST("/register", handler.Register)
		users.POST("/login", handler.Login)
		users.POST("/send-otp", handler.SendOTP)
		// Current user route (protected by auth middleware)
		users.GET("/me", AuthRequired(), handler.GetCurrentUser)

		// User listing route (admin only)
		users.GET("", AuthRequired(), handler.ListUsers)

		// User profile routes
		users.GET("/:id", handler.GetUser)
		users.PUT("/:id", handler.UpdateUser)
		users.PUT("/:id/password", handler.ChangePassword)

		// Address routes
		users.GET("/:id/addresses", handler.GetAddresses)
		users.POST("/:id/addresses", handler.AddAddress)
		users.PUT("/:id/addresses/:addressId", handler.UpdateAddress)
		users.DELETE("/:id/addresses/:addressId", handler.DeleteAddress)
		users.PUT("/:id/addresses/:addressId/default", handler.SetDefaultAddress)

		// User orders route
		users.GET("/:id/orders", handler.GetUserOrders)
	}
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	// Email    string `json:"email" binding:"required,email"`
	// Password string `json:"password" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// LoginRequestAdmin represents the request body for admin login
type LoginRequestAdmin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response for a successful login
type LoginResponse struct {
	User  *entities.User `json:"user"`
	Token string         `json:"token"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// ChangePasswordRequest represents the request body for changing a password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// AddressRequest represents the request body for adding/updating an address
type AddressRequest struct {
	Name         string `json:"name" binding:"required"`
	AddressLine1 string `json:"address_line1" binding:"required"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city" binding:"required"`
	State        string `json:"state" binding:"required"`
	Country      string `json:"country" binding:"required"`
	PostalCode   string `json:"postal_code" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	IsDefault    bool   `json:"is_default"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} entities.User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userService.Register(
		c.Request.Context(),
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already registered" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	// Don't return the password hash
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login a user
// @Description Authenticate a user and return a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	log.Println("Login triggered")
	// check if user is admin login with email and password
	userRole := c.Request.Header.Get("user_role")
	log.Println("userRole", userRole)
	if userRole != "" && userRole == string(entities.UserRoleAdmin) {
		// Handle admin login with email and password
		var req LoginRequestAdmin
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		user, token, err := h.userService.LoginAdmin(
			c.Request.Context(),
			req.Email,
			req.Password,
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			User:  user,
			Token: token,
		})
	} else {
		// Handle user login with phone number
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		user, token, err := h.userService.Login(
			c.Request.Context(),
			req.PhoneNumber,
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
			return
		}

		// Don't return the password hash
		user.PasswordHash = ""

		c.JSON(http.StatusOK, LoginResponse{
			User:  user,
			Token: token,
		})
	}

}

// Login godoc
// @Summary Send OTP
// @Description Saves OTP Inside redis and send send it to phone number
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Send OTP credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/send-otp [post]
func (h *UserHandler) SendOTP(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.userService.SendOTP(c.Request.Context(), req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get detailed information about a user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} entities.User
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Don't return the password hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UpdateUserRequest true "Updated user details"
// @Success 200 {object} entities.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// Get existing user
	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Parse request
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Update user fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName

	// Save changes
	err = h.userService.UpdateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Don't return the password hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, user)
}

// ChangePassword godoc
// @Summary Change a user's password
// @Description Change a user's password
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param passwords body ChangePasswordRequest true "Password change details"
// @Success 200 {string} string "Password changed successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id := c.Param("id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.userService.ChangePassword(
		c.Request.Context(),
		id,
		req.CurrentPassword,
		req.NewPassword,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "current password is incorrect" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// GetAddresses godoc
// @Summary Get user addresses
// @Description Get all addresses for a user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {array} entities.Address
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/addresses [get]
func (h *UserHandler) GetAddresses(c *gin.Context) {
	id := c.Param("id")

	addresses, err := h.userService.GetAddresses(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// AddAddress godoc
// @Summary Add a new address
// @Description Add a new address for a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param address body AddressRequest true "Address details"
// @Success 201 {object} entities.Address
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/addresses [post]
func (h *UserHandler) AddAddress(c *gin.Context) {
	id := c.Param("id")

	var req AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	address, err := h.userService.AddAddress(
		c.Request.Context(),
		id,
		req.Name,
		req.AddressLine1,
		req.AddressLine2,
		req.City,
		req.State,
		req.Country,
		req.PostalCode,
		req.Phone,
		req.IsDefault,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, address)
}

// UpdateAddress godoc
// @Summary Update an address
// @Description Update an existing address for a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param addressId path string true "Address ID"
// @Param address body AddressRequest true "Updated address details"
// @Success 200 {string} string "Address updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/addresses/{addressId} [put]
func (h *UserHandler) UpdateAddress(c *gin.Context) {
	userID := c.Param("id")
	addressIDStr := c.Param("addressId")
	addressID, err := primitive.ObjectIDFromHex(addressIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid address ID"})
		return
	}

	// Get existing addresses to verify ownership
	addresses, err := h.userService.GetAddresses(c.Request.Context(), userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	// Verify address belongs to user
	var existingAddress *entities.Address
	for _, addr := range addresses {
		if addr.ID == addressID {
			existingAddress = addr
			break
		}
	}

	if existingAddress == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "address not found"})
		return
	}

	// Parse request
	var req AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Update address fields
	existingAddress.Name = req.Name
	existingAddress.AddressLine1 = req.AddressLine1
	existingAddress.AddressLine2 = req.AddressLine2
	existingAddress.City = req.City
	existingAddress.State = req.State
	existingAddress.Country = req.Country
	existingAddress.PostalCode = req.PostalCode
	existingAddress.Phone = req.Phone
	existingAddress.IsDefault = req.IsDefault

	// Save changes
	err = h.userService.UpdateAddress(c.Request.Context(), existingAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address updated successfully"})
}

// DeleteAddress godoc
// @Summary Delete an address
// @Description Delete an address for a user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Param addressId path string true "Address ID"
// @Success 200 {string} string "Address deleted successfully"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/addresses/{addressId} [delete]
func (h *UserHandler) DeleteAddress(c *gin.Context) {
	userID := c.Param("id")
	addressIDStr := c.Param("addressId")
	addressID, err := primitive.ObjectIDFromHex(addressIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid address ID"})
		return
	}

	// Get existing addresses to verify ownership
	addresses, err := h.userService.GetAddresses(c.Request.Context(), userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	// Verify address belongs to user
	addressExists := false
	for _, addr := range addresses {
		if addr.ID == addressID {
			addressExists = true
			break
		}
	}

	if !addressExists {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "address not found"})
		return
	}

	// Delete address
	err = h.userService.DeleteAddress(c.Request.Context(), addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}

// SetDefaultAddress godoc
// @Summary Set default address
// @Description Set an address as the default for a user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Param addressId path string true "Address ID"
// @Success 200 {string} string "Default address set successfully"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/addresses/{addressId}/default [put]
func (h *UserHandler) SetDefaultAddress(c *gin.Context) {
	userID := c.Param("id")
	addressIDStr := c.Param("addressId")
	addressID, err := primitive.ObjectIDFromHex(addressIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid address ID"})
		return
	}

	// Get existing addresses to verify ownership
	addresses, err := h.userService.GetAddresses(c.Request.Context(), userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	// Verify address belongs to user
	addressExists := false
	for _, addr := range addresses {
		if addr.ID == addressID {
			addressExists = true
			break
		}
	}

	if !addressExists {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "address not found"})
		return
	}

	// Set as default
	err = h.userService.SetDefaultAddress(c.Request.Context(), userID, addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default address set successfully"})
}

// GetCurrentUser godoc
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} entities.User
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Convert to string
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "invalid user ID format"})
		return
	}

	// Get user by ID
	user, err := h.userService.GetUser(c.Request.Context(), userIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}

	// Don't return the password hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, user)
}

// GetUserOrders godoc
// @Summary Get user orders
// @Description Get all orders for a specific user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} PaginatedResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/orders [get]
func (h *UserHandler) GetUserOrders(c *gin.Context) {
	// Get user ID from path parameter
	userID := c.Param("id")

	// Verify user exists
	_, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Get orders from order service
	orders, total, err := h.orderService.GetCustomerOrders(c.Request.Context(), userID, page, limit)
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

// ListUsers godoc
// @Summary List users with filtering
// @Description Get a list of users with optional role filtering and pagination
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param role query string false "Filter by user role (customer or admin)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Get user ID and role from context (set by auth middleware)
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized: authentication required"})
		return
	}

	userRole, roleExists := c.Get("user_role")

	// Parse query parameters
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = limitNum
		}
	}

	// Create filter
	filter := map[string]interface{}{}

	// Check if user is admin
	isAdmin := roleExists && userRole == string(entities.UserRoleAdmin)

	// Apply role filtering
	if role := c.Query("role"); role != "" {
		// Convert role string to UserRole type for proper filtering
		switch role {
		case string(entities.UserRoleCustomer):
			filter["role"] = entities.UserRoleCustomer
		case string(entities.UserRoleAdmin):
			// Only allow admin to see other admins
			if isAdmin {
				filter["role"] = entities.UserRoleAdmin
			} else {
				c.JSON(http.StatusForbidden, ErrorResponse{Error: "unauthorized: admin access required to view admin users"})
				return
			}
		default:
			// If role is not recognized, don't filter by role
		}
	} else if !isAdmin {
		// If no role specified and user is not admin, only show customers
		filter["role"] = entities.UserRoleCustomer
	}

	// Get users from service
	users, total, err := h.userService.ListUsers(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Don't return password hashes
	for _, user := range users {
		user.PasswordHash = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
