package services

import (
	"context"
	"strings"
	"time"

	"github.com/skygenesisenterprise/kami-sama/server/src/models"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type PremiereService struct {
	repos *Repositories
}

func NewPremiereService(repos *Repositories) *PremiereService {
	return &PremiereService{repos: repos}
}

type UpdatePremiereInput struct {
	Season     *int    `json:"season"`
	EpisodeNum *int    `json:"episodeNum"`
	ScheduledAt *string `json:"scheduledAt"`
	Region     *string `json:"region"`
	Status     *string `json:"status"`
}

func (s *PremiereService) List(ctx context.Context) ([]models.Premiere, error) {
	return s.repos.Premieres().List(ctx)
}

func (s *PremiereService) GetByID(ctx context.Context, id string) (*models.Premiere, error) {
	return s.repos.Premieres().GetByID(ctx, id)
}

func (s *PremiereService) Create(ctx context.Context, userID, animeID string, season, episodeNum int, scheduledAt, region string) (*models.Premiere, error) {
	if strings.TrimSpace(animeID) == "" {
		return nil, utils.ErrValidationFailed
	}

	parsedAt, err := time.Parse(time.RFC3339, scheduledAt)
	if err != nil {
		parsedAt, err = time.Parse("2006-01-02", scheduledAt)
		if err != nil {
			return nil, utils.ErrValidationFailed
		}
	}

	now := time.Now().UTC()
	premiere := &models.Premiere{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		AnimeID:     animeID,
		Season:      season,
		EpisodeNum:  episodeNum,
		ScheduledAt: parsedAt,
		Region:      defaultString(region, "global"),
		Status:      "scheduled",
	}
	if err := s.repos.Premieres().Create(ctx, premiere); err != nil {
		return nil, err
	}
	return premiere, nil
}

func (s *PremiereService) Update(ctx context.Context, userID, id string, req UpdatePremiereInput) (*models.Premiere, error) {
	premiere, err := s.repos.Premieres().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Season != nil {
		premiere.Season = *req.Season
	}
	if req.EpisodeNum != nil {
		premiere.EpisodeNum = *req.EpisodeNum
	}
	if req.ScheduledAt != nil {
		parsedAt, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err != nil {
			return nil, utils.ErrValidationFailed
		}
		premiere.ScheduledAt = parsedAt
	}
	if req.Region != nil {
		premiere.Region = *req.Region
	}
	if req.Status != nil {
		premiere.Status = *req.Status
	}
	premiere.UpdatedAt = time.Now().UTC()
	if err := s.repos.Premieres().Update(ctx, premiere); err != nil {
		return nil, err
	}
	return premiere, nil
}

func (s *PremiereService) Delete(ctx context.Context, id string) error {
	return s.repos.Premieres().Delete(ctx, id)
}
