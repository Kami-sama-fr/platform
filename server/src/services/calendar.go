package services

import (
	"context"
	"strings"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
	"github.com/kami-sama-fr/platform/server/src/utils"
	"gorm.io/datatypes"
)

type CalendarService struct {
	repos *Repositories
}

func NewCalendarService(repos *Repositories) *CalendarService {
	return &CalendarService{repos: repos}
}

type UpdateCalendarEventInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	EventType   *string `json:"eventType"`
	AnimeID     *string `json:"animeId"`
	EpisodeID   *string `json:"episodeId"`
	StartAt     *string `json:"startAt"`
	EndAt       *string `json:"endAt"`
	AllDay      *bool   `json:"allDay"`
	Status      *string `json:"status"`
	Color       *string `json:"color"`
}

func (s *CalendarService) ListEvents(ctx context.Context, startDate, endDate string) ([]models.CalendarEvent, error) {
	return s.repos.CalendarEvents().ListByDateRange(ctx, startDate, endDate)
}

func (s *CalendarService) ListEventsAll(ctx context.Context) ([]models.CalendarEvent, error) {
	return s.repos.CalendarEvents().List(ctx)
}

func (s *CalendarService) GetEvent(ctx context.Context, id string) (*models.CalendarEvent, error) {
	return s.repos.CalendarEvents().GetByID(ctx, id)
}

func (s *CalendarService) CreateEvent(ctx context.Context, userID string, title, description, eventType string, animeID, episodeID *string, startAt, endAt string, allDay bool, color string) (*models.CalendarEvent, error) {
	if strings.TrimSpace(title) == "" {
		return nil, utils.ErrValidationFailed
	}

	parsedStart, err := time.Parse(time.RFC3339, startAt)
	if err != nil {
		parsedStart, err = time.Parse("2006-01-02T15:04:05", startAt)
		if err != nil {
			parsedStart, err = time.Parse("2006-01-02", startAt)
			if err != nil {
				return nil, utils.ErrValidationFailed
			}
		}
	}

	var parsedEnd *time.Time
	if endAt != "" {
		t, err := time.Parse(time.RFC3339, endAt)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05", endAt)
			if err != nil {
				t, err = time.Parse("2006-01-02", endAt)
				if err != nil {
					return nil, utils.ErrValidationFailed
				}
			}
		}
		parsedEnd = &t
	}

	now := time.Now().UTC()
	event := &models.CalendarEvent{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		Title:       title,
		Description: description,
		EventType:   defaultString(eventType, "release"),
		AnimeID:     animeID,
		EpisodeID:   episodeID,
		StartAt:     parsedStart,
		EndAt:       parsedEnd,
		AllDay:      allDay,
		Status:      "scheduled",
		Color:       color,
		Metadata:    datatypes.JSON("{}"),
	}
	if err := s.repos.CalendarEvents().Create(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *CalendarService) UpdateEvent(ctx context.Context, userID, id string, req UpdateCalendarEventInput) (*models.CalendarEvent, error) {
	event, err := s.repos.CalendarEvents().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.EventType != nil {
		event.EventType = *req.EventType
	}
	if req.AnimeID != nil {
		event.AnimeID = req.AnimeID
	}
	if req.EpisodeID != nil {
		event.EpisodeID = req.EpisodeID
	}
	if req.StartAt != nil {
		t, err := time.Parse(time.RFC3339, *req.StartAt)
		if err != nil {
			return nil, utils.ErrValidationFailed
		}
		event.StartAt = t
	}
	if req.EndAt != nil {
		if *req.EndAt == "" {
			event.EndAt = nil
		} else {
			t, err := time.Parse(time.RFC3339, *req.EndAt)
			if err != nil {
				return nil, utils.ErrValidationFailed
			}
			event.EndAt = &t
		}
	}
	if req.AllDay != nil {
		event.AllDay = *req.AllDay
	}
	if req.Color != nil {
		event.Color = *req.Color
	}
	if req.Status != nil {
		event.Status = *req.Status
	}
	event.UpdatedAt = time.Now().UTC()
	if err := s.repos.CalendarEvents().Update(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *CalendarService) DeleteEvent(ctx context.Context, id string) error {
	return s.repos.CalendarEvents().Delete(ctx, id)
}
