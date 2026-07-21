package services

import (
	"context"

	"github.com/kami-sama-fr/platform/server/src/models"
)

type SearchService struct {
	repos *Repositories
}

func NewSearchService(repos *Repositories) *SearchService {
	return &SearchService{repos: repos}
}

type SearchResult struct {
	Anime      []models.Anime     `json:"anime"`
	Characters []models.Character `json:"characters"`
	Studios    []models.Studio    `json:"studios"`
}

func (s *SearchService) Search(ctx context.Context, query string, limit int) (*SearchResult, error) {
	if limit <= 0 {
		limit = 10
	}
	anime, _ := s.repos.Anime().Search(ctx, query, limit)
	characters, _ := s.repos.Characters().Search(ctx, query, limit)
	studios, _ := s.repos.Studios().Search(ctx, query, limit)
	return &SearchResult{
		Anime:      anime,
		Characters: characters,
		Studios:    studios,
	}, nil
}

func (s *SearchService) SearchAnime(ctx context.Context, query string, limit int) ([]models.Anime, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repos.Anime().Search(ctx, query, limit)
}

func (s *SearchService) SearchCharacters(ctx context.Context, query string, limit int) ([]models.Character, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repos.Characters().Search(ctx, query, limit)
}

func (s *SearchService) SearchStudios(ctx context.Context, query string, limit int) ([]models.Studio, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repos.Studios().Search(ctx, query, limit)
}

func (s *SearchService) Suggestions(ctx context.Context, query string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}
	anime, err := s.repos.Anime().Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	suggestions := make([]string, 0, len(anime))
	for _, a := range anime {
		suggestions = append(suggestions, a.Title)
	}
	return suggestions, nil
}
