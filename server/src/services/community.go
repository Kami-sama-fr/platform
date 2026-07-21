package services

import (
	"context"
	"strings"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type CommunityService struct {
	repos *Repositories
}

func NewCommunityService(repos *Repositories) *CommunityService {
	return &CommunityService{repos: repos}
}

type UpdateReviewInput struct {
	Rating  *int     `json:"rating"`
	Title   *string  `json:"title"`
	Content *string  `json:"content"`
	Spoiler *bool    `json:"spoiler"`
}

type UpdateCommentInput struct {
	Content *string `json:"content"`
}

type UpdateReportInput struct {
	Status *string `json:"status"`
}

type UpdateWatchlistInput struct {
	Name *string `json:"name"`
	Type *string `json:"type"`
}

func (s *CommunityService) ListReviews(ctx context.Context, animeID string, limit, offset int) ([]models.Review, int64, error) {
	return s.repos.Reviews().ListByAnime(ctx, animeID, limit, offset)
}

func (s *CommunityService) GetReview(ctx context.Context, id string) (*models.Review, error) {
	return s.repos.Reviews().GetByID(ctx, id)
}

func (s *CommunityService) CreateReview(ctx context.Context, userID, animeID string, rating int, title, content string, spoiler bool) (*models.Review, error) {
	if strings.TrimSpace(animeID) == "" || rating < 1 || rating > 10 {
		return nil, utils.ErrValidationFailed
	}
	existing, _ := s.repos.Reviews().GetByUserAndAnime(ctx, userID, animeID)
	if existing != nil {
		return nil, utils.ErrDuplicateReview
	}
	now := time.Now().UTC()
	review := &models.Review{
		Common:  models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		UserID:  userID,
		AnimeID: animeID,
		Rating:  rating,
		Title:   title,
		Content: content,
		Spoiler: spoiler,
	}
	if err := s.repos.Reviews().Create(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *CommunityService) UpdateReview(ctx context.Context, userID, id string, req UpdateReviewInput) (*models.Review, error) {
	review, err := s.repos.Reviews().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if review.UserID != userID {
		return nil, utils.ErrForbidden
	}
	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 10 {
			return nil, utils.ErrValidationFailed
		}
		review.Rating = *req.Rating
	}
	if req.Title != nil {
		review.Title = *req.Title
	}
	if req.Content != nil {
		review.Content = *req.Content
	}
	if req.Spoiler != nil {
		review.Spoiler = *req.Spoiler
	}
	review.UpdatedAt = time.Now().UTC()
	if err := s.repos.Reviews().Update(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *CommunityService) DeleteReview(ctx context.Context, userID, id string) error {
	review, err := s.repos.Reviews().GetByID(ctx, id)
	if err != nil {
		return err
	}
	if review.UserID != userID {
		return utils.ErrForbidden
	}
	return s.repos.Reviews().Delete(ctx, id)
}

func (s *CommunityService) ListComments(ctx context.Context, animeID string, limit, offset int) ([]models.Comment, int64, error) {
	return s.repos.Comments().ListByAnime(ctx, animeID, limit, offset)
}

func (s *CommunityService) ListCommentsByReview(ctx context.Context, reviewID string, limit, offset int) ([]models.Comment, int64, error) {
	return s.repos.Comments().ListByReview(ctx, reviewID, limit, offset)
}

func (s *CommunityService) CreateComment(ctx context.Context, userID, animeID string, reviewID, parentID *string, content string) (*models.Comment, error) {
	if strings.TrimSpace(content) == "" {
		return nil, utils.ErrValidationFailed
	}
	now := time.Now().UTC()
	comment := &models.Comment{
		Common:   models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		UserID:   userID,
		AnimeID:  &animeID,
		ReviewID: reviewID,
		ParentID: parentID,
		Content:  content,
	}
	if err := s.repos.Comments().Create(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommunityService) UpdateComment(ctx context.Context, userID, id, content string) (*models.Comment, error) {
	if strings.TrimSpace(content) == "" {
		return nil, utils.ErrValidationFailed
	}
	comment, err := s.repos.Comments().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if comment.UserID != userID {
		return nil, utils.ErrForbidden
	}
	comment.Content = content
	comment.UpdatedAt = time.Now().UTC()
	if err := s.repos.Comments().Update(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommunityService) DeleteComment(ctx context.Context, userID, id string) error {
	comment, err := s.repos.Comments().GetByID(ctx, id)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return utils.ErrForbidden
	}
	return s.repos.Comments().Delete(ctx, id)
}

func (s *CommunityService) CreateReport(ctx context.Context, userID, targetType, targetID, reason, description string) (*models.Report, error) {
	if strings.TrimSpace(targetType) == "" || strings.TrimSpace(targetID) == "" || strings.TrimSpace(reason) == "" {
		return nil, utils.ErrValidationFailed
	}
	now := time.Now().UTC()
	report := &models.Report{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		UserID:      userID,
		TargetType:  targetType,
		TargetID:    targetID,
		Reason:      reason,
		Description: description,
		Status:      "pending",
	}
	if err := s.repos.Reports().Create(ctx, report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *CommunityService) ListReports(ctx context.Context, status string, limit, offset int) ([]models.Report, int64, error) {
	return s.repos.Reports().List(ctx, status, limit, offset)
}

func (s *CommunityService) UpdateReport(ctx context.Context, userID, id, status string) (*models.Report, error) {
	report, err := s.repos.Reports().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	report.Status = status
	report.ResolvedBy = &userID
	now := time.Now().UTC()
	report.ResolvedAt = &now
	report.UpdatedAt = now
	if err := s.repos.Reports().Update(ctx, report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *CommunityService) ListWatchlists(ctx context.Context, userID string) ([]models.Watchlist, error) {
	return s.repos.Watchlists().ListByUser(ctx, userID)
}

func (s *CommunityService) GetWatchlist(ctx context.Context, id string) (*models.Watchlist, error) {
	return s.repos.Watchlists().GetByID(ctx, id)
}

func (s *CommunityService) CreateWatchlist(ctx context.Context, userID, name, typ string) (*models.Watchlist, error) {
	if strings.TrimSpace(name) == "" {
		return nil, utils.ErrValidationFailed
	}
	now := time.Now().UTC()
	watchlist := &models.Watchlist{
		Common: models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		UserID: userID,
		Name:   strings.TrimSpace(name),
		Type:   typ,
	}
	if err := s.repos.Watchlists().Create(ctx, watchlist); err != nil {
		return nil, err
	}
	return watchlist, nil
}

func (s *CommunityService) UpdateWatchlist(ctx context.Context, userID, id, name, typ string) (*models.Watchlist, error) {
	watchlist, err := s.repos.Watchlists().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if watchlist.UserID != userID {
		return nil, utils.ErrForbidden
	}
	if strings.TrimSpace(name) != "" {
		watchlist.Name = strings.TrimSpace(name)
	}
	if typ != "" {
		watchlist.Type = typ
	}
	watchlist.UpdatedAt = time.Now().UTC()
	if err := s.repos.Watchlists().Update(ctx, watchlist); err != nil {
		return nil, err
	}
	return watchlist, nil
}

func (s *CommunityService) DeleteWatchlist(ctx context.Context, userID, id string) error {
	watchlist, err := s.repos.Watchlists().GetByID(ctx, id)
	if err != nil {
		return err
	}
	if watchlist.UserID != userID {
		return utils.ErrForbidden
	}
	return s.repos.Watchlists().Delete(ctx, id)
}

func (s *CommunityService) AddToWatchlist(ctx context.Context, userID, watchlistID, animeID string) error {
	watchlist, err := s.repos.Watchlists().GetByID(ctx, watchlistID)
	if err != nil {
		return err
	}
	if watchlist.UserID != userID {
		return utils.ErrForbidden
	}
	item := &models.WatchlistItem{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		WatchlistID: watchlistID,
		AnimeID:     animeID,
	}
	return s.repos.Watchlists().AddAnime(ctx, item)
}

func (s *CommunityService) RemoveFromWatchlist(ctx context.Context, userID, watchlistID, animeID string) error {
	watchlist, err := s.repos.Watchlists().GetByID(ctx, watchlistID)
	if err != nil {
		return err
	}
	if watchlist.UserID != userID {
		return utils.ErrForbidden
	}
	return s.repos.Watchlists().RemoveAnime(ctx, watchlistID, animeID)
}

func (s *CommunityService) ListWatchlistAnime(ctx context.Context, id string) ([]models.WatchlistItem, error) {
	return s.repos.Watchlists().ListAnime(ctx, id)
}
