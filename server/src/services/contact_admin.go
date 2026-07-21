package services

import (
	"context"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
)

type ContactAdminService struct {
	repos *Repositories
}

func NewContactAdminService(repos *Repositories) *ContactAdminService {
	return &ContactAdminService{repos: repos}
}

func (s *ContactAdminService) List(ctx context.Context, status string, page, limit int) ([]models.ContactMessage, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	return s.repos.ContactMessages().List(ctx, status, limit, offset)
}

func (s *ContactAdminService) GetByID(ctx context.Context, id string) (*models.ContactMessage, error) {
	return s.repos.ContactMessages().GetByID(ctx, id)
}

func (s *ContactAdminService) MarkAsRead(ctx context.Context, id string) (*models.ContactMessage, error) {
	msg, err := s.repos.ContactMessages().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	msg.Status = "read"
	msg.ReadAt = &now
	msg.UpdatedAt = now
	if err := s.repos.ContactMessages().Update(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *ContactAdminService) MarkAsReplied(ctx context.Context, id string) (*models.ContactMessage, error) {
	msg, err := s.repos.ContactMessages().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	msg.Status = "replied"
	msg.RepliedAt = &now
	msg.UpdatedAt = now
	if err := s.repos.ContactMessages().Update(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *ContactAdminService) Reply(ctx context.Context, id string) error {
	msg, err := s.repos.ContactMessages().GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	msg.Status = "replied"
	msg.RepliedAt = &now
	msg.UpdatedAt = now
	return s.repos.ContactMessages().Update(ctx, msg)
}

func (s *ContactAdminService) Delete(ctx context.Context, id string) error {
	return s.repos.ContactMessages().Delete(ctx, id)
}
