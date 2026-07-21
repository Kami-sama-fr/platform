package routes

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	redisclient "github.com/kami-sama-fr/platform/server/internal/redis"
	"github.com/kami-sama-fr/platform/server/src/config"
	"github.com/kami-sama-fr/platform/server/src/interfaces"
	"github.com/kami-sama-fr/platform/server/src/middleware"
	"github.com/kami-sama-fr/platform/server/src/services"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type Dependencies struct {
	Config              config.Config
	Logger              *slog.Logger
	RuntimeRole         string
	Database            interfaces.Database
	Redis               *redisclient.Client
	EventBus            interfaces.EventBus
	IdentityProvider    interfaces.IdentityProvider
	AuthService         *services.AuthService
	OAuthService        *services.OAuthService
	UserService         *services.UserService
	WorkspaceService    *services.WorkspaceService
	Repos               *services.Repositories
	AnimeService        *services.AnimeService
	EpisodeService      *services.EpisodeService
	GenreService        *services.GenreService
	StudioService       *services.StudioService
	CharacterService    *services.CharacterService
	MediaService        *services.MediaService
	CommunityService    *services.CommunityService
	WatchService        *services.WatchService
	SchedulingService   *services.SchedulingService
	NotificationService *services.NotificationService
	SearchService       *services.SearchService
	SettingsService     *services.SettingsService
}

