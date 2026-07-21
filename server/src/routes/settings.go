package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/middleware"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type SettingsHandler struct {
	deps Dependencies
}

func NewSettingsHandler(deps Dependencies) *SettingsHandler {
	return &SettingsHandler{deps: deps}
}

func (h *SettingsHandler) GetByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SettingsService.GetSetting(c.Request.Context(), key)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SettingsHandler) List(c *gin.Context) {
	category := c.Query("category")
	items, err := h.deps.SettingsService.ListSettings(c.Request.Context(), category)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SettingsHandler) Upsert(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	_ = principal
	var req struct {
		Key         string `json:"key"`
		Value       any    `json:"value"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Key == "" || req.Category == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SettingsService.UpsertSetting(c.Request.Context(), req.Key, req.Value, req.Category, req.Description)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SettingsHandler) Delete(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.SettingsService.DeleteSetting(c.Request.Context(), key); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

// --- SEO ---

func (h *SettingsHandler) GetSeoMeta(c *gin.Context) {
	pagePath := c.Query("pagePath")
	if pagePath == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SettingsService.GetSeoMeta(c.Request.Context(), pagePath)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SettingsHandler) UpsertSeoMeta(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	_ = principal
	var req struct {
		PagePath     string `json:"pagePath"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		Keywords     string `json:"keywords"`
		OgImage      string `json:"ogImage"`
		CanonicalUrl string `json:"canonicalUrl"`
	}
	if c.ShouldBindJSON(&req) != nil || req.PagePath == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SettingsService.UpsertSeoMeta(c.Request.Context(), req.PagePath, req.Title, req.Description, req.Keywords, req.OgImage, req.CanonicalUrl)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}
