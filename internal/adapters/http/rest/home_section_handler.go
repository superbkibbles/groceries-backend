package rest

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/superbkibbles/ecommerce/internal/adapters/http/middleware"
    "github.com/superbkibbles/ecommerce/internal/domain/entities"
    "github.com/superbkibbles/ecommerce/internal/domain/ports"
)

type HomeSectionHandler struct {
    service ports.HomeSectionService
}

func NewHomeSectionHandler(router *gin.Engine, service ports.HomeSectionService) *HomeSectionHandler {
    h := &HomeSectionHandler{service: service}
    api := router.Group("/api/v1")
    sections := api.Group("/home-sections")
    {
        // Public listing for frontend
        sections.GET("", h.ListActive)
    }

    // Admin endpoints
    admin := api.Group("/admin/home-sections")
    admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
    {
        admin.POST("", h.Create)
        admin.GET(":id", h.Get)
        admin.PUT(":id", h.Update)
        admin.DELETE(":id", h.Delete)
        admin.GET("", h.List)
    }
    return h
}

type createHomeSectionRequest struct {
    Type        entities.HomeSectionType `json:"type" binding:"required,oneof=products categories"`
    Title       map[string]string        `json:"title" binding:"required"`
    ProductIDs  []string                 `json:"product_ids"`
    CategoryIDs []string                 `json:"category_ids"`
    Order       int                      `json:"order"`
    Active      bool                     `json:"active"`
}

func (h *HomeSectionHandler) Create(c *gin.Context) {
    var req createHomeSectionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error(), Code: http.StatusBadRequest})
        return
    }
    section, err := h.service.Create(c.Request.Context(), req.Type, req.Title, req.ProductIDs, req.CategoryIDs, req.Order, req.Active)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error(), Code: http.StatusBadRequest})
        return
    }
    c.JSON(http.StatusCreated, section)
}

func (h *HomeSectionHandler) Get(c *gin.Context) {
    id := c.Param("id")
    section, err := h.service.Get(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error(), Code: http.StatusNotFound})
        return
    }
    c.JSON(http.StatusOK, section)
}

func (h *HomeSectionHandler) Update(c *gin.Context) {
    id := c.Param("id")
    section, err := h.service.Get(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error(), Code: http.StatusNotFound})
        return
    }

    var req createHomeSectionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error(), Code: http.StatusBadRequest})
        return
    }
    section.Type = req.Type
    section.Title = req.Title
    // product/category IDs are replaced as a whole
    section.ProductIDs = section.ProductIDs[:0]
    section.CategoryIDs = section.CategoryIDs[:0]
    // let service convert and update IDs by reusing Create conversion path
    // quick rebuild through service to reuse validation
    // simpler approach: call Create-conversion logic inline here
    // For brevity, we call service.Create-like conversion by creating a temp
    tmp, err2 := h.service.Create(c.Request.Context(), req.Type, req.Title, req.ProductIDs, req.CategoryIDs, req.Order, req.Active)
    if err2 != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err2.Error(), Code: http.StatusBadRequest})
        return
    }
    section.ProductIDs = tmp.ProductIDs
    section.CategoryIDs = tmp.CategoryIDs
    section.Order = req.Order
    section.Active = req.Active

    if err := h.service.Update(c.Request.Context(), section); err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
        return
    }
    c.JSON(http.StatusOK, section)
}

func (h *HomeSectionHandler) Delete(c *gin.Context) {
    id := c.Param("id")
    if err := h.service.Delete(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *HomeSectionHandler) List(c *gin.Context) {
    sections, err := h.service.List(c.Request.Context(), false)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
        return
    }
    c.JSON(http.StatusOK, sections)
}

func (h *HomeSectionHandler) ListActive(c *gin.Context) {
    sections, err := h.service.List(c.Request.Context(), true)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
        return
    }
    c.JSON(http.StatusOK, sections)
}