func SetupRoutes(router *gin.Engine, deps Dependencies) {
	handler := &apiHandler{deps: deps}
	platform := NewPlatformHandler(deps)
	anime := NewAnimeHandler(deps)
	episode := NewEpisodeHandler(deps)
	genre := NewGenreHandler(deps)
	studio := NewStudioHandler(deps)
	character := NewCharacterHandler(deps)
	media := NewMediaHandler(deps)
	community := NewCommunityHandler(deps)
	watch := NewWatchHandler(deps)
	scheduling := NewSchedulingHandler(deps)
	notification := NewNotificationHandler(deps)
	search := NewSearchHandler(deps)
	settings := NewSettingsHandler(deps)

	router.GET("/health/live", handler.live)
	router.GET("/health/ready", handler.ready)
	router.GET("/metrics", handler.metrics)

	api := router.Group("/api/v1")
	api.GET("/health", handler.health)
	api.GET("/ready", handler.ready)
	api.POST("/webhooks/:provider/:integrationId", handler.webhook)

	auth := api.Group("/auth")
	{
		auth.POST("/register", handler.register)
		auth.POST("/login", handler.login)
		auth.POST("/refresh", handler.refresh)
		auth.POST("/logout", handler.logout)
		auth.POST("/forgot-password", handler.forgotPassword)
		auth.POST("/reset-password", handler.resetPassword)
		auth.POST("/verify-email", handler.verifyEmail)
		auth.POST("/resend-verification", handler.resendVerification)
		auth.GET("/oauth/:provider", handler.oauthLogin)
		auth.GET("/oauth/:provider/callback", handler.oauthCallback)
		authProtected := auth.Group("")
		authProtected.Use(middleware.Auth(deps.IdentityProvider))
		{
			authProtected.POST("/logout-all", handler.logoutAll)
			authProtected.POST("/change-password", handler.changePassword)
			authProtected.GET("/me", handler.authMe)
			authProtected.GET("/sessions", handler.listAuthSessions)
			authProtected.DELETE("/sessions/:sessionId", handler.deleteAuthSession)
			authProtected.GET("/oauth/accounts", handler.listOAuthAccounts)
			authProtected.DELETE("/oauth/:provider", handler.unlinkOAuthAccount)
		}
	}

	protected := api.Group("")
	protected.Use(middleware.Auth(deps.IdentityProvider))
	protected.Use(middleware.WorkspaceContext())
	{
		protected.GET("/me", handler.me)
		protected.PATCH("/me", handler.updateMe)

		protected.GET("/workspaces", handler.listWorkspaces)
		protected.POST("/workspaces", handler.createWorkspace)
		protected.GET("/workspaces/:workspaceId", handler.getWorkspace)
		protected.PATCH("/workspaces/:workspaceId", handler.updateWorkspace)
		protected.DELETE("/workspaces/:workspaceId", handler.deleteWorkspace)
		protected.GET("/workspaces/:workspaceId/members", handler.listWorkspaceMembers)
		protected.POST("/workspaces/:workspaceId/members", handler.createWorkspaceMember)
		protected.POST("/workspaces/:workspaceId/members/provision", handler.provisionWorkspaceUser)
		protected.PATCH("/workspaces/:workspaceId/members/:userId", handler.updateWorkspaceMember)
		protected.DELETE("/workspaces/:workspaceId/members/:userId", handler.deleteWorkspaceMember)

		platformGroup := protected.Group("/platform")
		{
			platformGroup.GET("/home", platform.GetHomeData)
			platformGroup.GET("/wallet", platform.GetWalletInfo)
			platformGroup.GET("/user-info", platform.GetUserInfo)
			platformGroup.PATCH("/user-info", platform.UpdateUserInfo)
			platformGroup.GET("/security", platform.GetSecurityInfo)
			platformGroup.POST("/password/change", platform.ChangePassword)
			platformGroup.GET("/applications", platform.ListApplications)
			platformGroup.GET("/applications/:applicationId", platform.GetApplicationDetails)
			platformGroup.GET("/data-privacy", platform.GetDataPrivacyInfo)
			platformGroup.PATCH("/data-privacy", platform.UpdatePrivacySettings)
			platformGroup.GET("/contacts", platform.ListContacts)
			platformGroup.GET("/contacts/:contactId", platform.GetContactDetails)
			platformGroup.GET("/family", platform.ListFamilyMembers)
			platformGroup.POST("/family/invite", platform.InviteFamilyMember)
			platformGroup.GET("/storage", platform.GetStorageInfo)
			platformGroup.GET("/storage/files", platform.ListFiles)
			platformGroup.GET("/settings", platform.GetSettings)
			platformGroup.PATCH("/settings", platform.UpdateSettings)
		}

		animeGroup := protected.Group("/anime")
		{
			animeGroup.GET("", anime.List)
			animeGroup.POST("", anime.Create)
			animeGroup.GET("/:animeId", anime.GetByID)
			animeGroup.GET("/slug/:slug", anime.GetBySlug)
			animeGroup.PATCH("/:animeId", anime.Update)
			animeGroup.DELETE("/:animeId", anime.Delete)
		}

		genreGroup := protected.Group("/genres")
		{
			genreGroup.GET("", genre.List)
			genreGroup.POST("", genre.Create)
			genreGroup.GET("/:genreId", genre.GetByID)
			genreGroup.GET("/slug/:slug", genre.GetBySlug)
			genreGroup.PATCH("/:genreId", genre.Update)
			genreGroup.DELETE("/:genreId", genre.Delete)
		}

		studioGroup := protected.Group("/studios")
		{
			studioGroup.GET("", studio.List)
			studioGroup.POST("", studio.Create)
			studioGroup.GET("/:studioId", studio.GetByID)
			studioGroup.GET("/slug/:slug", studio.GetBySlug)
			studioGroup.PATCH("/:studioId", studio.Update)
			studioGroup.DELETE("/:studioId", studio.Delete)
		}

		characterGroup := protected.Group("/characters")
		{
			characterGroup.GET("", character.List)
			characterGroup.POST("", character.Create)
			characterGroup.GET("/:characterId", character.GetByID)
			characterGroup.GET("/slug/:slug", character.GetBySlug)
			characterGroup.PATCH("/:characterId", character.Update)
			characterGroup.DELETE("/:characterId", character.Delete)
		}

		episodeGroup := protected.Group("/anime/:animeId/episodes")
		{
			episodeGroup.GET("", episode.ListByAnime)
			episodeGroup.POST("", episode.Create)
			episodeGroup.GET("/:episodeId", episode.GetByID)
			episodeGroup.GET("/number/:episodeNumber", episode.GetNumber)
			episodeGroup.PATCH("/:episodeId", episode.Update)
			episodeGroup.DELETE("/:episodeId", episode.Delete)
		}

		mediaGroup := protected.Group("/media")
		{
			mediaGroup.GET("", media.List)
			mediaGroup.POST("", media.Create)
			mediaGroup.GET("/:mediaId", media.GetByID)
			mediaGroup.PATCH("/:mediaId", media.Update)
			mediaGroup.DELETE("/:mediaId", media.Delete)
			mediaGroup.GET("/encoding-jobs", media.ListEncodingJobs)
			mediaGroup.GET("/encoding-jobs/:jobId", media.GetEncodingJob)
		}

		communityGroup := protected.Group("/community")
		{
			communityGroup.GET("/reviews", community.ListReviews)
			communityGroup.POST("/reviews", community.CreateReview)
			communityGroup.GET("/reviews/:reviewId", community.GetReview)
			communityGroup.PATCH("/reviews/:reviewId", community.UpdateReview)
			communityGroup.DELETE("/reviews/:reviewId", community.DeleteReview)
			communityGroup.GET("/reviews/:reviewId/comments", community.ListCommentsByReview)
			communityGroup.POST("/reviews/:reviewId/comments", community.CreateComment)
			communityGroup.PATCH("/comments/:commentId", community.UpdateComment)
			communityGroup.DELETE("/comments/:commentId", community.DeleteComment)
			communityGroup.GET("/watchlists", community.ListWatchlists)
			communityGroup.POST("/watchlists", community.CreateWatchlist)
			communityGroup.GET("/watchlists/:watchlistId", community.GetWatchlist)
			communityGroup.PATCH("/watchlists/:watchlistId", community.UpdateWatchlist)
			communityGroup.DELETE("/watchlists/:watchlistId", community.DeleteWatchlist)
			communityGroup.GET("/watchlists/:watchlistId/anime", community.ListWatchlistAnime)
			communityGroup.POST("/watchlists/:watchlistId/anime", community.AddToWatchlist)
			communityGroup.DELETE("/watchlists/:watchlistId/anime/:animeId", community.RemoveFromWatchlist)
			communityGroup.GET("/reports", community.ListReports)
			communityGroup.POST("/reports", community.CreateReport)
			communityGroup.PATCH("/reports/:reportId", community.UpdateReport)
		}

		watchGroup := protected.Group("/watch")
		{
			watchGroup.GET("/progress/:episodeId", watch.GetProgress)
			watchGroup.PUT("/progress/:episodeId", watch.UpsertProgress)
			watchGroup.GET("/progress", watch.ListProgress)
			watchGroup.GET("/continue", watch.ContinueWatching)
			watchGroup.GET("/history", watch.ListHistory)
			watchGroup.POST("/history", watch.AddHistory)
		}

		schedulingGroup := protected.Group("/scheduling")
		{
			schedulingGroup.GET("/simulcasts", scheduling.ListSimulcasts)
			schedulingGroup.POST("/simulcasts", scheduling.CreateSimulcast)
			schedulingGroup.GET("/simulcasts/:simulcastId", scheduling.GetSimulcast)
			schedulingGroup.PATCH("/simulcasts/:simulcastId", scheduling.UpdateSimulcast)
			schedulingGroup.DELETE("/simulcasts/:simulcastId", scheduling.DeleteSimulcast)
			schedulingGroup.GET("/simulcasts/week/:weekday", scheduling.GetSimulcastByWeek)
			schedulingGroup.GET("/upcoming", scheduling.ListUpcomingReleases)
			schedulingGroup.POST("/schedules", scheduling.CreateReleaseSchedule)
			schedulingGroup.PATCH("/schedules/:scheduleId", scheduling.UpdateReleaseSchedule)
		}

		notificationGroup := protected.Group("/notifications")
		{
			notificationGroup.GET("", notification.List)
			notificationGroup.GET("/:notificationId", notification.GetByID)
			notificationGroup.PATCH("/:notificationId/read", notification.MarkRead)
			notificationGroup.POST("/read-all", notification.MarkAllRead)
			notificationGroup.DELETE("/:notificationId", notification.Delete)
			notificationGroup.GET("/unread-count", notification.UnreadCount)
			notificationGroup.GET("/preferences", notification.GetPreferences)
			notificationGroup.PUT("/preferences", notification.UpdatePreferences)
		}

		searchGroup := protected.Group("/search")
		{
			searchGroup.GET("", search.Search)
			searchGroup.GET("/anime", search.SearchAnime)
			searchGroup.GET("/characters", search.SearchCharacters)
			searchGroup.GET("/studios", search.SearchStudios)
			searchGroup.GET("/suggestions", search.Suggestions)
		}

		settingsGroup := protected.Group("/settings")
		{
			settingsGroup.GET("", settings.List)
			settingsGroup.GET("/seo", settings.GetSeoMeta)
			settingsGroup.PUT("/seo", settings.UpsertSeoMeta)
			settingsGroup.GET("/:key", settings.GetByKey)
			settingsGroup.PUT("/:key", settings.Upsert)
			settingsGroup.DELETE("/:key", settings.Delete)
		}
	}
}

