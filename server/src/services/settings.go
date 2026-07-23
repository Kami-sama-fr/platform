package services

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/skygenesisenterprise/kami-sama/server/src/models"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
	"gorm.io/datatypes"
)

type SettingsService struct {
	repos *Repositories
}

func NewSettingsService(repos *Repositories) *SettingsService {
	return &SettingsService{repos: repos}
}

func (s *SettingsService) GetSetting(ctx context.Context, key string) (*models.SystemSetting, error) {
	return s.repos.SystemSettings().GetByKey(ctx, key)
}

func (s *SettingsService) ListSettings(ctx context.Context, category string) ([]models.SystemSetting, error) {
	return s.repos.SystemSettings().List(ctx, category)
}

func (s *SettingsService) UpsertSetting(ctx context.Context, key string, value any, category, description string) (*models.SystemSetting, error) {
	if strings.TrimSpace(key) == "" {
		return nil, utils.ErrValidationFailed
	}
	var jsonValue []byte
	switch v := value.(type) {
	case []byte:
		jsonValue = v
	case string:
		jsonValue = []byte(v)
	default:
		var err error
		jsonValue, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}
	now := time.Now().UTC()
	setting := &models.SystemSetting{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		Key:         key,
		Category:    category,
		Description: description,
		Value:       datatypes.JSON(jsonValue),
	}
	if err := s.repos.SystemSettings().Upsert(ctx, setting); err != nil {
		return nil, err
	}
	return setting, nil
}

func (s *SettingsService) DeleteSetting(ctx context.Context, key string) error {
	return s.repos.SystemSettings().Delete(ctx, key)
}

func (s *SettingsService) GetSeoMeta(ctx context.Context, pagePath string) (*models.SeoMeta, error) {
	var meta models.SeoMeta
	err := s.repos.db.WithContext(ctx).Where("page_path = ?", pagePath).First(&meta).Error
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

func (s *SettingsService) UpsertSeoMeta(ctx context.Context, pagePath, title, description, keywords, ogImage, canonicalUrl string) (*models.SeoMeta, error) {
	now := time.Now().UTC()
	meta := &models.SeoMeta{
		Common:       models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		PagePath:     pagePath,
		Title:        title,
		Description:  description,
		Keywords:     keywords,
		OgImage:      ogImage,
		CanonicalUrl: canonicalUrl,
	}
	if err := s.repos.db.WithContext(ctx).Where("page_path = ?", pagePath).Save(meta).Error; err != nil {
		return nil, err
	}
	return meta, nil
}
