package services

import (
	"context"
	"strings"
	"time"

	"github.com/skygenesisenterprise/kami-sama/server/src/models"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type EpisodeService struct {
	repos *Repositories
}

func NewEpisodeService(repos *Repositories) *EpisodeService {
	return &EpisodeService{repos: repos}
}

type UpdateEpisodeInput struct {
	Number       *int     `json:"number"`
	Title        *string  `json:"title"`
	Synopsis     *string  `json:"synopsis"`
	ThumbnailUrl *string  `json:"thumbnailUrl"`
	Duration     *float64 `json:"duration"`
	AirDate      *string  `json:"airDate"`
	IsSubbed     *bool    `json:"isSubbed"`
	IsDubbed     *bool    `json:"isDubbed"`
}

func (s *EpisodeService) Create(ctx context.Context, userID, animeID string, seasonID *string, number int, title, synopsis, thumbnailUrl string, duration float64, airDate string, isSubbed, isDubbed bool) (*models.Episode, error) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(animeID) == "" {
		return nil, utils.ErrValidationFailed
	}
	now := time.Now().UTC()
	episode := &models.Episode{
		Common:       models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		AnimeID:      animeID,
		SeasonID:     seasonID,
		Number:       number,
		Title:        strings.TrimSpace(title),
		Synopsis:     synopsis,
		ThumbnailUrl: thumbnailUrl,
		Duration:     duration,
		IsSubbed:     isSubbed,
		IsDubbed:     isDubbed,
	}
	if airDate != "" {
		t, err := time.Parse("2006-01-02", airDate)
		if err == nil {
			episode.AirDate = &t
		}
	}
	if err := s.repos.Episodes().Create(ctx, episode); err != nil {
		return nil, err
	}
	return episode, nil
}

func (s *EpisodeService) GetByID(ctx context.Context, id string) (*models.Episode, error) {
	return s.repos.Episodes().GetByID(ctx, id)
}

func (s *EpisodeService) GetByNumber(ctx context.Context, animeID string, number int) (*models.Episode, error) {
	var episode models.Episode
	err := s.repos.db.WithContext(ctx).
		Where("anime_id = ? AND number = ?", animeID, number).
		First(&episode).Error
	if err != nil {
		return nil, err
	}
	return &episode, nil
}

func (s *EpisodeService) ListByAnime(ctx context.Context, animeID string, seasonID *string) ([]models.Episode, error) {
	return s.repos.Episodes().ListByAnime(ctx, animeID, seasonID)
}

func (s *EpisodeService) Update(ctx context.Context, userID, id string, input UpdateEpisodeInput) (*models.Episode, error) {
	episode, err := s.repos.Episodes().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Number != nil {
		episode.Number = *input.Number
	}
	if input.Title != nil {
		episode.Title = strings.TrimSpace(*input.Title)
	}
	if input.Synopsis != nil {
		episode.Synopsis = *input.Synopsis
	}
	if input.ThumbnailUrl != nil {
		episode.ThumbnailUrl = *input.ThumbnailUrl
	}
	if input.Duration != nil {
		episode.Duration = *input.Duration
	}
	if input.AirDate != nil {
		t, err := time.Parse("2006-01-02", *input.AirDate)
		if err == nil {
			episode.AirDate = &t
		}
	}
	if input.IsSubbed != nil {
		episode.IsSubbed = *input.IsSubbed
	}
	if input.IsDubbed != nil {
		episode.IsDubbed = *input.IsDubbed
	}
	episode.UpdatedAt = time.Now().UTC()
	if err := s.repos.Episodes().Update(ctx, episode); err != nil {
		return nil, err
	}
	return episode, nil
}

func (s *EpisodeService) Delete(ctx context.Context, id string) error {
	return s.repos.Episodes().Delete(ctx, id)
}