type apiHandler struct {
	deps Dependencies
}

func (h *apiHandler) principal(c *gin.Context) (interfaces.Principal, bool) {
	return middleware.PrincipalFromGin(c)
}

func (h *apiHandler) live(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{
		"status":  "alive",
		"role":    h.runtimeRole(),
		"version": h.deps.Config.App.Version,
	})
}

func (h *apiHandler) runtimeRole() string {
	if h.deps.RuntimeRole == "" {
		return "api"
	}
	return h.deps.RuntimeRole
}

func (h *apiHandler) health(c *gin.Context) {
	status := "healthy"
	redisStatus := "disabled"
	if err := h.deps.Database.Ping(c.Request.Context()); err != nil {
		status = "degraded"
	}
	if h.deps.Redis != nil {
		redisHealth := h.deps.Redis.Health(c.Request.Context())
		redisStatus = string(redisHealth.Status)
		if redisHealth.Status != redisclient.StatusHealthy {
			status = "degraded"
		}
	}
	utils.Success(c, http.StatusOK, gin.H{
		"status":   status,
		"database": "healthy",
		"redis":    redisStatus,
		"role":     h.runtimeRole(),
		"version":  h.deps.Config.App.Version,
	})
}

func (h *apiHandler) ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	result := gin.H{"database": "healthy", "redis": "disabled"}
	if err := h.deps.Database.Ping(ctx); err != nil {
		result["database"] = "unhealthy"
		utils.Error(c, utils.ErrDependencyUnavailable)
		return
	}
	if h.deps.Config.Redis.Enabled {
		if h.deps.Redis == nil || !h.deps.Redis.IsAvailable() {
			if h.deps.Config.Redis.Required {
				result["redis"] = "unhealthy"
				utils.Error(c, utils.ErrDependencyUnavailable)
				return
			}
			result["redis"] = "unavailable"
		} else {
			result["redis"] = "healthy"
		}
	}
	utils.Success(c, http.StatusOK, gin.H{
		"status":   "ready",
		"database": result["database"],
		"redis":    result["redis"],
		"role":     h.runtimeRole(),
		"version":  h.deps.Config.App.Version,
	})
}

