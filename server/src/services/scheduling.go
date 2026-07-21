package services

import (
	"context"
	"strings"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type SchedulingService struct {
	repos *Repositories
}

func NewSchedulingService(repos *Repositories) *SchedulingService {
	return &SchedulingService{repos: repos}
}

type UpdateSimulcastReq struct {
	AirDay      *string `json:"airDay"`
	AirTime     *string `json:"airTime"`
	AirTimezone *string `json:"airTimezone"`
	Region      *string `json:"region"`
	Platform    *string `json:"platform"`
	Status      *string `json:"status"`
}

type UpdateReleaseScheduleReq struct {
	ScheduledAt *string `json:"scheduledAt"`
	Status      *string `json:"status"`
}

func (s *SchedulingService) ListActiveSimulcasts(ctx context.Context) ([]models.Simulcast, error) {
	return s.repos.Simulcasts().ListActive(ctx)
}

func (s *SchedulingService) GetSimulcast(ctx context.Context, id string) (*models.Simulcast, error) {
	return s.repos.Simulcasts().GetByID(ctx, id)
}

func (s *SchedulingService) ListSimulcastsByWeek(ctx context.Context) (map[string][]models.Simulcast, error) {
	return s.repos.Simulcasts().ListByWeek(ctx)
}

func (s *SchedulingService) CreateSimulcast(ctx context.Context, userID, animeID, airDay, airTime, airTimezone, region, platform string, episodeNumber, season int, startsAt string) (*models.Simulcast, error) {
	if strings.TrimSpace(airDay) == "" || strings.TrimSpace(animeID) == "" {
		return nil, utils.ErrValidationFailed
	}
	parsedStartsAt, err := time.Parse(time.RFC3339, startsAt)
	if err != nil {
		parsedStartsAt, err = time.Parse("2006-01-02", startsAt)
		if err != nil {
			return nil, utils.ErrValidationFailed
		}
	}
	now := time.Now().UTC()
	simulcast := &models.Simulcast{
		Common:        models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		AnimeID:       animeID,
		AirDay:        airDay,
		AirTime:       airTime,
		AirTimezone:   airTimezone,
		Region:        defaultString(region, "global"),
		Platform:      platform,
		EpisodeNumber: episodeNumber,
		Season:        season,
		StartsAt:      parsedStartsAt,
		Status:        "scheduled",
	}
	if err := s.repos.Simulcasts().Create(ctx, simulcast); err != nil {
		return nil, err
	}
	return simulcast, nil
}

func (s *SchedulingService) UpdateSimulcast(ctx context.Context, userID, id string, req UpdateSimulcastReq) (*models.Simulcast, error) {
	simulcast, err := s.repos.Simulcasts().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.AirDay != nil {
		simulcast.AirDay = *req.AirDay
	}
	if req.AirTime != nil {
		simulcast.AirTime = *req.AirTime
	}
	if req.AirTimezone != nil {
		simulcast.AirTimezone = *req.AirTimezone
	}
	if req.Region != nil {
		simulcast.Region = *req.Region
	}
	if req.Platform != nil {
		simulcast.Platform = *req.Platform
	}
	if req.Status != nil {
		simulcast.Status = *req.Status
	}
	simulcast.UpdatedAt = time.Now().UTC()
	if err := s.repos.Simulcasts().Update(ctx, simulcast); err != nil {
		return nil, err
	}
	return simulcast, nil
}

func (s *SchedulingService) DeleteSimulcast(ctx context.Context, id string) error {
	return s.repos.Simulcasts().Delete(ctx, id)
}

func (s *SchedulingService) ListUpcomingReleases(ctx context.Context, limit int) ([]models.ReleaseSchedule, error) {
	return s.repos.ReleaseSchedules().ListUpcoming(ctx, limit)
}

func (s *SchedulingService) CreateReleaseSchedule(ctx context.Context, userID, animeID string, episodeID *string, scheduledAt, region string) (*models.ReleaseSchedule, error) {
	parsedAt, err := time.Parse(time.RFC3339, scheduledAt)
	if err != nil {
		parsedAt, err = time.Parse("2006-01-02", scheduledAt)
		if err != nil {
			return nil, utils.ErrValidationFailed
		}
	}
	now := time.Now().UTC()
	schedule := &models.ReleaseSchedule{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		AnimeID:     animeID,
		EpisodeID:   episodeID,
		ScheduledAt: parsedAt,
		Region:      defaultString(region, "global"),
		Status:      "scheduled",
	}
	if err := s.repos.ReleaseSchedules().Create(ctx, schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

func (s *SchedulingService) UpdateReleaseSchedule(ctx context.Context, userID, id string, req UpdateReleaseScheduleReq) (*models.ReleaseSchedule, error) {
	schedule, err := s.repos.ReleaseSchedules().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.ScheduledAt != nil {
		parsedAt, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err != nil {
			parsedAt, err = time.Parse("2006-01-02", *req.ScheduledAt)
			if err != nil {
				return nil, utils.ErrValidationFailed
			}
		}
		schedule.ScheduledAt = parsedAt
	}
	if req.Status != nil {
		schedule.Status = *req.Status
	}
	schedule.UpdatedAt = time.Now().UTC()
	if err := s.repos.ReleaseSchedules().Update(ctx, schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

func (s *SchedulingService) GetReleaseSchedule(ctx context.Context, id string) (*models.ReleaseSchedule, error) {
	return s.repos.ReleaseSchedules().GetByID(ctx, id)
}

func (s *SchedulingService) UpdateReleaseScheduleStatus(ctx context.Context, id, status string) (*models.ReleaseSchedule, error) {
	schedule, err := s.repos.ReleaseSchedules().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	schedule.Status = status
	schedule.UpdatedAt = time.Now().UTC()
	if err := s.repos.ReleaseSchedules().Update(ctx, schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

// ListReleases is an alias for ListUpcomingReleases without limit
func (s *SchedulingService) ListReleases(ctx context.Context, limit int) ([]models.ReleaseSchedule, error) {
	return s.repos.ReleaseSchedules().ListUpcoming(ctx, limit)
}
