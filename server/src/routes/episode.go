package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/middleware"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type EpisodeHandler struct {
	deps Dependencies
}

func NewEpisodeHandler(deps Dependencies) *EpisodeHandler {
	return &EpisodeHandler{deps: deps}
}

func (h *EpisodeHandler) ListByAnime(c *gin.Context) {
	animeID := c.Param("animeId")
	if animeID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	var seasonID *string
	if s := c.Query("seasonId"); s != "" {
		seasonID = &s
	}
	items, err := h.deps.EpisodeService.ListByAnime(c.Request.Context(), animeID, seasonID)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *EpisodeHandler) GetByID(c *gin.Context) {
	id := c.Param("episodeId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.EpisodeService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *EpisodeHandler) Create(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	animeID := c.Param("animeId")
	if animeID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	var req struct {
		Number       int     `json:"number"`
		Title        string  `json:"title"`
		Synopsis     string  `json:"synopsis"`
		ThumbnailUrl string  `json:"thumbnailUrl"`
		Duration     float64 `json:"duration"`
		AirDate      string  `json:"airDate"`
		IsSubbed     bool    `json:"isSubbed"`
		IsDubbed     bool    `json:"isDubbed"`
		SeasonID     string  `json:"seasonId"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.EpisodeService.Create(c.Request.Context(), principal.UserID, animeID, &req.SeasonID, req.Number, req.Title, req.Synopsis, req.ThumbnailUrl, req.Duration, req.AirDate, req.IsSubbed, req.IsDubbed)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *EpisodeHandler) Update(c *gin.Context) {
	id := c.Param("episodeId")
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
		Number       *int     `json:"number"`
		Title        *string  `json:"title"`
		Synopsis     *string  `json:"synopsis"`
		ThumbnailUrl *string  `json:"thumbnailUrl"`
		Duration     *float64 `json:"duration"`
		AirDate      *string  `json:"airDate"`
		IsSubbed     *bool    `json:"isSubbed"`
		IsDubbed     *bool    `json:"isDubbed"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.EpisodeService.Update(c.Request.Context(), principal.UserID, id, req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *EpisodeHandler) Delete(c *gin.Context) {
	id := c.Param("episodeId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.EpisodeService.Delete(c.Request.Context(), id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

func (h *EpisodeHandler) GetNumber(c *gin.Context) {
	animeID := c.Param("animeId")
	numStr := c.Param("number")
	if animeID == "" || numStr == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	num, err := strconv.Atoi(numStr)
	if err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.EpisodeService.GetByNumber(c.Request.Context(), animeID, num)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}