func (h *apiHandler) me(c *gin.Context) {
	principal, _ := h.principal(c)
	user, err := h.deps.UserService.GetMe(c.Request.Context(), principal)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{
		"id":                user.ID,
		"email":             user.Email,
		"displayName":       user.DisplayName,
		"avatarUrl":         user.AvatarURL,
		"status":            user.Status,
		"presenceStatus":    user.PresenceStatus,
		"roles":             principal.Roles,
		"permissions":       principal.Permissions,
		"workspaceId":       principal.WorkspaceID,
		"createdAt":         user.CreatedAt,
		"updatedAt":         user.UpdatedAt,
		"lastSeenAt":        user.LastSeenAt,
		"disabledAt":        user.DisabledAt,
		"emailVerifiedAt":   user.EmailVerifiedAt,
		"passwordChangedAt": user.PasswordChangedAt,
	})
}

func (h *apiHandler) updateMe(c *gin.Context) {
	var req struct {
		DisplayName string `json:"displayName"`
		AvatarURL   string `json:"avatarUrl"`
		Status      string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	principal, _ := h.principal(c)
	user, err := h.deps.UserService.UpdateMe(c.Request.Context(), principal, req.DisplayName, req.AvatarURL, req.Status)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{
		"id":                user.ID,
		"email":             user.Email,
		"displayName":       user.DisplayName,
		"avatarUrl":         user.AvatarURL,
		"status":            user.Status,
		"presenceStatus":    user.PresenceStatus,
		"roles":             principal.Roles,
		"permissions":       principal.Permissions,
		"workspaceId":       principal.WorkspaceID,
		"createdAt":         user.CreatedAt,
		"updatedAt":         user.UpdatedAt,
		"lastSeenAt":        user.LastSeenAt,
		"disabledAt":        user.DisabledAt,
		"emailVerifiedAt":   user.EmailVerifiedAt,
		"passwordChangedAt": user.PasswordChangedAt,
	})
}

func (h *apiHandler) listWorkspaces(c *gin.Context) {
	principal, _ := h.principal(c)
	items, err := h.deps.WorkspaceService.List(c.Request.Context(), principal)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.List(c, items, "", false)
}

func (h *apiHandler) createWorkspace(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	if c.ShouldBindJSON(&req) != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	principal, _ := h.principal(c)
	item, err := h.deps.WorkspaceService.Create(c.Request.Context(), principal, req.Name, req.Slug, req.Description)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *apiHandler) getWorkspace(c *gin.Context)    { h.workspaceResource(c, "get") }
func (h *apiHandler) updateWorkspace(c *gin.Context) { h.workspaceResource(c, "update") }
func (h *apiHandler) deleteWorkspace(c *gin.Context) { h.workspaceResource(c, "delete") }

