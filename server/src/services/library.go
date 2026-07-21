package services

import (
	"context"
	"strings"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
	"github.com/kami-sama-fr/platform/server/src/utils"
	"gorm.io/datatypes"
)

type LibraryService struct {
	repos *Repositories
}

func NewLibraryService(repos *Repositories) *LibraryService {
	return &LibraryService{repos: repos}
}

func (s *LibraryService) Create(ctx context.Context, userID, sourceType string, enabled bool, config datatypes.JSON) (*models.SourceConfig, error) {
	sourceType = strings.TrimSpace(sourceType)
	if sourceType == "" {
		return nil, utils.ErrValidationFailed
	}
	now := time.Now().UTC()
	library := &models.SourceConfig{
		Common:     models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		SourceType: sourceType,
		Enabled:    enabled,
		Config:     config,
	}
	if err := s.repos.Libraries().Create(ctx, library); err != nil {
		return nil, err
	}
	return library, nil
}

func (s *LibraryService) List(ctx context.Context) ([]models.SourceConfig, error) {
	return s.repos.Libraries().List(ctx)
}

func (s *LibraryService) GetByID(ctx context.Context, id string) (*models.SourceConfig, error) {
	return s.repos.Libraries().GetByID(ctx, id)
}

func (s *LibraryService) GetBySourceType(ctx context.Context, sourceType string) (*models.SourceConfig, error) {
	return s.repos.Libraries().GetBySourceType(ctx, sourceType)
}

func (s *LibraryService) Update(ctx context.Context, userID, id string, req struct {
	SourceType *string         `json:"sourceType"`
	Enabled    *bool           `json:"enabled"`
	Config     datatypes.JSON  `json:"config"`
}) (*models.SourceConfig, error) {
	library, err := s.repos.Libraries().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.SourceType != nil {
		library.SourceType = *req.SourceType
	}
	if req.Enabled != nil {
		library.Enabled = *req.Enabled
	}
	if req.Config != nil {
		library.Config = req.Config
	}
	library.UpdatedAt = time.Now().UTC()
	if err := s.repos.Libraries().Update(ctx, library); err != nil {
		return nil, err
	}
	return library, nil
}

func (s *LibraryService) Delete(ctx context.Context, id string) error {
	return s.repos.Libraries().Delete(ctx, id)
}
