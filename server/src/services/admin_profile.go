package services

import (
	"context"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
)

type AdminProfileService struct {
	repos *Repositories
}

func NewAdminProfileService(repos *Repositories) *AdminProfileService {
	return &AdminProfileService{repos: repos}
}

type AdminProfileItem struct {
	ID          string  `json:"id"`
	DisplayName string  `json:"displayName"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
	Status      string  `json:"status"`
}

func (s *AdminProfileService) List(ctx context.Context, page, limit int) ([]AdminProfileItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	db := s.repos.db.WithContext(ctx).Model(&models.User{}).Where("status <> ?", "deleted")

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []models.User
	offset := (page - 1) * limit
	if err := db.Order("display_name ASC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	items := make([]AdminProfileItem, 0, len(users))
	for _, u := range users {
		items = append(items, AdminProfileItem{
			ID:          u.ID,
			DisplayName: u.DisplayName,
			AvatarURL:   u.AvatarURL,
			Status:      u.Status,
		})
	}

	return items, total, nil
}

func (s *AdminProfileService) GetByUserID(ctx context.Context, userID string) (*AdminProfileItem, error) {
	user, err := s.repos.Users().GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &AdminProfileItem{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Status:      user.Status,
	}, nil
}

type UpdateAdminProfileInput struct {
	DisplayName *string `json:"displayName"`
	AvatarURL   *string `json:"avatarUrl"`
	Bio         *string `json:"bio"`
}

func (s *AdminProfileService) Update(ctx context.Context, userID string, input UpdateAdminProfileInput) (*AdminProfileItem, error) {
	user, err := s.repos.Users().GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if input.DisplayName != nil {
		user.DisplayName = *input.DisplayName
	}
	if input.AvatarURL != nil {
		user.AvatarURL = input.AvatarURL
	}
	user.UpdatedAt = time.Now().UTC()
	if err := s.repos.Users().Update(ctx, user); err != nil {
		return nil, err
	}
	return &AdminProfileItem{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Status:      user.Status,
	}, nil
}

func (s *AdminProfileService) Delete(ctx context.Context, userID string) error {
	user, err := s.repos.Users().GetByID(ctx, userID)
	if err != nil {
		return err
	}
	user.Status = "deleted"
	user.UpdatedAt = time.Now().UTC()
	return s.repos.Users().Update(ctx, user)
}
