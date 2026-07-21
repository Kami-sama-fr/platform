package services

import (
	"context"
	"strings"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
	"github.com/kami-sama-fr/platform/server/src/utils"
)

type MediaService struct {
	repos *Repositories
}

func NewMediaService(repos *Repositories) *MediaService {
	return &MediaService{repos: repos}
}

type UpdateMediaAssetInput struct {
	Type     *string  `json:"type"`
	Status   *string  `json:"status"`
	CdnUrl   *string  `json:"cdnUrl"`
	Duration *float64 `json:"duration"`
	Quality  *string  `json:"quality"`
	Codec    *string  `json:"codec"`
}

func (s *MediaService) List(ctx context.Context, mediaType string, limit, offset int) ([]models.MediaAsset, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repos.MediaAssets().List(ctx, mediaType, limit, offset)
}

func (s *MediaService) GetByID(ctx context.Context, id string) (*models.MediaAsset, error) {
	return s.repos.MediaAssets().GetByID(ctx, id)
}

func (s *MediaService) Create(ctx context.Context, userID string, animeID string, episodeID *string, typ, filename, originalUrl, mimeType string, size int64, duration float64, width, height int, language, quality, codec string) (*models.MediaAsset, error) {
	if strings.TrimSpace(filename) == "" || strings.TrimSpace(animeID) == "" {
		return nil, utils.ErrValidationFailed
	}
	now := time.Now().UTC()
	asset := &models.MediaAsset{
		Common:      models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		AnimeId:     animeID,
		EpisodeId:   episodeID,
		Type:        typ,
		Filename:    filename,
		OriginalUrl: originalUrl,
		MimeType:    mimeType,
		Size:        size,
		Duration:    duration,
		Width:       width,
		Height:      height,
		Language:    language,
		Quality:     quality,
		Codec:       codec,
		Status:      "pending",
	}
	if err := s.repos.MediaAssets().Create(ctx, asset); err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *MediaService) Update(ctx context.Context, userID, id string, req UpdateMediaAssetInput) (*models.MediaAsset, error) {
	asset, err := s.repos.MediaAssets().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Type != nil {
		asset.Type = *req.Type
	}
	if req.Status != nil {
		asset.Status = *req.Status
	}
	if req.CdnUrl != nil {
		asset.CdnUrl = *req.CdnUrl
	}
	if req.Duration != nil {
		asset.Duration = *req.Duration
	}
	if req.Quality != nil {
		asset.Quality = *req.Quality
	}
	if req.Codec != nil {
		asset.Codec = *req.Codec
	}
	asset.UpdatedAt = time.Now().UTC()
	if err := s.repos.MediaAssets().Update(ctx, asset); err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *MediaService) Delete(ctx context.Context, id string) error {
	return s.repos.MediaAssets().Delete(ctx, id)
}

func (s *MediaService) ListEncodingJobs(ctx context.Context, status string, limit, offset int) ([]models.EncodingJob, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repos.EncodingJobs().List(ctx, status, limit, offset)
}

func (s *MediaService) GetEncodingJob(ctx context.Context, id string) (*models.EncodingJob, error) {
	return s.repos.EncodingJobs().GetByID(ctx, id)
}

func (s *MediaService) TriggerEncoding(ctx context.Context, mediaAssetID string, profile string) (*models.EncodingJob, error) {
	_, err := s.repos.MediaAssets().GetByID(ctx, mediaAssetID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	job := &models.EncodingJob{
		Common:       models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		MediaAssetId: mediaAssetID,
		Status:       "pending",
		Profile:      profile,
	}
	if err := s.repos.EncodingJobs().Create(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}
