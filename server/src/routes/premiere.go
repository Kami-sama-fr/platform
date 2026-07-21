package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/middleware"
	"github.com/kami-sama-fr/platform/server/src/services"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type PremiereHandler struct {
	deps Dependencies
}

func NewPremiereHandler(deps Dependencies) *PremiereHandler {
	return &PremiereHandler{deps: deps}
}

func (h *PremiereHandler) List(c *gin.Context) {
	items, err := h.deps.PremiereService.List(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *PremiereHandler) Create(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req struct {
		AnimeID     string `json:"animeId"`
		Season      int    `json:"season"`
		EpisodeNum  int    `json:"episodeNum"`
		ScheduledAt string `json:"scheduledAt"`
		Region      string `json:"region"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.PremiereService.Create(c.Request.Context(), principal.UserID, req.AnimeID, req.Season, req.EpisodeNum, req.ScheduledAt, req.Region)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *PremiereHandler) GetByID(c *gin.Context) {
	id := c.Param("premiereId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.PremiereService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *PremiereHandler) Update(c *gin.Context) {
	id := c.Param("premiereId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req services.UpdatePremiereInput
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.PremiereService.Update(c.Request.Context(), principal.UserID, id, req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *PremiereHandler) Delete(c *gin.Context) {
	id := c.Param("premiereId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.PremiereService.Delete(c.Request.Context(), id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}
