package services

import (
	"context"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
	"gorm.io/gorm"
)

type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

type DashboardStats struct {
	ActiveSubscribers     int64   `json:"activeSubscribers"`
	ActiveSubscribersDelta int64  `json:"activeSubscribersDelta"`
	DailyViews            int64   `json:"dailyViews"`
	DailyViewsDelta       int64   `json:"dailyViewsDelta"`
	MonthlyRevenue        float64 `json:"monthlyRevenue"`
	MonthlyRevenueDelta   float64 `json:"monthlyRevenueDelta"`
	AnimeCount            int64   `json:"animeCount"`
	AnimeCountDelta       int64   `json:"animeCountDelta"`
}

type WeeklyViewsEntry struct {
	Date  string `json:"date"`
	Views int64  `json:"views"`
}

type SubscriptionDistributionEntry struct {
	Plan        string  `json:"plan"`
	MemberCount int64   `json:"memberCount"`
	Percentage  float64 `json:"percentage"`
}

type TopAnimeEntry struct {
	AnimeID   string  `json:"animeId"`
	Title     string  `json:"title"`
	Views     int64   `json:"views"`
	Rating    float64 `json:"rating"`
}

type RecentUploadEntry struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Episode     string    `json:"episode"`
	Language    string    `json:"language"`
	Status      string    `json:"status"`
	UploadedAt  time.Time `json:"uploadedAt"`
}

func (s *DashboardService) GetStats(ctx context.Context) (*DashboardStats, error) {
	now := time.Now().UTC()
	todayStart := now.Truncate(24 * time.Hour)
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastMonthStart := monthStart.AddDate(0, -1, 0)

	var stats DashboardStats

	// Active subscribers (status = 'active')
	err := s.db.WithContext(ctx).Model(&models.Subscription{}).Where("status = ?", "active").Count(&stats.ActiveSubscribers).Error
	if err != nil {
		return nil, err
	}
	var lastMonthSubscribers int64
	err = s.db.WithContext(ctx).Model(&models.Subscription{}).Where("status = ? AND start_date < ?", "active", monthStart).Count(&lastMonthSubscribers).Error
	if err != nil {
		return nil, err
	}
	stats.ActiveSubscribersDelta = stats.ActiveSubscribers - lastMonthSubscribers

	// Daily views (today vs yesterday)
	err = s.db.WithContext(ctx).Model(&models.WatchHistory{}).Where("watched_at >= ?", todayStart).Count(&stats.DailyViews).Error
	if err != nil {
		return nil, err
	}
	var yesterdayViews int64
	err = s.db.WithContext(ctx).Model(&models.WatchHistory{}).Where("watched_at >= ? AND watched_at < ?", yesterdayStart, todayStart).Count(&yesterdayViews).Error
	if err != nil {
		return nil, err
	}
	stats.DailyViewsDelta = stats.DailyViews - yesterdayViews

	// Monthly revenue (sum of active subscription amounts for current month)
	type revenueResult struct {
		Total float64
	}
	var currentRevenue revenueResult
	s.db.WithContext(ctx).Model(&models.Subscription{}).Select("COALESCE(SUM(amount), 0) as total").Where("start_date >= ?", monthStart).Scan(&currentRevenue)
	stats.MonthlyRevenue = currentRevenue.Total

	var lastMonthRevenue revenueResult
	s.db.WithContext(ctx).Model(&models.Subscription{}).Select("COALESCE(SUM(amount), 0) as total").Where("start_date >= ? AND start_date < ?", lastMonthStart, monthStart).Scan(&lastMonthRevenue)
	stats.MonthlyRevenueDelta = currentRevenue.Total - lastMonthRevenue.Total

	// Anime count
	err = s.db.WithContext(ctx).Model(&models.Anime{}).Where("deleted_at IS NULL").Count(&stats.AnimeCount).Error
	if err != nil {
		return nil, err
	}
	var lastMonthAnimeCount int64
	err = s.db.WithContext(ctx).Model(&models.Anime{}).Where("deleted_at IS NULL AND created_at < ?", monthStart).Count(&lastMonthAnimeCount).Error
	if err != nil {
		return nil, err
	}
	stats.AnimeCountDelta = stats.AnimeCount - lastMonthAnimeCount

	return &stats, nil
}

