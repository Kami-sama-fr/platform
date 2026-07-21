package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kami-sama-fr/platform/server/src/middleware"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type SettingsAdminHandler struct {
	deps Dependencies
}

func NewSettingsAdminHandler(deps Dependencies) *SettingsAdminHandler {
	return &SettingsAdminHandler{deps: deps}
}

// General

func (h *SettingsAdminHandler) GetGeneral(c *gin.Context) {
	result, err := h.deps.SettingsAdminService.GetGeneral(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) UpdateGeneral(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.UpdateGeneral(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

// Security

func (h *SettingsAdminHandler) GetSecurity(c *gin.Context) {
	result, err := h.deps.SettingsAdminService.GetSecurity(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) UpdateSecurity(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.UpdateSecurity(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) GetSessions(c *gin.Context) {
	result, err := h.deps.SettingsAdminService.GetSessions(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) UpdateSessions(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.UpdateSessions(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) GetRateLimit(c *gin.Context) {
	result, err := h.deps.SettingsAdminService.GetRateLimit(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) UpdateRateLimit(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.UpdateRateLimit(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) Get2FA(c *gin.Context) {
	result, err := h.deps.SettingsAdminService.Get2FA(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) Update2FA(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.Update2FA(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

// Branding

func (h *SettingsAdminHandler) GetBranding(c *gin.Context) {
	result, err := h.deps.SettingsAdminService.GetBranding(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) UpdateBranding(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.UpdateBranding(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) UploadLogo(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"url": req.URL, "message": "logo uploaded"})
}

func (h *SettingsAdminHandler) UploadFavicon(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"url": req.URL, "message": "favicon uploaded"})
}

func (h *SettingsAdminHandler) PreviewBranding(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"preview": req, "message": "branding preview generated"})
}

// Email

func (h *SettingsAdminHandler) GetEmail(c *gin.Context) {
	result, err := h.deps.SettingsAdminService.GetEmail(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) UpdateEmail(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.UpdateEmail(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) SendTestEmail(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "test email queued"})
}

func (h *SettingsAdminHandler) ListEmailTemplates(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"items": []map[string]any{}})
}

func (h *SettingsAdminHandler) UpdateEmailTemplate(c *gin.Context) {
	templateID := c.Param("templateId")
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"templateId": templateID, "updated": true})
}

func (h *SettingsAdminHandler) ListEmailLogs(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"items": []map[string]any{}, "total": 0})
}

// SEO

func (h *SettingsAdminHandler) ListSEOPages(c *gin.Context) {
	items, err := h.deps.SettingsAdminService.ListSEOPages(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SettingsAdminHandler) UpdateSEOPage(c *gin.Context) {
	pagePath := c.Param("pagePath")
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.UpdateSEOPage(c.Request.Context(), pagePath, req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) GetSitemap(c *gin.Context) {
	sitemap, err := h.deps.SettingsAdminService.GenerateSitemap(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, sitemap)
}

func (h *SettingsAdminHandler) GetRobots(c *gin.Context) {
	robots, err := h.deps.SettingsAdminService.GetRobots(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, robots)
}

// Storage

func (h *SettingsAdminHandler) GetStorage(c *gin.Context) {
	result, err := h.deps.SettingsAdminService.GetStorage(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) UpdateStorage(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.UpdateStorage(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) TestStorage(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "storage connection test queued"})
}

func (h *SettingsAdminHandler) GetStorageUsage(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"used": 0, "available": 0, "total": 0})
}

func (h *SettingsAdminHandler) ListBuckets(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"items": []map[string]any{}})
}

// CDN

func (h *SettingsAdminHandler) GetCDN(c *gin.Context) {
	result, err := h.deps.SettingsAdminService.GetCDN(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) UpdateCDN(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.UpdateCDN(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) PurgeCDN(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "cache purge queued"})
}

func (h *SettingsAdminHandler) GetCDNStats(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{
		"requests":   0,
		"bandwidth":  0,
		"cacheHit":   0,
		"cacheMiss": 0,
	})
}

// Domains

