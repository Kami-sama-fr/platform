package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/kami-sama/server/src/middleware"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type SupportHandler struct {
	deps Dependencies
}

func NewSupportHandler(deps Dependencies) *SupportHandler {
	return &SupportHandler{deps: deps}
}

func (h *SupportHandler) ListTickets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, total, err := h.deps.SupportService.ListTickets(
		c.Request.Context(),
		c.DefaultQuery("status", ""),
		c.DefaultQuery("priority", ""),
		c.DefaultQuery("category", ""),
		page, limit,
	)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *SupportHandler) CreateTicket(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req struct {
		Subject     string `json:"subject"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		Category    string `json:"category"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Subject == "" || req.Description == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SupportService.CreateTicket(c.Request.Context(), principal.UserID, req.Subject, req.Description, req.Priority, req.Category)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *SupportHandler) GetTicket(c *gin.Context) {
	id := c.Param("ticketId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SupportService.GetTicket(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SupportHandler) UpdateTicket(c *gin.Context) {
	id := c.Param("ticketId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	var req struct {
		Status     *string `json:"status"`
		AssignedTo *string `json:"assignedTo"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SupportService.UpdateTicket(c.Request.Context(), id, req.Status, req.AssignedTo)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SupportHandler) ReplyToTicket(c *gin.Context) {
	id := c.Param("ticketId")
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
		Content string `json:"content"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Content == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SupportService.ReplyToTicket(c.Request.Context(), id, principal.UserID, req.Content, true)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *SupportHandler) CloseTicket(c *gin.Context) {
	id := c.Param("ticketId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SupportService.CloseTicket(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SupportHandler) EscalateTicket(c *gin.Context) {
	id := c.Param("ticketId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SupportService.EscalateTicket(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SupportHandler) SupportLogs(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"items": []any{}})
}

func (h *SupportHandler) ExportLogs(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=support_logs.csv")
	c.String(http.StatusOK, "timestamp,action,userId,detail\n")
}
