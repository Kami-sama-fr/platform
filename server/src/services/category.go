package services

import (
	"context"
	"strings"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type CategoryService struct {
	repos *Repositories
}

func NewCategoryService(repos *Repositories) *CategoryService {
	return &CategoryService{repos: repos}
}

func (s *CategoryService) Create(ctx context.Context, userID, name, slug, description string, parentID *string) (*models.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, utils.ErrValidationFailed
	}
	if slug == "" {
		slug = generateSlug(name)
	}
	now := time.Now().UTC()
	category := &models.Category{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		Name:        name,
		Slug:        slug,
		Description: description,
		ParentID:    parentID,
		IsActive:    true,
	}
	if err := s.repos.Categories().Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) List(ctx context.Context) ([]models.Category, error) {
	return s.repos.Categories().List(ctx)
}

func (s *CategoryService) GetByID(ctx context.Context, id string) (*models.Category, error) {
	return s.repos.Categories().GetByID(ctx, id)
}

func (s *CategoryService) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	return s.repos.Categories().GetBySlug(ctx, slug)
}

func (s *CategoryService) Update(ctx context.Context, userID, id string, req struct {
	Name        *string  `json:"name"`
	Slug        *string  `json:"slug"`
	Description *string  `json:"description"`
	ParentID    *string  `json:"parentId"`
	SortOrder   *int     `json:"sortOrder"`
	IsActive    *bool    `json:"isActive"`
}) (*models.Category, error) {
	category, err := s.repos.Categories().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	category.UpdatedAt = time.Now().UTC()
	if err := s.repos.Categories().Update(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	return s.repos.Categories().Delete(ctx, id)
}
