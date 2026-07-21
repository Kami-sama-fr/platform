package services

import (
	"context"
	"strings"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type TagService struct {
	repos *Repositories
}

func NewTagService(repos *Repositories) *TagService {
	return &TagService{repos: repos}
}

func (s *TagService) Create(ctx context.Context, userID, name, slug string) (*models.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, utils.ErrValidationFailed
	}
	if slug == "" {
		slug = generateSlug(name)
	}
	now := time.Now().UTC()
	tag := &models.Tag{
		Common: models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		Name:   name,
		Slug:   slug,
	}
	if err := s.repos.Tags().Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) List(ctx context.Context) ([]models.Tag, error) {
	return s.repos.Tags().List(ctx)
}

func (s *TagService) GetByID(ctx context.Context, id string) (*models.Tag, error) {
	return s.repos.Tags().GetByID(ctx, id)
}

func (s *TagService) GetBySlug(ctx context.Context, slug string) (*models.Tag, error) {
	return s.repos.Tags().GetBySlug(ctx, slug)
}

func (s *TagService) Update(ctx context.Context, userID, id string, req struct {
	Name *string `json:"name"`
	Slug *string `json:"slug"`
}) (*models.Tag, error) {
	tag, err := s.repos.Tags().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		tag.Name = *req.Name
	}
	if req.Slug != nil {
		tag.Slug = *req.Slug
	}
	tag.UpdatedAt = time.Now().UTC()
	if err := s.repos.Tags().Update(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) Delete(ctx context.Context, id string) error {
	return s.repos.Tags().Delete(ctx, id)
}
