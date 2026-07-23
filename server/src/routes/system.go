package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type SystemHandler struct {
	deps Dependencies
}

func NewSystemHandler(deps Dependencies) *SystemHandler {
	return &SystemHandler{deps: deps}
}

// --- Health ---

func (h *SystemHandler) HealthServices(c *gin.Context) {
	result := h.deps.SystemService.GetServiceHealth(c.Request.Context())
	utils.Success(c, http.StatusOK, result)
}

func (h *SystemHandler) HealthUptime(c *gin.Context) {
	result := h.deps.SystemService.GetUptime()
	utils.Success(c, http.StatusOK, result)
}

func (h *SystemHandler) HealthMetrics(c *gin.Context) {
	result := h.deps.SystemService.GetMetrics()
	utils.Success(c, http.StatusOK, result)
}

// --- Logs ---

func (h *SystemHandler) ListLogs(c *gin.Context) {
	level := c.Query("level")
	source := c.Query("source")
	date := c.Query("date")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	items, total, err := h.deps.SystemService.GetLogs(c.Request.Context(), level, source, date, limit, offset)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *SystemHandler) GetLogByID(c *gin.Context) {
	logID := c.Param("logId")
	if logID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SystemService.GetLogByID(c.Request.Context(), logID)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SystemHandler) SearchLogs(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	items, err := h.deps.SystemService.SearchLogs(c.Request.Context(), query)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

// --- Queue ---

func (h *SystemHandler) GetQueueStatus(c *gin.Context) {
	result := h.deps.SystemService.GetQueueStatus(c.Request.Context())
	utils.Success(c, http.StatusOK, result)
}

func (h *SystemHandler) ListQueueJobs(c *gin.Context) {
	status := c.Query("status")
	jobType := c.Query("type")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	items, total, err := h.deps.SystemService.ListQueueJobs(c.Request.Context(), status, jobType, limit, offset)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *SystemHandler) RetryQueueJob(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "operation accepted"})
}

func (h *SystemHandler) CancelQueueJob(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "operation accepted"})
}

func (h *SystemHandler) FlushQueue(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "operation accepted"})
}

// --- Cache ---

func (h *SystemHandler) GetCacheStatus(c *gin.Context) {
	result := h.deps.SystemService.GetCacheStatus(c.Request.Context())
	utils.Success(c, http.StatusOK, result)
}

func (h *SystemHandler) FlushCache(c *gin.Context) {
	if err := h.deps.SystemService.FlushCache(c.Request.Context()); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "cache flushed"})
}

func (h *SystemHandler) FlushCacheByPattern(c *gin.Context) {
	pattern := c.Param("pattern")
	if pattern == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.SystemService.FlushCacheByPattern(c.Request.Context(), pattern); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "cache pattern flushed"})
}

func (h *SystemHandler) ListCacheKeys(c *gin.Context) {
	keys, err := h.deps.SystemService.ListCacheKeys(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": keys})
}

func (h *SystemHandler) DeleteCacheKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.SystemService.DeleteCacheKey(c.Request.Context(), key); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

// --- Search ---

func (h *SystemHandler) GetSearchStatus(c *gin.Context) {
	result := h.deps.SystemService.GetSearchStatus(c.Request.Context())
	utils.Success(c, http.StatusOK, result)
}

func (h *SystemHandler) TriggerReindex(c *gin.Context) {
	if err := h.deps.SystemService.TriggerReindex(c.Request.Context()); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "reindex triggered"})
}

func (h *SystemHandler) ListSearchIndexes(c *gin.Context) {
	items, err := h.deps.SystemService.ListSearchIndexes(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SystemHandler) GetSearchIndexStats(c *gin.Context) {
	indexName := c.Param("indexName")
	if indexName == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SystemService.GetSearchIndexStats(c.Request.Context(), indexName)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SystemHandler) UpdateSearchIndex(c *gin.Context) {
	indexName := c.Param("indexName")
	if indexName == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.SystemService.UpdateSearchIndex(c.Request.Context(), indexName); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "index update triggered"})
}

// --- Background Jobs ---

func (h *SystemHandler) ListBackgroundJobs(c *gin.Context) {
	items, err := h.deps.SystemService.ListBackgroundJobs(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SystemHandler) RunBackgroundJob(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.SystemService.RunBackgroundJob(c.Request.Context(), jobID); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "job triggered"})
}

func (h *SystemHandler) PauseBackgroundJob(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.SystemService.PauseBackgroundJob(c.Request.Context(), jobID); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "job paused"})
}

func (h *SystemHandler) ResumeBackgroundJob(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.SystemService.ResumeBackgroundJob(c.Request.Context(), jobID); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "job resumed"})
}

func (h *SystemHandler) GetBackgroundJobHistory(c *gin.Context) {
	items, err := h.deps.SystemService.GetBackgroundJobHistory(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}
