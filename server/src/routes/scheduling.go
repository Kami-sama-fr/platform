package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/middleware"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type SchedulingHandler struct {
	deps Dependencies
}

func NewSchedulingHandler(deps Dependencies) *SchedulingHandler {
	return &SchedulingHandler{deps: deps}
}

// --- Simulcasts ---

func (h *SchedulingHandler) ListSimulcasts(c *gin.Context) {
	items, err := h.deps.SchedulingService.ListActiveSimulcasts(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SchedulingHandler) GetSimulcast(c *gin.Context) {
	id := c.Param("simulcastId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SchedulingService.GetSimulcast(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SchedulingHandler) GetSimulcastByWeek(c *gin.Context) {
	items, err := h.deps.SchedulingService.ListSimulcastsByWeek(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SchedulingHandler) CreateSimulcast(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req struct {
		AnimeID       string `json:"animeId"`
		AirDay        string `json:"airDay"`
		AirTime       string `json:"airTime"`
		AirTimezone   string `json:"airTimezone"`
		Region        string `json:"region"`
		Platform      string `json:"platform"`
		EpisodeNumber int    `json:"episodeNumber"`
		Season        int    `json:"season"`
		StartsAt      string `json:"startsAt"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SchedulingService.CreateSimulcast(c.Request.Context(), principal.UserID, req.AnimeID, req.AirDay, req.AirTime, req.AirTimezone, req.Region, req.Platform, req.EpisodeNumber, req.Season, req.StartsAt)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *SchedulingHandler) UpdateSimulcast(c *gin.Context) {
	id := c.Param("simulcastId")
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
		AirDay      *string `json:"airDay"`
		AirTime     *string `json:"airTime"`
		AirTimezone *string `json:"airTimezone"`
		Region      *string `json:"region"`
		Platform    *string `json:"platform"`
		Status      *string `json:"status"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SchedulingService.UpdateSimulcast(c.Request.Context(), principal.UserID, id, req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SchedulingHandler) DeleteSimulcast(c *gin.Context) {
	id := c.Param("simulcastId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.SchedulingService.DeleteSimulcast(c.Request.Context(), id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

// --- Release Schedules ---

func (h *SchedulingHandler) ListUpcomingReleases(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	items, err := h.deps.SchedulingService.ListUpcomingReleases(c.Request.Context(), limit)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SchedulingHandler) CreateReleaseSchedule(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req struct {
		AnimeID     string `json:"animeId"`
		EpisodeID   string `json:"episodeId"`
		ScheduledAt string `json:"scheduledAt"`
		Region      string `json:"region"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SchedulingService.CreateReleaseSchedule(c.Request.Context(), principal.UserID, req.AnimeID, &req.EpisodeID, req.ScheduledAt, req.Region)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *SchedulingHandler) UpdateReleaseSchedule(c *gin.Context) {
	id := c.Param("scheduleId")
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
		ScheduledAt *string `json:"scheduledAt"`
		Status      *string `json:"status"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SchedulingService.UpdateReleaseSchedule(c.Request.Context(), principal.UserID, id, req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}
