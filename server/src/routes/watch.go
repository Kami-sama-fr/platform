package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/middleware"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type WatchHandler struct {
	deps Dependencies
}

func NewWatchHandler(deps Dependencies) *WatchHandler {
	return &WatchHandler{deps: deps}
}

func (h *WatchHandler) GetProgress(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	episodeID := c.Param("episodeId")
	if episodeID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.WatchService.GetProgress(c.Request.Context(), principal.UserID, episodeID)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *WatchHandler) UpsertProgress(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	episodeID := c.Param("episodeId")
	if episodeID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	var req struct {
		AnimeID   string  `json:"animeId"`
		Progress  float64 `json:"progress"`
		Duration  float64 `json:"duration"`
		Percentage float64 `json:"percentage"`
		Completed bool    `json:"completed"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.WatchService.UpsertProgress(c.Request.Context(), principal.UserID, episodeID, req.AnimeID, req.Progress, req.Duration, req.Percentage, req.Completed)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *WatchHandler) ListProgress(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	items, err := h.deps.WatchService.ListProgress(c.Request.Context(), principal.UserID)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *WatchHandler) ContinueWatching(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}
	items, err := h.deps.WatchService.GetContinueWatching(c.Request.Context(), principal.UserID, limit)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *WatchHandler) ListHistory(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.deps.WatchService.ListHistory(c.Request.Context(), principal.UserID, limit, offset)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *WatchHandler) AddHistory(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req struct {
		EpisodeID string  `json:"episodeId"`
		AnimeID   string  `json:"animeId"`
		Duration  float64 `json:"duration"`
	}
	if c.ShouldBindJSON(&req) != nil || req.EpisodeID == "" || req.AnimeID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.WatchService.AddHistory(c.Request.Context(), principal.UserID, req.EpisodeID, req.AnimeID, req.Duration)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}