func (s *DashboardService) GetWeeklyViews(ctx context.Context) ([]WeeklyViewsEntry, error) {
	now := time.Now().UTC()
	weekAgo := now.AddDate(0, 0, -7)

	type dayViews struct {
		Date  string `gorm:"date"`
		Views int64  `gorm:"views"`
	}
	var results []dayViews
	err := s.db.WithContext(ctx).
		Model(&models.WatchHistory{}).
		Select("DATE(watched_at) as date, COUNT(*) as views").
		Where("watched_at >= ?", weekAgo).
		Group("DATE(watched_at)").
		Order("date ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	entries := make([]WeeklyViewsEntry, len(results))
	for i, r := range results {
		entries[i] = WeeklyViewsEntry{Date: r.Date, Views: r.Views}
	}
	return entries, nil
}

func (s *DashboardService) GetSubscriptionDistribution(ctx context.Context) ([]SubscriptionDistributionEntry, error) {
	type planResult struct {
		Plan        string `gorm:"plan_name"`
		MemberCount int64  `gorm:"member_count"`
	}
	var results []planResult
	err := s.db.WithContext(ctx).
		Model(&models.Subscription{}).
		Select("plan_name, COUNT(*) as member_count").
		Where("status = ?", "active").
		Group("plan_name").
		Order("member_count DESC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var total int64
	for _, r := range results {
		total += r.MemberCount
	}

	entries := make([]SubscriptionDistributionEntry, len(results))
	for i, r := range results {
		pct := 0.0
		if total > 0 {
			pct = float64(r.MemberCount) / float64(total) * 100
		}
		entries[i] = SubscriptionDistributionEntry{
			Plan:        r.Plan,
			MemberCount: r.MemberCount,
			Percentage:  pct,
		}
	}
	return entries, nil
}

func (s *DashboardService) GetTopAnime(ctx context.Context, limit int) ([]TopAnimeEntry, error) {
	type topResult struct {
		AnimeID string  `gorm:"anime_id"`
		Views   int64   `gorm:"views"`
	}
	var results []topResult
	err := s.db.WithContext(ctx).
		Model(&models.WatchHistory{}).
		Select("anime_id, COUNT(*) as views").
		Group("anime_id").
		Order("views DESC").
		Limit(limit).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []TopAnimeEntry{}, nil
	}

	entries := make([]TopAnimeEntry, len(results))
	for i, r := range results {
		var anime models.Anime
		s.db.WithContext(ctx).First(&anime, "id = ?", r.AnimeID)
		entries[i] = TopAnimeEntry{
			AnimeID: r.AnimeID,
			Title:   anime.Title,
			Views:   r.Views,
			Rating:  anime.Rating,
		}
	}
	return entries, nil
}

func (s *DashboardService) GetRecentUploads(ctx context.Context, limit int) ([]RecentUploadEntry, error) {
	type uploadResult struct {
		ID         string    `gorm:"id"`
		Title      string    `gorm:"title"`
		Episode    string    `gorm:"episode"`
		Language   string    `gorm:"language"`
		Status     string    `gorm:"status"`
		UploadedAt time.Time `gorm:"uploaded_at"`
	}
	var results []uploadResult
	err := s.db.WithContext(ctx).
		Model(&models.MediaAsset{}).
		Select("id, filename as title, episode_id as episode, language, status, created_at as uploaded_at").
		Order("created_at DESC").
		Limit(limit).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	entries := make([]RecentUploadEntry, len(results))
	for i, r := range results {
		entries[i] = RecentUploadEntry{
			ID:         r.ID,
			Title:      r.Title,
			Episode:    r.Episode,
			Language:   r.Language,
			Status:     r.Status,
			UploadedAt: r.UploadedAt,
		}
	}
	return entries, nil
}
