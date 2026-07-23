package services

import (
	"context"
	"time"

	"github.com/skygenesisenterprise/kami-sama/server/src/models"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type SupportService struct {
	repos *Repositories
}

func NewSupportService(repos *Repositories) *SupportService {
	return &SupportService{repos: repos}
}

func (s *SupportService) ListTickets(ctx context.Context, status, priority, category string, page, limit int) ([]models.Ticket, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	return s.repos.Tickets().List(ctx, status, priority, category, limit, offset)
}

func (s *SupportService) CreateTicket(ctx context.Context, userID, subject, description, priority, category string) (*models.Ticket, error) {
	ticket := &models.Ticket{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		UserID:      userID,
		Subject:     subject,
		Description: description,
		Status:      "open",
		Priority:    priority,
		Category:    category,
	}
	if ticket.Priority == "" {
		ticket.Priority = "medium"
	}
	if ticket.Category == "" {
		ticket.Category = "general"
	}
	if err := s.repos.Tickets().Create(ctx, ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

func (s *SupportService) GetTicket(ctx context.Context, id string) (*models.Ticket, error) {
	return s.repos.Tickets().GetByID(ctx, id)
}

func (s *SupportService) UpdateTicket(ctx context.Context, id string, status, assignedTo *string) (*models.Ticket, error) {
	ticket, err := s.repos.Tickets().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if status != nil {
		ticket.Status = *status
		if *status == "resolved" {
			now := time.Now()
			ticket.ResolvedAt = &now
		}
		if *status == "closed" {
			now := time.Now()
			ticket.ClosedAt = &now
		}
	}
	if assignedTo != nil {
		ticket.AssignedTo = assignedTo
	}
	ticket.UpdatedAt = time.Now()
	if err := s.repos.Tickets().Update(ctx, ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

func (s *SupportService) ReplyToTicket(ctx context.Context, ticketID, userID, content string, isStaff bool) (*models.TicketReply, error) {
	ticket, err := s.repos.Tickets().GetByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	_ = ticket
	reply := &models.TicketReply{
		Common:   models.Common{ID: utils.NewID(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		TicketID: ticketID,
		UserID:   userID,
		Content:  content,
		IsStaff:  isStaff,
	}
	if err := s.repos.TicketReplies().Create(ctx, reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (s *SupportService) CloseTicket(ctx context.Context, id string) (*models.Ticket, error) {
	ticket, err := s.repos.Tickets().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	ticket.Status = "closed"
	ticket.ClosedAt = &now
	ticket.UpdatedAt = now
	if err := s.repos.Tickets().Update(ctx, ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

func (s *SupportService) EscalateTicket(ctx context.Context, id string) (*models.Ticket, error) {
	ticket, err := s.repos.Tickets().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	ticket.Priority = "urgent"
	ticket.UpdatedAt = time.Now()
	if err := s.repos.Tickets().Update(ctx, ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

func (s *SupportService) ListTicketReplies(ctx context.Context, ticketID string, page, limit int) ([]models.TicketReply, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit
	return s.repos.TicketReplies().ListByTicket(ctx, ticketID, limit, offset)
}
