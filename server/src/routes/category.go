package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/middleware"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type CategoryHandler struct {
	deps Dependencies
}

func NewCategoryHandler(deps Dependencies) *CategoryHandler {
	return &CategoryHandler{deps: deps}
}

func (h *CategoryHandler) List(c *gin.Context) {
	items, err := h.deps.CategoryService.List(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("categoryId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CategoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CategoryHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CategoryService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req struct {
		Name        string  `json:"name"`
		Slug        string  `json:"slug"`
		Description string  `json:"description"`
		ParentID    *string `json:"parentId"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CategoryService.Create(c.Request.Context(), principal.UserID, req.Name, req.Slug, req.Description, req.ParentID)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("categoryId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req struct {
		Name        *string `json:"name"`
		Slug        *string `json:"slug"`
		Description *string `json:"description"`
		ParentID    *string `json:"parentId"`
		SortOrder   *int    `json:"sortOrder"`
		IsActive    *bool   `json:"isActive"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CategoryService.Update(c.Request.Context(), principal.UserID, id, req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("categoryId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.CategoryService.Delete(c.Request.Context(), id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}
