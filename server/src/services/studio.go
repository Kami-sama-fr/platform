package services

import (
	"context"
	"strings"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type StudioService struct {
	repos *Repositories
}

func NewStudioService(repos *Repositories) *StudioService {
	return &StudioService{repos: repos}
}

func (s *StudioService) Create(ctx context.Context, userID, name, slug, logoUrl, description, website string) (*models.Studio, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, utils.ErrValidationFailed
	}
	if slug == "" {
		slug = generateSlug(name)
	}
	now := time.Now().UTC()
	studio := &models.Studio{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		Name:        name,
		Slug:        slug,
		LogoUrl:     logoUrl,
		Description: description,
		Website:     website,
	}
	if err := s.repos.Studios().Create(ctx, studio); err != nil {
		return nil, err
	}
	return studio, nil
}

func (s *StudioService) List(ctx context.Context) ([]models.Studio, error) {
	return s.repos.Studios().List(ctx)
}

func (s *StudioService) GetByID(ctx context.Context, id string) (*models.Studio, error) {
	return s.repos.Studios().GetByID(ctx, id)
}

func (s *StudioService) GetBySlug(ctx context.Context, slug string) (*models.Studio, error) {
	return s.repos.Studios().GetBySlug(ctx, slug)
}

func (s *StudioService) Update(ctx context.Context, userID, id string, req struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	LogoUrl     *string `json:"logoUrl"`
	Description *string `json:"description"`
	Website     *string `json:"website"`
}) (*models.Studio, error) {
	studio, err := s.repos.Studios().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		studio.Name = *req.Name
	}
	if req.Slug != nil {
		studio.Slug = *req.Slug
	}
	if req.LogoUrl != nil {
		studio.LogoUrl = *req.LogoUrl
	}
	if req.Description != nil {
		studio.Description = *req.Description
	}
	if req.Website != nil {
		studio.Website = *req.Website
	}
	studio.UpdatedAt = time.Now().UTC()
	if err := s.repos.Studios().Update(ctx, studio); err != nil {
		return nil, err
	}
	return studio, nil
}

func (s *StudioService) Delete(ctx context.Context, id string) error {
	return s.repos.Studios().Delete(ctx, id)
}
