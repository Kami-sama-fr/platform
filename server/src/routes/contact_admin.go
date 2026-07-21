package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type ContactAdminHandler struct {
	deps Dependencies
}

func NewContactAdminHandler(deps Dependencies) *ContactAdminHandler {
	return &ContactAdminHandler{deps: deps}
}

func (h *ContactAdminHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, total, err := h.deps.ContactAdminService.List(c.Request.Context(), c.DefaultQuery("status", ""), page, limit)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *ContactAdminHandler) GetByID(c *gin.Context) {
	id := c.Param("messageId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.ContactAdminService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *ContactAdminHandler) Update(c *gin.Context) {
	id := c.Param("messageId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	var item interface{}
	var err error
	switch req.Status {
	case "read":
		item, err = h.deps.ContactAdminService.MarkAsRead(c.Request.Context(), id)
	case "replied":
		item, err = h.deps.ContactAdminService.MarkAsReplied(c.Request.Context(), id)
	default:
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *ContactAdminHandler) Reply(c *gin.Context) {
	id := c.Param("messageId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.ContactAdminService.Reply(c.Request.Context(), id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"sent": true})
}

func (h *ContactAdminHandler) Delete(c *gin.Context) {
	id := c.Param("messageId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.ContactAdminService.Delete(c.Request.Context(), id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}
