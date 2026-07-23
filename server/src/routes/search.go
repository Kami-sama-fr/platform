package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type SearchHandler struct {
	deps Dependencies
}

func NewSearchHandler(deps Dependencies) *SearchHandler {
	return &SearchHandler{deps: deps}
}

func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 50 {
		limit = 20
	}

	results, err := h.deps.SearchService.Search(c.Request.Context(), q, limit)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": results})
}

func (h *SearchHandler) SearchAnime(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 50 {
		limit = 20
	}

	items, err := h.deps.SearchService.SearchAnime(c.Request.Context(), q, limit)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SearchHandler) SearchCharacters(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 50 {
		limit = 20
	}

	items, err := h.deps.SearchService.SearchCharacters(c.Request.Context(), q, limit)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SearchHandler) SearchStudios(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 50 {
		limit = 20
	}

	items, err := h.deps.SearchService.SearchStudios(c.Request.Context(), q, limit)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SearchHandler) Suggestions(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 20 {
		limit = 10
	}

	items, err := h.deps.SearchService.Suggestions(c.Request.Context(), q, limit)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}
