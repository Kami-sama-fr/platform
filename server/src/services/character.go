package services

import (
	"context"
	"strings"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type CharacterService struct {
	repos *Repositories
}

func NewCharacterService(repos *Repositories) *CharacterService {
	return &CharacterService{repos: repos}
}

func (s *CharacterService) Create(ctx context.Context, userID, name, slug, imageUrl, description, gender string) (*models.Character, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, utils.ErrValidationFailed
	}
	if slug == "" {
		slug = generateSlug(name)
	}
	now := time.Now().UTC()
	character := &models.Character{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		Name:        name,
		Slug:        slug,
		ImageUrl:    imageUrl,
		Description: description,
		Gender:      gender,
	}
	if err := s.repos.Characters().Create(ctx, character); err != nil {
		return nil, err
	}
	return character, nil
}

func (s *CharacterService) List(ctx context.Context) ([]models.Character, error) {
	return s.repos.Characters().List(ctx)
}

func (s *CharacterService) GetByID(ctx context.Context, id string) (*models.Character, error) {
	return s.repos.Characters().GetByID(ctx, id)
}

func (s *CharacterService) GetBySlug(ctx context.Context, slug string) (*models.Character, error) {
	return s.repos.Characters().GetBySlug(ctx, slug)
}

func (s *CharacterService) Update(ctx context.Context, userID, id string, req struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	ImageUrl    *string `json:"imageUrl"`
	Description *string `json:"description"`
	Gender      *string `json:"gender"`
}) (*models.Character, error) {
	character, err := s.repos.Characters().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		character.Name = *req.Name
	}
	if req.Slug != nil {
		character.Slug = *req.Slug
	}
	if req.ImageUrl != nil {
		character.ImageUrl = *req.ImageUrl
	}
	if req.Description != nil {
		character.Description = *req.Description
	}
	if req.Gender != nil {
		character.Gender = *req.Gender
	}
	character.UpdatedAt = time.Now().UTC()
	if err := s.repos.Characters().Update(ctx, character); err != nil {
		return nil, err
	}
	return character, nil
}

func (s *CharacterService) Delete(ctx context.Context, id string) error {
	return s.repos.Characters().Delete(ctx, id)
}