func (h *apiHandler) workspaceResource(c *gin.Context, action string) {
	principal, _ := h.principal(c)
	id := c.Param("workspaceId")
	switch action {
	case "get":
		item, err := h.deps.WorkspaceService.Get(c.Request.Context(), principal, id)
		if err != nil {
			utils.Error(c, err)
			return
		}
		utils.Success(c, http.StatusOK, item)
	case "update":
		var req struct{ Name, Description string }
		if c.ShouldBindJSON(&req) != nil {
			utils.Error(c, utils.ErrValidationFailed)
			return
		}
		item, err := h.deps.WorkspaceService.Update(c.Request.Context(), principal, id, req.Name, req.Description)
		if err != nil {
			utils.Error(c, err)
			return
		}
		utils.Success(c, http.StatusOK, item)
	case "delete":
		if err := h.deps.WorkspaceService.Archive(c.Request.Context(), principal, id); err != nil {
			utils.Error(c, err)
			return
		}
		utils.Success(c, http.StatusOK, gin.H{"deleted": true})
	}
}

func (h *apiHandler) listWorkspaceMembers(c *gin.Context)   { h.membersResource(c, "list") }
func (h *apiHandler) createWorkspaceMember(c *gin.Context)  { h.membersResource(c, "create") }
func (h *apiHandler) provisionWorkspaceUser(c *gin.Context) { h.membersResource(c, "provision") }
func (h *apiHandler) updateWorkspaceMember(c *gin.Context)  { h.membersResource(c, "update") }
func (h *apiHandler) deleteWorkspaceMember(c *gin.Context)  { h.membersResource(c, "delete") }

func (h *apiHandler) membersResource(c *gin.Context, action string) {
	principal, _ := h.principal(c)
	workspaceID := c.Param("workspaceId")
	switch action {
	case "list":
		items, err := h.deps.WorkspaceService.ListMembers(c.Request.Context(), principal, workspaceID)
		if err != nil {
			utils.Error(c, err)
			return
		}
		utils.List(c, items, "", false)
	case "create":
		var req struct {
			UserID string `json:"userId"`
			Email  string `json:"email"`
			Role   string `json:"role"`
		}
		if c.ShouldBindJSON(&req) != nil {
			utils.Error(c, utils.ErrValidationFailed)
			return
		}
		item, err := h.deps.WorkspaceService.AddMember(c.Request.Context(), principal, workspaceID, req.UserID, req.Email, req.Role)
		if err != nil {
			utils.Error(c, err)
			return
		}
		utils.Success(c, http.StatusCreated, item)
	case "provision":
		var req struct {
			Email             string `json:"email"`
			DisplayName       string `json:"displayName"`
			Role              string `json:"role"`
			TemporaryPassword string `json:"temporaryPassword"`
		}
		if c.ShouldBindJSON(&req) != nil {
			utils.Error(c, utils.ErrValidationFailed)
			return
		}
		item, err := h.deps.WorkspaceService.ProvisionWorkspaceUser(
			c.Request.Context(),
			principal,
			workspaceID,
			services.ProvisionWorkspaceUserInput{
				Email:             req.Email,
				DisplayName:       req.DisplayName,
				Role:              req.Role,
				TemporaryPassword: req.TemporaryPassword,
			},
		)
		if err != nil {
			utils.Error(c, err)
			return
		}
		utils.Success(c, http.StatusCreated, item)
	case "update":
		var req struct{ Role string }
		if c.ShouldBindJSON(&req) != nil {
			utils.Error(c, utils.ErrValidationFailed)
			return
		}
		item, err := h.deps.WorkspaceService.UpdateMember(c.Request.Context(), principal, workspaceID, c.Param("userId"), req.Role)
		if err != nil {
			utils.Error(c, err)
			return
		}
		utils.Success(c, http.StatusOK, item)
	case "delete":
		if err := h.deps.WorkspaceService.RemoveMember(c.Request.Context(), principal, workspaceID, c.Param("userId")); err != nil {
			utils.Error(c, err)
			return
		}
		utils.Success(c, http.StatusOK, gin.H{"deleted": true})
	}
}

func (h *apiHandler) metrics(c *gin.Context) {
	c.String(http.StatusOK, "")
}

func (h *apiHandler) webhook(c *gin.Context) {
	payload, _ := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	_ = payload
	utils.Success(c, http.StatusAccepted, gin.H{"accepted": true})
}
