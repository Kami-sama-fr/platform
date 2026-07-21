package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type AnalyticsHandler struct {
	deps Dependencies
}

func NewAnalyticsHandler(deps Dependencies) *AnalyticsHandler {
	return &AnalyticsHandler{deps: deps}
}

func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	data, err := h.deps.AnalyticsService.GetOverview(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, data)
}

func (h *AnalyticsHandler) GetOverviewByPeriod(c *gin.Context) {
	period := c.DefaultQuery("period", "day")
	items, err := h.deps.AnalyticsService.GetOverviewByPeriod(c.Request.Context(), period)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetWatchTime(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetWatchTime(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetWatchTimeByAnime(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetWatchTimeByAnime(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetWatchTimeByEpisode(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetWatchTimeByEpisode(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetWatchTimeHistogram(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetWatchTimeHistogram(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetDevices(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetDevices(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetDevicesBrowsers(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetDevicesBrowsers(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetDevicesOS(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetDevicesOS(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetPopular(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetPopular(c.Request.Context(), 20)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetPopularTrending(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetPopularTrending(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetPopularNew(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetPopularNew(c.Request.Context(), 20)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetGeography(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetGeography(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetGeographyTopCountries(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetGeographyTopCountries(c.Request.Context(), 10)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetActiveUsers(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetActiveUsers(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetActiveUsersRetention(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetActiveUsersRetention(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *AnalyticsHandler) GetActiveUsersSessions(c *gin.Context) {
	items, err := h.deps.AnalyticsService.GetActiveUsersSessions(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}
