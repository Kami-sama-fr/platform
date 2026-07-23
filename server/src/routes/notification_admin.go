package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type NotificationAdminHandler struct {
	deps Dependencies
}

func NewNotificationAdminHandler(deps Dependencies) *NotificationAdminHandler {
	return &NotificationAdminHandler{deps: deps}
}

func (h *NotificationAdminHandler) ListTemplates(c *gin.Context) {
	items, err := h.deps.NotificationAdminService.ListTemplates(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *NotificationAdminHandler) CreateTemplate(c *gin.Context) {
	var req struct {
		Type    string `json:"type"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
		IsHTML  bool   `json:"isHtml"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Type == "" || req.Subject == "" || req.Body == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.NotificationAdminService.CreateTemplate(c.Request.Context(), req.Type, req.Subject, req.Body, req.IsHTML)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *NotificationAdminHandler) UpdateTemplate(c *gin.Context) {
	id := c.Param("templateId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	var req struct {
		Subject *string `json:"subject"`
		Body    *string `json:"body"`
		IsHTML  *bool   `json:"isHtml"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.NotificationAdminService.UpdateTemplate(c.Request.Context(), id, req.Subject, req.Body, req.IsHTML)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *NotificationAdminHandler) DeleteTemplate(c *gin.Context) {
	id := c.Param("templateId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.NotificationAdminService.DeleteTemplate(c.Request.Context(), id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

func (h *NotificationAdminHandler) Send(c *gin.Context) {
	var req struct {
		Type          string   `json:"type"`
		Title         string   `json:"title"`
		Body          string   `json:"body"`
		TargetUserIDs []string `json:"targetUserIds"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Type == "" || req.Title == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	sent, err := h.deps.NotificationAdminService.Send(c.Request.Context(), req.Type, req.Title, req.Body, req.TargetUserIDs)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"sent": sent})
}

func (h *NotificationAdminHandler) GetHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, total, err := h.deps.NotificationAdminService.GetHistory(c.Request.Context(), page, limit)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}
