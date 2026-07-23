package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/skygenesisenterprise/kami-sama/server/src/models"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SettingsAdminService struct {
	db    *gorm.DB
	repos *Repositories
}

func NewSettingsAdminService(db *gorm.DB, repos *Repositories) *SettingsAdminService {
	return &SettingsAdminService{db: db, repos: repos}
}

// General

func (s *SettingsAdminService) GetGeneral(ctx context.Context) (map[string]any, error) {
	items, _ := s.repos.SystemSettings().List(ctx, "general")
	result := map[string]any{}
	for _, item := range items {
		var v any
		json.Unmarshal(item.Value, &v)
		result[item.Key] = v
	}
	if _, ok := result["name"]; !ok {
		result["name"] = ""
	}
	if _, ok := result["description"]; !ok {
		result["description"] = ""
	}
	if _, ok := result["locale"]; !ok {
		result["locale"] = "en"
	}
	return result, nil
}

func (s *SettingsAdminService) UpdateGeneral(ctx context.Context, data map[string]any) (map[string]any, error) {
	for k, v := range data {
		desc := ""
		switch k {
		case "name":
			desc = "Site name"
		case "description":
			desc = "Site description"
		case "locale":
			desc = "Default locale"
		}
		if err := s.repos.SystemSettings().Upsert(ctx, &models.SystemSetting{
			Common:      models.Common{ID: utils.NewID(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
			Key:         k,
			Category:    "general",
			Description: desc,
			Value:       datatypes.JSON(mustMarshal(v)),
		}); err != nil {
			return nil, err
		}
	}
	return s.GetGeneral(ctx)
}

// Security

func (s *SettingsAdminService) GetSecurity(ctx context.Context) (map[string]any, error) {
	return s.getCategoryMap(ctx, "security")
}

func (s *SettingsAdminService) UpdateSecurity(ctx context.Context, data map[string]any) (map[string]any, error) {
	return s.updateCategory(ctx, "security", data)
}

func (s *SettingsAdminService) GetSessions(ctx context.Context) (map[string]any, error) {
	return s.getCategoryMap(ctx, "sessions")
}

func (s *SettingsAdminService) UpdateSessions(ctx context.Context, data map[string]any) (map[string]any, error) {
	return s.updateCategory(ctx, "sessions", data)
}

func (s *SettingsAdminService) GetRateLimit(ctx context.Context) (map[string]any, error) {
	return s.getCategoryMap(ctx, "rate_limit")
}

func (s *SettingsAdminService) UpdateRateLimit(ctx context.Context, data map[string]any) (map[string]any, error) {
	return s.updateCategory(ctx, "rate_limit", data)
}

func (s *SettingsAdminService) Get2FA(ctx context.Context) (map[string]any, error) {
	return s.getCategoryMap(ctx, "2fa")
}

func (s *SettingsAdminService) Update2FA(ctx context.Context, data map[string]any) (map[string]any, error) {
	return s.updateCategory(ctx, "2fa", data)
}

// Branding

func (s *SettingsAdminService) GetBranding(ctx context.Context) (map[string]any, error) {
	return s.getCategoryMap(ctx, "branding")
}

func (s *SettingsAdminService) UpdateBranding(ctx context.Context, data map[string]any) (map[string]any, error) {
	return s.updateCategory(ctx, "branding", data)
}

// Email

func (s *SettingsAdminService) GetEmail(ctx context.Context) (map[string]any, error) {
	return s.getCategoryMap(ctx, "email")
}

func (s *SettingsAdminService) UpdateEmail(ctx context.Context, data map[string]any) (map[string]any, error) {
	return s.updateCategory(ctx, "email", data)
}

// Storage

func (s *SettingsAdminService) GetStorage(ctx context.Context) (map[string]any, error) {
	return s.getCategoryMap(ctx, "storage")
}

func (s *SettingsAdminService) UpdateStorage(ctx context.Context, data map[string]any) (map[string]any, error) {
	return s.updateCategory(ctx, "storage", data)
}

// CDN

func (s *SettingsAdminService) GetCDN(ctx context.Context) (map[string]any, error) {
	return s.getCategoryMap(ctx, "cdn")
}

func (s *SettingsAdminService) UpdateCDN(ctx context.Context, data map[string]any) (map[string]any, error) {
	return s.updateCategory(ctx, "cdn", data)
}

// Maintenance

func (s *SettingsAdminService) GetMaintenance(ctx context.Context) (map[string]any, error) {
	result, _ := s.getCategoryMap(ctx, "maintenance")
	if _, ok := result["enabled"]; !ok {
		result["enabled"] = false
	}
	if _, ok := result["message"]; !ok {
		result["message"] = ""
	}
	return result, nil
}

func (s *SettingsAdminService) UpdateMaintenance(ctx context.Context, data map[string]any) (map[string]any, error) {
	return s.updateCategory(ctx, "maintenance", data)
}

// Domains

func (s *SettingsAdminService) ListDomains(ctx context.Context) ([]models.DomainConfig, error) {
	return s.repos.DomainConfigs().List(ctx)
}

func (s *SettingsAdminService) CreateDomain(ctx context.Context, domain string) (*models.DomainConfig, error) {
	item := &models.DomainConfig{
		Common:     models.Common{ID: utils.NewID(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		Domain:     domain,
		IsVerified: false,
		IsSSL:      false,
	}
	if err := s.repos.DomainConfigs().Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *SettingsAdminService) GetDomain(ctx context.Context, domainID string) (*models.DomainConfig, error) {
	return s.repos.DomainConfigs().GetByID(ctx, domainID)
}

func (s *SettingsAdminService) DeleteDomain(ctx context.Context, domainID string) error {
	return s.repos.DomainConfigs().Delete(ctx, domainID)
}

// API Keys

func (s *SettingsAdminService) ListAPIKeys(ctx context.Context) ([]models.ApiKey, error) {
	return s.repos.ApiKeys().List(ctx)
}

func (s *SettingsAdminService) CreateAPIKey(ctx context.Context, name string, userID string, expiresAt *time.Time) (*models.ApiKey, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, "", err
	}
	rawKey := fmt.Sprintf("sk_%s", hex.EncodeToString(b))
	hash := sha256Hex(rawKey)

	item := &models.ApiKey{
		Common:    models.Common{ID: utils.NewID(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		Name:      name,
		KeyHash:   hash,
		KeyPrefix: "sk_" + rawKey[3:11],
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	if err := s.repos.ApiKeys().Create(ctx, item); err != nil {
		return nil, "", err
	}
	return item, rawKey, nil
}

func (s *SettingsAdminService) GetAPIKey(ctx context.Context, keyID string) (*models.ApiKey, error) {
	return s.repos.ApiKeys().GetByID(ctx, keyID)
}

func (s *SettingsAdminService) UpdateAPIKey(ctx context.Context, keyID string, data map[string]any) (*models.ApiKey, error) {
	item, err := s.repos.ApiKeys().GetByID(ctx, keyID)
	if err != nil {
		return nil, err
	}
	if v, ok := data["name"]; ok {
		name, _ := v.(string)
		item.Name = name
	}
	if v, ok := data["permissions"]; ok {
		b, _ := json.Marshal(v)
		item.Permissions = datatypes.JSON(b)
	}
	if v, ok := data["expiresAt"]; ok {
		if str, ok2 := v.(string); ok2 && str != "" {
			t, err2 := time.Parse(time.RFC3339, str)
			if err2 != nil {
				return nil, utils.NewError(http.StatusBadRequest, "VALIDATION_ERROR", "Invalid expiresAt format", nil)
			}
			item.ExpiresAt = &t
		} else {
			item.ExpiresAt = nil
		}
	}
	item.UpdatedAt = time.Now().UTC()
	if err := s.repos.ApiKeys().Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *SettingsAdminService) DeleteAPIKey(ctx context.Context, keyID string) error {
	key, err := s.repos.ApiKeys().GetByID(ctx, keyID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	key.RevokedAt = &now
	return s.repos.ApiKeys().Update(ctx, key)
}
func (s *SettingsAdminService) GetAPIKeyWithUsage(ctx context.Context, keyID string) (*models.ApiKey, error) {
	return s.repos.ApiKeys().GetByID(ctx, keyID)
}

// OAuth

func (s *SettingsAdminService) UpdateOAuth(ctx context.Context, provider string, data map[string]any) (map[string]any, error) {
	return s.updateCategory(ctx, fmt.Sprintf("oauth_%s", provider), data)
}

func (s *SettingsAdminService) GetCallbackURL(ctx context.Context, provider string) string {
	return fmt.Sprintf("/auth/oauth/%s/callback", provider)
}

// Integrations

func (s *SettingsAdminService) ListIntegrations(ctx context.Context) ([]models.Integration, error) {
	return s.repos.Integrations().List(ctx)
}

func (s *SettingsAdminService) CreateIntegration(ctx context.Context, name, provider, integrationType string) (*models.Integration, error) {
	item := &models.Integration{
		Common:   models.Common{ID: utils.NewID(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		Name:     name,
		Provider: provider,
		Type:     integrationType,
		Status:   "inactive",
	}
	if err := s.repos.Integrations().Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *SettingsAdminService) GetIntegration(ctx context.Context, integrationID string) (*models.Integration, error) {
	return s.repos.Integrations().GetByID(ctx, integrationID)
}

func (s *SettingsAdminService) UpdateIntegration(ctx context.Context, integrationID string, data map[string]any) (*models.Integration, error) {
	item, err := s.repos.Integrations().GetByID(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	if v, ok := data["name"]; ok {
		name, _ := v.(string)
		item.Name = name
	}
	if v, ok := data["status"]; ok {
		status, _ := v.(string)
		item.Status = status
	}
	if v, ok := data["config"]; ok {
		b, _ := json.Marshal(v)
		item.Config = datatypes.JSON(b)
	}
	if v, ok := data["metadata"]; ok {
		b, _ := json.Marshal(v)
		item.Metadata = datatypes.JSON(b)
	}
	item.UpdatedAt = time.Now().UTC()
	if err := s.repos.Integrations().Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *SettingsAdminService) DeleteIntegration(ctx context.Context, integrationID string) error {
	return s.repos.Integrations().Delete(ctx, integrationID)
}

// SEO

func (s *SettingsAdminService) ListSEOPages(ctx context.Context) ([]models.SeoMeta, error) {
	var items []models.SeoMeta
	err := s.db.WithContext(ctx).Order("page_path ASC").Find(&items).Error
	return items, err
}

func (s *SettingsAdminService) UpdateSEOPage(ctx context.Context, pagePath string, data map[string]string) (*models.SeoMeta, error) {
	var meta models.SeoMeta
	err := s.db.WithContext(ctx).Where("page_path = ?", pagePath).First(&meta).Error
	now := time.Now().UTC()
	if err != nil {
		meta = models.SeoMeta{
			Common: models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		}
		meta.PagePath = pagePath
	}
	if v, ok := data["title"]; ok {
		meta.Title = v
	}
	if v, ok := data["description"]; ok {
		meta.Description = v
	}
	if v, ok := data["keywords"]; ok {
		meta.Keywords = v
	}
	if v, ok := data["ogImage"]; ok {
		meta.OgImage = v
	}
	if v, ok := data["canonicalUrl"]; ok {
		meta.CanonicalUrl = v
	}
	meta.UpdatedAt = now
	if err := s.db.WithContext(ctx).Save(&meta).Error; err != nil {
		return nil, err
	}
	return &meta, nil
}

// Sitemap & robots

func (s *SettingsAdminService) GenerateSitemap(ctx context.Context) (string, error) {
	var paths []string
	s.db.WithContext(ctx).Model(&models.SeoMeta{}).Pluck("page_path", &paths)
	sb := &strings.Builder{}
	sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	sb.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	for _, p := range paths {
		sb.WriteString(fmt.Sprintf("  <url><loc>%s</loc></url>\n", p))
	}
	sb.WriteString("</urlset>\n")
	return sb.String(), nil
}

func (s *SettingsAdminService) GetRobots(ctx context.Context) (string, error) {
	setting, err := s.repos.SystemSettings().GetByKey(ctx, "robots_txt")
	if err != nil {
		return "User-agent: *\nAllow: /\n", nil
	}
	var v string
	if err := json.Unmarshal(setting.Value, &v); err != nil {
		return "User-agent: *\nAllow: /\n", nil
	}
	return v, nil
}

// helpers

func (s *SettingsAdminService) getCategoryMap(ctx context.Context, category string) (map[string]any, error) {
	items, _ := s.repos.SystemSettings().List(ctx, category)
	result := map[string]any{}
	for _, item := range items {
		var v any
		json.Unmarshal(item.Value, &v)
		result[item.Key] = v
	}
	return result, nil
}

func (s *SettingsAdminService) updateCategory(ctx context.Context, category string, data map[string]any) (map[string]any, error) {
	for k, v := range data {
		if err := s.repos.SystemSettings().Upsert(ctx, &models.SystemSetting{
			Common:   models.Common{ID: utils.NewID(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
			Key:      k,
			Category: category,
			Value:    datatypes.JSON(mustMarshal(v)),
		}); err != nil {
			return nil, err
		}
	}
	return s.getCategoryMap(ctx, category)
}

func mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		b, _ = json.Marshal(fmt.Sprintf("%v", v))
	}
	return b
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
