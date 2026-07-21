package services

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type LocalMediaSource struct {
	db *gorm.DB
}

func NewLocalMediaSource(db *gorm.DB) *LocalMediaSource {
	return &LocalMediaSource{db: db}
}

func (s *LocalMediaSource) Name() string { return "local" }

func (s *LocalMediaSource) ListLibraries(ctx context.Context) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{
			"id":        "local-default",
			"name":      "Local Library",
			"type":      "tvshows",
			"itemCount": 0,
		},
	}, nil
}

func (s *LocalMediaSource) GetLibrary(ctx context.Context, id string) (map[string]interface{}, error) {
	if id != "local-default" {
		return nil, fmt.Errorf("library %s not found", id)
	}
	return map[string]interface{}{
		"id":        "local-default",
		"name":      "Local Library",
		"type":      "tvshows",
		"itemCount": 0,
	}, nil
}

func (s *LocalMediaSource) ListItems(ctx context.Context, libraryID string, limit, offset int, sortBy, query string) ([]map[string]interface{}, int, error) {
	type animeRecord struct {
		ID            string
		Slug          string
		Title         string
		JapaneseTitle string
		Synopsis      string
		CoverImageUrl string
		Status        string
		Rating        float64
		TotalEpisodes int
		ReleaseYear   int
		Season        string
		Source        string
		AgeRating     string
		IsFeatured    bool
		IsTrending    bool
	}

	var records []animeRecord
	var total int64

	tx := s.db.Model(&struct{}{}).Table("anime")
	if query != "" {
		tx = tx.Where("title ILIKE ? OR japanese_title ILIKE ?", "%"+query+"%", "%"+query+"%")
	}
	tx.Count(&total)
	if err := tx.Offset(offset).Limit(limit).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	items := make([]map[string]interface{}, 0, len(records))
	for _, r := range records {
		items = append(items, map[string]interface{}{
			"id":            r.ID,
			"sourceId":      r.ID,
			"name":          r.Title,
			"originalTitle": r.JapaneseTitle,
			"type":          "Series",
			"year":          r.ReleaseYear,
			"rating":        r.Rating,
			"overview":      r.Synopsis,
			"genres":        []string{},
		})
	}
	return items, int(total), nil
}

func (s *LocalMediaSource) GetItem(ctx context.Context, id string) (map[string]interface{}, error) {
	type animeRecord struct {
		ID            string
		Title         string
		JapaneseTitle string
		Synopsis      string
		CoverImageUrl string
		Rating        float64
		TotalEpisodes int
		ReleaseYear   int
	}
	var r animeRecord
	if err := s.db.Table("anime").Where("id = ?", id).First(&r).Error; err != nil {
		return nil, fmt.Errorf("item %s not found", id)
	}
	return map[string]interface{}{
		"id":            r.ID,
		"sourceId":      r.ID,
		"name":          r.Title,
		"originalTitle": r.JapaneseTitle,
		"type":          "Series",
		"year":          r.ReleaseYear,
		"rating":        r.Rating,
		"overview":      r.Synopsis,
	}, nil
}

func (s *LocalMediaSource) SearchItems(ctx context.Context, query string, limit int) ([]map[string]interface{}, error) {
	items, _, err := s.ListItems(ctx, "local-default", limit, 0, "", query)
	return items, err
}

func (s *LocalMediaSource) GetStreamURL(ctx context.Context, itemID string, static bool) (string, error) {
	return "", fmt.Errorf("local source does not support streaming")
}

func (s *LocalMediaSource) GetPlaybackInfo(ctx context.Context, itemID string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("local source does not support playback info")
}

func (s *LocalMediaSource) ReportPlaybackProgress(ctx context.Context, itemID string, positionTicks int64, stopped bool) error {
	return nil
}

func (s *LocalMediaSource) SyncLibrary(ctx context.Context, libraryID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"libraryId":    libraryID,
		"itemsCreated": 0,
		"itemsUpdated": 0,
		"itemsRemoved": 0,
		"startedAt":    time.Now().UTC(),
		"completedAt":  time.Now().UTC(),
	}, nil
}

func (s *LocalMediaSource) GetSyncStatus(ctx context.Context, libraryID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"libraryId":  libraryID,
		"lastSyncAt": nil,
		"status":     "idle",
		"itemCount":  0,
	}, nil
}
