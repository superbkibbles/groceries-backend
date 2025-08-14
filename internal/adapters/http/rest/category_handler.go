package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	categoryService ports.CategoryService
}

// NewCategoryHandler creates a new category handler and registers routes
func NewCategoryHandler(router *gin.RouterGroup, categoryService ports.CategoryService) {
	handler := &CategoryHandler{
		categoryService: categoryService,
	}

	categories := router.Group("/categories")
	{
		categories.POST("", handler.CreateCategory)
		categories.GET("", handler.ListCategories)
		categories.GET("/root", handler.GetRootCategories)
		categories.GET("/:id", handler.GetCategory)
		categories.GET("/slug/:slug", handler.GetCategoryBySlug)
		categories.PUT("/:id", handler.UpdateCategory)
		categories.DELETE("/:id", handler.DeleteCategory)
		categories.GET("/:id/children", handler.GetChildCategories)
		categories.GET("/:id/tree", handler.GetCategoryTree)
		categories.GET("/:id/products", handler.GetProductsByCategory)
	}
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	ParentID    string `json:"parent_id,omitempty"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category with the provided details
// @Tags categories
// @Accept json
// @Produce json
// @Param category body CreateCategoryRequest true "Category details"
// @Success 201 {object} entities.Category
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	category, err := h.categoryService.CreateCategory(c.Request.Context(), req.Name, req.Description, req.Slug, req.ParentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategory godoc
// @Summary Get a category by ID
// @Description Get detailed information about a category by its ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} entities.Category
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")

	category, err := h.categoryService.GetCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategoryBySlug godoc
// @Summary Get a category by slug
// @Description Get detailed information about a category by its slug
// @Tags categories
// @Produce json
// @Param slug path string true "Category Slug"
// @Success 200 {object} entities.Category
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.categoryService.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category with the provided details
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body UpdateCategoryRequest true "Updated category details"
// @Success 200 {object} entities.Category
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	// Get existing category
	existing, err := h.categoryService.GetCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Update fields
	existing.Update(req.Name, req.Description, req.Slug)

	// Save changes
	err = h.categoryService.UpdateCategory(c.Request.Context(), existing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by its ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	err := h.categoryService.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListCategories godoc
// @Summary List categories
// @Description Get a paginated list of categories with optional filtering
// @Tags categories
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Create empty filter for now (can be extended later)
	filter := map[string]any{}

	categories, total, err := h.categoryService.ListCategories(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       categories,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + limit - 1) / limit,
	})
}

// GetRootCategories godoc
// @Summary Get root categories
// @Description Get all top-level categories (with no parent)
// @Tags categories
// @Produce json
// @Success 200 {array} entities.Category
// @Failure 500 {object} ErrorResponse
// @Router /categories/root [get]
func (h *CategoryHandler) GetRootCategories(c *gin.Context) {
	categories, err := h.categoryService.GetRootCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetChildCategories godoc
// @Summary Get child categories
// @Description Get all direct child categories of a parent category
// @Tags categories
// @Produce json
// @Param id path string true "Parent Category ID"
// @Success 200 {array} entities.Category
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id}/children [get]
func (h *CategoryHandler) GetChildCategories(c *gin.Context) {
	parentID := c.Param("id")

	categories, err := h.categoryService.GetChildCategories(c.Request.Context(), parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryTree godoc
// @Summary Get category tree
// @Description Get a complete category tree starting from a root category
// @Tags categories
// @Produce json
// @Param id path string true "Root Category ID"
// @Success 200 {object} entities.Category
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id}/tree [get]
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	rootID := c.Param("id")

	category, err := h.categoryService.GetCategoryTree(c.Request.Context(), rootID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Description Get products in a category and optionally its subcategories
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Param include_subcategories query bool false "Include products from subcategories"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} PaginatedResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id}/products [get]
func (h *CategoryHandler) GetProductsByCategory(c *gin.Context) {
	categoryID := c.Param("id")
	includeSubcategories, _ := strconv.ParseBool(c.DefaultQuery("include_subcategories", "false"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.categoryService.GetProductsByCategory(c.Request.Context(), categoryID, includeSubcategories, page, limit)
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