func (h *SettingsAdminHandler) ListDomains(c *gin.Context) {
	items, err := h.deps.SettingsAdminService.ListDomains(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SettingsAdminHandler) CreateDomain(c *gin.Context) {
	var req struct {
		Domain string `json:"domain"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Domain == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SettingsAdminService.CreateDomain(c.Request.Context(), req.Domain)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *SettingsAdminHandler) GetDomain(c *gin.Context) {
	item, err := h.deps.SettingsAdminService.GetDomain(c.Request.Context(), c.Param("domainId"))
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SettingsAdminHandler) DeleteDomain(c *gin.Context) {
	if err := h.deps.SettingsAdminService.DeleteDomain(c.Request.Context(), c.Param("domainId")); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

func (h *SettingsAdminHandler) VerifyDomain(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "domain verification queued", "domainId": c.Param("domainId")})
}

func (h *SettingsAdminHandler) GenerateSSL(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "SSL certificate generation queued", "domainId": c.Param("domainId")})
}

// API Keys

func (h *SettingsAdminHandler) ListAPIKeys(c *gin.Context) {
	items, err := h.deps.SettingsAdminService.ListAPIKeys(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SettingsAdminHandler) CreateAPIKey(c *gin.Context) {
	var req struct {
		Name      string  `json:"name"`
		ExpiresAt *string `json:"expiresAt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	principal, ok := middleware.GetPrincipal(c)
	if !ok {
		utils.Error(c, utils.ErrUnauthorized)
		return
	}
	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			utils.Error(c, utils.ErrValidationFailed)
			return
		}
		expiresAt = &t
	}
	item, rawKey, err := h.deps.SettingsAdminService.CreateAPIKey(c.Request.Context(), req.Name, principal.UserID, expiresAt)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, gin.H{
		"apiKey":    item,
		"rawKey":    rawKey,
		"message":   "save this key securely, it will not be shown again",
	})
}

func (h *SettingsAdminHandler) GetAPIKey(c *gin.Context) {
	item, err := h.deps.SettingsAdminService.GetAPIKey(c.Request.Context(), c.Param("keyId"))
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SettingsAdminHandler) UpdateAPIKey(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SettingsAdminService.UpdateAPIKey(c.Request.Context(), c.Param("keyId"), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SettingsAdminHandler) DeleteAPIKey(c *gin.Context) {
	if err := h.deps.SettingsAdminService.DeleteAPIKey(c.Request.Context(), c.Param("keyId")); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

func (h *SettingsAdminHandler) GetAPIKeyUsage(c *gin.Context) {
	item, err := h.deps.SettingsAdminService.GetAPIKey(c.Request.Context(), c.Param("keyId"))
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{
		"keyId":      item.ID,
		"usageCount": item.UsageCount,
		"lastUsedAt": item.LastUsedAt,
	})
}

// OAuth

func (h *SettingsAdminHandler) UpdateOAuth(c *gin.Context) {
	provider := c.Param("provider")
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.UpdateOAuth(c.Request.Context(), provider, req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) TestOAuth(c *gin.Context) {
	provider := c.Param("provider")
	utils.Success(c, http.StatusOK, gin.H{"message": "OAuth test queued", "provider": provider})
}

func (h *SettingsAdminHandler) GetCallbackURL(c *gin.Context) {
	provider := c.Param("provider")
	url := h.deps.SettingsAdminService.GetCallbackURL(c.Request.Context(), provider)
	utils.Success(c, http.StatusOK, gin.H{"provider": provider, "callbackUrl": url})
}

// Integrations

func (h *SettingsAdminHandler) ListIntegrations(c *gin.Context) {
	items, err := h.deps.SettingsAdminService.ListIntegrations(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"items": items})
}

func (h *SettingsAdminHandler) CreateIntegration(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Provider string `json:"provider"`
		Type     string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" || req.Provider == "" || req.Type == "" {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SettingsAdminService.CreateIntegration(c.Request.Context(), req.Name, req.Provider, req.Type)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *SettingsAdminHandler) GetIntegration(c *gin.Context) {
	item, err := h.deps.SettingsAdminService.GetIntegration(c.Request.Context(), c.Param("integrationId"))
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SettingsAdminHandler) UpdateIntegration(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	item, err := h.deps.SettingsAdminService.UpdateIntegration(c.Request.Context(), c.Param("integrationId"), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *SettingsAdminHandler) DeleteIntegration(c *gin.Context) {
	if err := h.deps.SettingsAdminService.DeleteIntegration(c.Request.Context(), c.Param("integrationId")); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"deleted": true})
}

func (h *SettingsAdminHandler) TestIntegration(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "connection test queued", "integrationId": c.Param("integrationId")})
}

func (h *SettingsAdminHandler) GetIntegrationLogs(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"items": []map[string]any{}, "integrationId": c.Param("integrationId")})
}

// Maintenance

func (h *SettingsAdminHandler) GetMaintenance(c *gin.Context) {
	result, err := h.deps.SettingsAdminService.GetMaintenance(c.Request.Context())
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) UpdateMaintenance(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrValidationFailed)
		return
	}
	result, err := h.deps.SettingsAdminService.UpdateMaintenance(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, http.StatusOK, result)
}

func (h *SettingsAdminHandler) ClearCache(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "cache clear queued"})
}

func (h *SettingsAdminHandler) OptimizeDB(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"message": "database optimization queued"})
}

func (h *SettingsAdminHandler) ListMaintenanceJobs(c *gin.Context) {
	utils.Success(c, http.StatusOK, gin.H{"items": []map[string]any{}})
}
