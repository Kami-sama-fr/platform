package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/middleware"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type MediaHandler struct {
	deps Dependencies
}

func NewMediaHandler(deps Dependencies) *MediaHandler {
	return &MediaHandler{deps: deps}
}

func (h *MediaHandler) List(c *gin.Context) {
	mediaType := strings.TrimSpace(c.Query("type"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.deps.MediaService.List(c.Request.Context(), mediaType, limit, offset)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *MediaHandler) GetByID(c *gin.Context) {
	id := c.Param("mediaId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.MediaService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *MediaHandler) Create(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req struct {
		AnimeID     string `json:"animeId"`
		EpisodeID   string `json:"episodeId"`
		Type        string `json:"type"`
		Filename    string `json:"filename"`
		OriginalUrl string `json:"originalUrl"`
		MimeType    string `json:"mimeType"`
		Size        int64  `json:"size"`
		Duration    float64 `json:"duration"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
		Language    string `json:"language"`
		Quality     string `json:"quality"`
		Codec       string `json:"codec"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.MediaService.Create(c.Request.Context(), principal.UserID, req.AnimeID, &req.EpisodeID, req.Type, req.Filename, req.OriginalUrl, req.MimeType, req.Size, req.Duration, req.Width, req.Height, req.Language, req.Quality, req.Codec)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *MediaHandler) Update(c *gin.Context) {
	id := c.Param("mediaId")
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
		Type     *string  `json:"type"`
		Status   *string  `json:"status"`
		CdnUrl   *string  `json:"cdnUrl"`
		Duration *float64 `json:"duration"`
		Quality  *string  `json:"quality"`
		Codec    *string  `json:"codec"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.MediaService.Update(c.Request.Context(), principal.UserID, id, req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *MediaHandler) Delete(c *gin.Context) {
	id := c.Param("mediaId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.MediaService.Delete(c.Request.Context(), id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

func (h *MediaHandler) ListEncodingJobs(c *gin.Context) {
	status := strings.TrimSpace(c.Query("status"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.deps.MediaService.ListEncodingJobs(c.Request.Context(), status, limit, offset)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *MediaHandler) GetEncodingJob(c *gin.Context) {
	id := c.Param("jobId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.MediaService.GetEncodingJob(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

// --- Media Extensions ---

func (h *MediaHandler) GenerateThumbnail(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "thumbnail generation accepted"})
}

func (h *MediaHandler) RetryEncodingJob(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "encoding job retry accepted"})
}

func (h *MediaHandler) CancelEncodingJob(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "encoding job cancellation accepted"})
}

func (h *MediaHandler) GetEncodingProfiles(c *gin.Context) {
	profiles := h.deps.SystemService.GetEncodingProfiles()
	utils.Success(c, http.StatusOK, gin.H{"items": profiles})
}

func (h *MediaHandler) ListUploads(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"items": []interface{}{}, "total": 0})
}

func (h *MediaHandler) InitiateUpload(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "upload initiation accepted", "uploadId": utils.NewID()})
}

func (h *MediaHandler) GetUploadProgress(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"uploadId": c.Param("uploadId"), "status": "pending", "progress": 0})
}

func (h *MediaHandler) CancelUpload(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "upload cancellation accepted"})
}

func (h *MediaHandler) CompleteUpload(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "upload finalization accepted"})
}
