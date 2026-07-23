package services

import (
	"context"
	"time"

	"github.com/skygenesisenterprise/kami-sama/server/src/models"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type FAQService struct {
	repos *Repositories
}

func NewFAQService(repos *Repositories) *FAQService {
	return &FAQService{repos: repos}
}

func (s *FAQService) List(ctx context.Context, category string, activeOnly bool) ([]models.FAQ, error) {
	return s.repos.FAQs().List(ctx, category, activeOnly)
}

func (s *FAQService) Create(ctx context.Context, question, answer, category string, sortOrder int) (*models.FAQ, error) {
	faq := &models.FAQ{
		Common:    models.Common{ID: utils.NewID(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Question:  question,
		Answer:    answer,
		Category:  category,
		SortOrder: sortOrder,
		IsActive:  true,
	}
	if faq.Category == "" {
		faq.Category = "general"
	}
	if err := s.repos.FAQs().Create(ctx, faq); err != nil {
		return nil, err
	}
	return faq, nil
}

func (s *FAQService) GetByID(ctx context.Context, id string) (*models.FAQ, error) {
	return s.repos.FAQs().GetByID(ctx, id)
}

func (s *FAQService) Update(ctx context.Context, id string, question, answer, category *string, sortOrder *int, isActive *bool) (*models.FAQ, error) {
	faq, err := s.repos.FAQs().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if question != nil {
		faq.Question = *question
	}
	if answer != nil {
		faq.Answer = *answer
	}
	if category != nil {
		faq.Category = *category
	}
	if sortOrder != nil {
		faq.SortOrder = *sortOrder
	}
	if isActive != nil {
		faq.IsActive = *isActive
	}
	faq.UpdatedAt = time.Now()
	if err := s.repos.FAQs().Update(ctx, faq); err != nil {
		return nil, err
	}
	return faq, nil
}

func (s *FAQService) Delete(ctx context.Context, id string) error {
	return s.repos.FAQs().Delete(ctx, id)
}

func (s *FAQService) Reorder(ctx context.Context, orders map[string]int) error {
	return s.repos.FAQs().UpdateSortOrders(ctx, orders)
}
