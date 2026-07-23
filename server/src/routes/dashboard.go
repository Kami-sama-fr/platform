package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type DashboardHandler struct {
	deps Dependencies
}

func NewDashboardHandler(deps Dependencies) *DashboardHandler {
	return &DashboardHandler{deps: deps}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.deps.DashboardService.GetStats(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, stats)
}

func (h *DashboardHandler) GetWeeklyViews(c *gin.Context) {
	items, err := h.deps.DashboardService.GetWeeklyViews(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *DashboardHandler) GetSubscriptionDistribution(c *gin.Context) {
	items, err := h.deps.DashboardService.GetSubscriptionDistribution(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *DashboardHandler) GetTopAnime(c *gin.Context) {
	items, err := h.deps.DashboardService.GetTopAnime(c.Request.Context(), 10)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *DashboardHandler) GetRecentUploads(c *gin.Context) {
	items, err := h.deps.DashboardService.GetRecentUploads(c.Request.Context(), 10)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}
