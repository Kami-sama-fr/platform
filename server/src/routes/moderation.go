package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type ModerationHandler struct {
	deps Dependencies
}

func NewModerationHandler(deps Dependencies) *ModerationHandler {
	return &ModerationHandler{deps: deps}
}

func (h *ModerationHandler) GetQueue(c *gin.Context) {
	items, err := h.deps.ModerationService.GetQueue(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": len(items)})
}

func (h *ModerationHandler) GetItem(c *gin.Context) {
	id := c.Param("moderationId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	itemType := c.DefaultQuery("type", "")
	if itemType == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.ModerationService.GetItem(c.Request.Context(), itemType, id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *ModerationHandler) Approve(c *gin.Context) {
	id := c.Param("moderationId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	itemType := c.DefaultQuery("type", "")
	if itemType == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.ModerationService.Approve(c.Request.Context(), itemType, id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"approved": true})
}

func (h *ModerationHandler) Reject(c *gin.Context) {
	id := c.Param("moderationId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	itemType := c.DefaultQuery("type", "")
	if itemType == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.ModerationService.Reject(c.Request.Context(), itemType, id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"rejected": true})
}

func (h *ModerationHandler) Escalate(c *gin.Context) {
	id := c.Param("moderationId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	itemType := c.DefaultQuery("type", "")
	if itemType == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.ModerationService.Escalate(c.Request.Context(), itemType, id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"escalated": true})
}
