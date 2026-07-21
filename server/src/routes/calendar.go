package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/middleware"
	"github.com/kami-sama-fr/platform/server/src/services"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type CalendarHandler struct {
	deps Dependencies
}

func NewCalendarHandler(deps Dependencies) *CalendarHandler {
	return &CalendarHandler{deps: deps}
}

func (h *CalendarHandler) ListEvents(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	var items any
	var err error
	if startDate != "" || endDate != "" {
		items, err = h.deps.CalendarService.ListEvents(c.Request.Context(), startDate, endDate)
	} else {
		items, err = h.deps.CalendarService.ListEventsAll(c.Request.Context())
	}
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *CalendarHandler) CreateEvent(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		EventType   string  `json:"eventType"`
		AnimeID     *string `json:"animeId"`
		EpisodeID   *string `json:"episodeId"`
		StartAt     string  `json:"startAt"`
		EndAt       string  `json:"endAt"`
		AllDay      bool    `json:"allDay"`
		Color       string  `json:"color"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CalendarService.CreateEvent(c.Request.Context(), principal.UserID, req.Title, req.Description, req.EventType, req.AnimeID, req.EpisodeID, req.StartAt, req.EndAt, req.AllDay, req.Color)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *CalendarHandler) GetEvent(c *gin.Context) {
	id := c.Param("eventId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CalendarService.GetEvent(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CalendarHandler) UpdateEvent(c *gin.Context) {
	id := c.Param("eventId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req services.UpdateCalendarEventInput
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CalendarService.UpdateEvent(c.Request.Context(), principal.UserID, id, req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CalendarHandler) DeleteEvent(c *gin.Context) {
	id := c.Param("eventId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.CalendarService.DeleteEvent(c.Request.Context(), id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}
