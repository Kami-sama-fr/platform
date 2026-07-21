package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/middleware"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type CommunityHandler struct {
	deps Dependencies
}

func NewCommunityHandler(deps Dependencies) *CommunityHandler {
	return &CommunityHandler{deps: deps}
}

// --- Reviews ---

func (h *CommunityHandler) ListReviews(c *gin.Context) {
	animeID := c.Param("animeId")
	if animeID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.deps.CommunityService.ListReviews(c.Request.Context(), animeID, limit, offset)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *CommunityHandler) GetReview(c *gin.Context) {
	id := c.Param("reviewId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CommunityService.GetReview(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CommunityHandler) CreateReview(c *gin.Context) {
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
		Rating  int    `json:"rating"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Spoiler bool   `json:"spoiler"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CommunityService.CreateReview(c.Request.Context(), principal.UserID, animeID, req.Rating, req.Title, req.Content, req.Spoiler)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *CommunityHandler) UpdateReview(c *gin.Context) {
	id := c.Param("reviewId")
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
		Rating  *int    `json:"rating"`
		Title   *string `json:"title"`
		Content *string `json:"content"`
		Spoiler *bool   `json:"spoiler"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CommunityService.UpdateReview(c.Request.Context(), principal.UserID, id, req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CommunityHandler) DeleteReview(c *gin.Context) {
	id := c.Param("reviewId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	if err := h.deps.CommunityService.DeleteReview(c.Request.Context(), principal.UserID, id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

// --- Comments ---

func (h *CommunityHandler) ListComments(c *gin.Context) {
	animeID := c.Param("animeId")
	if animeID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.deps.CommunityService.ListComments(c.Request.Context(), animeID, limit, offset)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *CommunityHandler) ListCommentsByReview(c *gin.Context) {
	reviewID := c.Param("reviewId")
	if reviewID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.deps.CommunityService.ListCommentsByReview(c.Request.Context(), reviewID, limit, offset)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *CommunityHandler) CreateComment(c *gin.Context) {
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
		ReviewID string `json:"reviewId"`
		ParentID string `json:"parentId"`
		Content  string `json:"content"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Content == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CommunityService.CreateComment(c.Request.Context(), principal.UserID, animeID, &req.ReviewID, &req.ParentID, req.Content)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *CommunityHandler) UpdateComment(c *gin.Context) {
	id := c.Param("commentId")
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
	item, err := h.deps.CommunityService.UpdateComment(c.Request.Context(), principal.UserID, id, req.Content)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CommunityHandler) DeleteComment(c *gin.Context) {
	id := c.Param("commentId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	if err := h.deps.CommunityService.DeleteComment(c.Request.Context(), principal.UserID, id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

// --- Reports ---

func (h *CommunityHandler) CreateReport(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req struct {
		TargetType  string `json:"targetType"`
		TargetID    string `json:"targetId"`
		Reason      string `json:"reason"`
		Description string `json:"description"`
	}
	if c.ShouldBindJSON(&req) != nil || req.TargetType == "" || req.TargetID == "" || req.Reason == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CommunityService.CreateReport(c.Request.Context(), principal.UserID, req.TargetType, req.TargetID, req.Reason, req.Description)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *CommunityHandler) ListReports(c *gin.Context) {
	status := c.DefaultQuery("status", "")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.deps.CommunityService.ListReports(c.Request.Context(), status, limit, offset)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *CommunityHandler) UpdateReport(c *gin.Context) {
	id := c.Param("reportId")
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
		Status string `json:"status"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CommunityService.UpdateReport(c.Request.Context(), principal.UserID, id, req.Status)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

// --- Watchlists ---

func (h *CommunityHandler) ListWatchlists(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	items, err := h.deps.CommunityService.ListWatchlists(c.Request.Context(), principal.UserID)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *CommunityHandler) GetWatchlist(c *gin.Context) {
	id := c.Param("watchlistId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CommunityService.GetWatchlist(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CommunityHandler) CreateWatchlist(c *gin.Context) {
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var req struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Name == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if req.Type == "" {
		req.Type = "custom"
	}
	item, err := h.deps.CommunityService.CreateWatchlist(c.Request.Context(), principal.UserID, req.Name, req.Type)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *CommunityHandler) UpdateWatchlist(c *gin.Context) {
	id := c.Param("watchlistId")
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
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.CommunityService.UpdateWatchlist(c.Request.Context(), principal.UserID, id, req.Name, req.Type)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CommunityHandler) DeleteWatchlist(c *gin.Context) {
	id := c.Param("watchlistId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	if err := h.deps.CommunityService.DeleteWatchlist(c.Request.Context(), principal.UserID, id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

func (h *CommunityHandler) AddToWatchlist(c *gin.Context) {
	id := c.Param("watchlistId")
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
		AnimeID string `json:"animeId"`
	}
	if c.ShouldBindJSON(&req) != nil || req.AnimeID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.CommunityService.AddToWatchlist(c.Request.Context(), principal.UserID, id, req.AnimeID); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, gin.H{"added": true})
}

func (h *CommunityHandler) RemoveFromWatchlist(c *gin.Context) {
	id := c.Param("watchlistId")
	animeID := c.Param("animeId")
	if id == "" || animeID == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	if err := h.deps.CommunityService.RemoveFromWatchlist(c.Request.Context(), principal.UserID, id, animeID); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"removed": true})
}

func (h *CommunityHandler) ListWatchlistAnime(c *gin.Context) {
	id := c.Param("watchlistId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	items, err := h.deps.CommunityService.ListWatchlistAnime(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

// --- Admin Comments ---

func (h *CommunityHandler) AdminListComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, total, err := h.deps.Repos.Comments().ListAll(
		c.Request.Context(),
		c.DefaultQuery("status", ""),
		c.DefaultQuery("author", ""),
		c.DefaultQuery("anime", ""),
		limit,
		(page-1)*limit,
	)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *CommunityHandler) AdminGetComment(c *gin.Context) {
	id := c.Param("commentId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.Repos.Comments().GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CommunityHandler) AdminModerateComment(c *gin.Context) {
	id := c.Param("commentId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	comment, err := h.deps.Repos.Comments().GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	comment.UpdatedAt = time.Now()
	if err := h.deps.Repos.Comments().Update(c.Request.Context(), comment); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, comment)
}

func (h *CommunityHandler) AdminDeleteComment(c *gin.Context) {
	id := c.Param("commentId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	if err := h.deps.Repos.Comments().Delete(c.Request.Context(), id); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

func (h *CommunityHandler) AdminApproveComment(c *gin.Context) {
	id := c.Param("commentId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	comment, err := h.deps.Repos.Comments().GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	comment.UpdatedAt = time.Now()
	if err := h.deps.Repos.Comments().Update(c.Request.Context(), comment); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, comment)
}

func (h *CommunityHandler) AdminFlagComment(c *gin.Context) {
	id := c.Param("commentId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	comment, err := h.deps.Repos.Comments().GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	comment.UpdatedAt = time.Now()
	if err := h.deps.Repos.Comments().Update(c.Request.Context(), comment); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, comment)
}

// --- Admin Reviews ---

func (h *CommunityHandler) AdminListReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	rating, _ := strconv.Atoi(c.DefaultQuery("rating", "0"))
	items, total, err := h.deps.Repos.Reviews().ListAll(
		c.Request.Context(),
		rating,
		c.DefaultQuery("author", ""),
		c.DefaultQuery("anime", ""),
		limit,
		(page-1)*limit,
	)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *CommunityHandler) AdminFeatureReview(c *gin.Context) {
	id := c.Param("reviewId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	review, err := h.deps.Repos.Reviews().GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	review.UpdatedAt = time.Now()
	if err := h.deps.Repos.Reviews().Update(c.Request.Context(), review); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"featured": true, "review": review})
}

// --- Admin Reports ---

func (h *CommunityHandler) AdminListReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, total, err := h.deps.Repos.Reports().ListAll(
		c.Request.Context(),
		c.DefaultQuery("status", ""),
		c.DefaultQuery("type", ""),
		limit,
		(page-1)*limit,
	)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items, "total": total})
}

func (h *CommunityHandler) AdminGetReport(c *gin.Context) {
	id := c.Param("reportId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.Repos.Reports().GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *CommunityHandler) AdminProcessReport(c *gin.Context) {
	id := c.Param("reportId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Status == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	report, err := h.deps.Repos.Reports().GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	report.Status = req.Status
	report.UpdatedAt = time.Now()
	if err := h.deps.Repos.Reports().Update(c.Request.Context(), report); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, report)
}

func (h *CommunityHandler) AdminResolveReport(c *gin.Context) {
	id := c.Param("reportId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	report, err := h.deps.Repos.Reports().GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	now := time.Now()
	report.Status = "resolved"
	report.ResolvedAt = &now
	report.UpdatedAt = now
	if err := h.deps.Repos.Reports().Update(c.Request.Context(), report); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, report)
}

func (h *CommunityHandler) AdminDismissReport(c *gin.Context) {
	id := c.Param("reportId")
	if id == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	report, err := h.deps.Repos.Reports().GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Error(c, err)
		return
	}
	report.Status = "dismissed"
	report.UpdatedAt = time.Now()
	if err := h.deps.Repos.Reports().Update(c.Request.Context(), report); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, report)
}

// --- Admin Watchlist Stats ---

func (h *CommunityHandler) AdminWatchlistStats(c *gin.Context) {
	var totalWatchlists int64
	var totalItems int64
	items, err := h.deps.CommunityService.AdminListAllWatchlists(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	totalWatchlists = int64(len(items))
	for _, wl := range items {
		animeItems, err := h.deps.CommunityService.ListWatchlistAnime(c.Request.Context(), wl.ID)
		if err == nil {
			totalItems += int64(len(animeItems))
		}
	}
	utils.Success(c, http.StatusOK, gin.H{
		"totalWatchlists": totalWatchlists,
		"totalItems":      totalItems,
	})
}
