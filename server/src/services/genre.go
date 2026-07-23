package services

import (
	"context"
	"strings"
	"time"

	"github.com/skygenesisenterprise/kami-sama/server/src/models"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type GenreService struct {
	repos *Repositories
}

func NewGenreService(repos *Repositories) *GenreService {
	return &GenreService{repos: repos}
}

func (s *GenreService) Create(ctx context.Context, userID, name, slug, description string) (*models.Genre, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, utils.ErrValidationFailed
	}
	if slug == "" {
		slug = generateSlug(name)
	}
	now := time.Now().UTC()
	genre := &models.Genre{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		Name:        name,
		Slug:        slug,
		Description: description,
	}
	if err := s.repos.Genres().Create(ctx, genre); err != nil {
		return nil, err
	}
	return genre, nil
}

func (s *GenreService) List(ctx context.Context) ([]models.Genre, error) {
	return s.repos.Genres().List(ctx)
}

func (s *GenreService) GetByID(ctx context.Context, id string) (*models.Genre, error) {
	return s.repos.Genres().GetByID(ctx, id)
}

func (s *GenreService) GetBySlug(ctx context.Context, slug string) (*models.Genre, error) {
	return s.repos.Genres().GetBySlug(ctx, slug)
}

func (s *GenreService) Update(ctx context.Context, userID, id string, req struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
}) (*models.Genre, error) {
	genre, err := s.repos.Genres().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		genre.Name = *req.Name
	}
	if req.Slug != nil {
		genre.Slug = *req.Slug
	}
	if req.Description != nil {
		genre.Description = *req.Description
	}
	genre.UpdatedAt = time.Now().UTC()
	if err := s.repos.Genres().Update(ctx, genre); err != nil {
		return nil, err
	}
	return genre, nil
}

func (s *GenreService) Delete(ctx context.Context, id string) error {
	return s.repos.Genres().Delete(ctx, id)
}
