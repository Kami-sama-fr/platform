package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type FAQHandler struct {
	deps Dependencies
}

func NewFAQHandler(deps Dependencies) *FAQHandler {
	return &FAQHandler{deps: deps}
}

func (h *FAQHandler) List(c *gin.Context) {
	category := c.DefaultQuery("category", "")
	activeOnly := c.DefaultQuery("activeOnly", "false") == "true"
	items, err := h.deps.FAQService.List(c.Request.Context(), category, activeOnly)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *FAQHandler) Create(c *gin.Context) {
	var req struct {
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		Category  string `json:"category"`
		SortOrder int    `json:"sortOrder"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Question == "" || req.Answer == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.FAQService.Create(c.Request.Context(), req.Question, req.Answer, req.Category, req.SortOrder)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *FAQHandler) GetByID(c *gin.Context) {
	id := c.Param("faqId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.FAQService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *FAQHandler) Update(c *gin.Context) {
	id := c.Param("faqId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	var req struct {
		Question  *string `json:"question"`
		Answer    *string `json:"answer"`
		Category  *string `json:"category"`
		SortOrder *int    `json:"sortOrder"`
		IsActive  *bool   `json:"isActive"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.FAQService.Update(c.Request.Context(), id, req.Question, req.Answer, req.Category, req.SortOrder, req.IsActive)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *FAQHandler) Delete(c *gin.Context) {
	id := c.Param("faqId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.FAQService.Delete(c.Request.Context(), id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

func (h *FAQHandler) Reorder(c *gin.Context) {
	var req struct {
		Orders []struct {
			ID    string `json:"id"`
			Order int    `json:"order"`
		} `json:"orders"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	orders := make(map[string]int, len(req.Orders))
	for _, o := range req.Orders {
		orders[o.ID] = o.Order
	}
	if err := h.deps.FAQService.Reorder(c.Request.Context(), orders); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"reordered": true})
}

func (h *FAQHandler) ListPublic(c *gin.Context) {
	category := c.DefaultQuery("category", "")
	items, err := h.deps.FAQService.List(c.Request.Context(), category, true)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *FAQHandler) GetByIDPublic(c *gin.Context) {
	id := c.Param("faqId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.FAQService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}
