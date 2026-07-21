package services

import (
	"context"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
	"gorm.io/gorm"
)

type AnalyticsService struct {
	db *gorm.DB
}

func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

type AnalyticsOverview struct {
	TotalViews       int64   `json:"totalViews"`
	ActiveUsers      int64   `json:"activeUsers"`
	AvgWatchTime     float64 `json:"avgWatchTime"`
	CompletionRate   float64 `json:"completionRate"`
}

type PeriodMetrics struct {
	Date        string  `json:"date"`
	TotalViews  int64   `json:"totalViews"`
	ActiveUsers int64   `json:"activeUsers"`
	AvgWatchTime float64 `json:"avgWatchTime"`
	Trend       float64 `json:"trend"`
}

type WatchTimeEntry struct {
	Date         string  `json:"date"`
	AvgWatchTime float64 `json:"avgWatchTime"`
	TotalTime    float64 `json:"totalTime"`
	SessionCount int64   `json:"sessionCount"`
}

type WatchTimeByAnimeEntry struct {
	AnimeID      string  `json:"animeId"`
	Title        string  `json:"title"`
	AvgWatchTime float64 `json:"avgWatchTime"`
	TotalViews   int64   `json:"totalViews"`
}

type WatchTimeByEpisodeEntry struct {
	AnimeID      string  `json:"animeId"`
	EpisodeID    string  `json:"episodeId"`
	EpisodeTitle string  `json:"episodeTitle"`
	AvgWatchTime float64 `json:"avgWatchTime"`
	TotalViews   int64   `json:"totalViews"`
}

type HistogramBucket struct {
	Range      string `json:"range"`
	Count      int64  `json:"count"`
	Percentage float64 `json:"percentage"`
}

type DeviceEntry struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
	Pct   float64 `json:"percentage"`
}

type BrowserEntry struct {
	Browser string `json:"browser"`
	Count   int64  `json:"count"`
	Pct     float64 `json:"percentage"`
}

type OSEntry struct {
	OS    string `json:"os"`
	Count int64  `json:"count"`
	Pct   float64 `json:"percentage"`
}

type PopularEntry struct {
	ContentID    string  `json:"contentId"`
	Title        string  `json:"title"`
	Views        int64   `json:"views"`
	AvgWatchTime float64 `json:"avgWatchTime"`
}

type TrendingEntry struct {
	ContentID string  `json:"contentId"`
	Title     string  `json:"title"`
	Views     int64   `json:"views"`
	Growth    float64 `json:"growth"`
}

type NewPopularEntry struct {
	ContentID   string `json:"contentId"`
	Title       string `json:"title"`
	Views       int64  `json:"views"`
	UploadedAt  string `json:"uploadedAt"`
}

type GeographyEntry struct {
	Country string `json:"country"`
	Views   int64  `json:"views"`
	Pct     float64 `json:"percentage"`
}

type ActiveUsersEntry struct {
	Date   string `json:"date"`
	Count  int64  `json:"count"`
}

type RetentionEntry struct {
	Period    string  `json:"period"`
	Rate      float64 `json:"rate"`
	UserCount int64   `json:"userCount"`
}

type SessionEntry struct {
	Period       string  `json:"period"`
	AvgSessions  float64 `json:"avgSessions"`
	TotalUsers   int64   `json:"totalUsers"`
}

func (s *AnalyticsService) GetOverview(ctx context.Context) (*AnalyticsOverview, error) {
	var overview AnalyticsOverview

	s.db.WithContext(ctx).Model(&models.WatchHistory{}).Count(&overview.TotalViews)

	var uniqueUsers int64
	s.db.WithContext(ctx).Model(&models.WatchHistory{}).Select("COUNT(DISTINCT user_id)").Scan(&uniqueUsers)
	overview.ActiveUsers = uniqueUsers

	s.db.WithContext(ctx).Model(&models.WatchHistory{}).Select("COALESCE(AVG(duration), 0)").Scan(&overview.AvgWatchTime)

	var completed, total int64
	s.db.WithContext(ctx).Model(&models.WatchProgress{}).Count(&total)
	s.db.WithContext(ctx).Model(&models.WatchProgress{}).Where("completed = ?", true).Count(&completed)
	overview.CompletionRate = 0
	if total > 0 {
		overview.CompletionRate = float64(completed) / float64(total) * 100
	}

	return &overview, nil
}

func (s *AnalyticsService) GetOverviewByPeriod(ctx context.Context, period string) ([]PeriodMetrics, error) {
	now := time.Now().UTC()
	var startDate time.Time
	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	default:
		startDate = now.AddDate(0, 0, -30)
	}

	type periodResult struct {
		Date        string  `gorm:"date"`
		TotalViews  int64   `gorm:"total_views"`
		ActiveUsers int64   `gorm:"active_users"`
		AvgWatchTime float64 `gorm:"avg_watch_time"`
	}
	var results []periodResult
	err := s.db.WithContext(ctx).
		Model(&models.WatchHistory{}).
		Select("DATE(watched_at) as date, COUNT(*) as total_views, COUNT(DISTINCT user_id) as active_users, COALESCE(AVG(duration), 0) as avg_watch_time").
		Where("watched_at >= ?", startDate).
		Group("DATE(watched_at)").
		Order("date ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	entries := make([]PeriodMetrics, len(results))
	for i, r := range results {
		trend := 0.0
		if i > 0 && entries[i-1].TotalViews > 0 {
			trend = float64(r.TotalViews-entries[i-1].TotalViews) / float64(entries[i-1].TotalViews) * 100
		}
		entries[i] = PeriodMetrics{
			Date:         r.Date,
			TotalViews:   r.TotalViews,
			ActiveUsers:  r.ActiveUsers,
			AvgWatchTime: r.AvgWatchTime,
			Trend:        trend,
		}
	}
	return entries, nil
}

func (s *AnalyticsService) GetWatchTime(ctx context.Context) ([]WatchTimeEntry, error) {
	now := time.Now().UTC()
	weekAgo := now.AddDate(0, 0, -7)

	type wtResult struct {
		Date         string  `gorm:"date"`
		AvgWatchTime float64 `gorm:"avg_watch_time"`
		TotalTime    float64 `gorm:"total_time"`
		SessionCount int64   `gorm:"session_count"`
	}
	var results []wtResult
	err := s.db.WithContext(ctx).
		Model(&models.WatchHistory{}).
		Select("DATE(watched_at) as date, COALESCE(AVG(duration), 0) as avg_watch_time, COALESCE(SUM(duration), 0) as total_time, COUNT(*) as session_count").
		Where("watched_at >= ?", weekAgo).
		Group("DATE(watched_at)").
		Order("date ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	entries := make([]WatchTimeEntry, len(results))
	for i, r := range results {
		entries[i] = WatchTimeEntry{
			Date:         r.Date,
			AvgWatchTime: r.AvgWatchTime,
			TotalTime:    r.TotalTime,
			SessionCount: r.SessionCount,
		}
	}
	return entries, nil
}

func (s *AnalyticsService) GetWatchTimeByAnime(ctx context.Context) ([]WatchTimeByAnimeEntry, error) {
	type wtAnime struct {
		AnimeID      string  `gorm:"anime_id"`
		AvgWatchTime float64 `gorm:"avg_watch_time"`
		TotalViews   int64   `gorm:"total_views"`
	}
	var results []wtAnime
	err := s.db.WithContext(ctx).
		Model(&models.WatchHistory{}).
		Select("anime_id, COALESCE(AVG(duration), 0) as avg_watch_time, COUNT(*) as total_views").
		Group("anime_id").
		Order("total_views DESC").
		Limit(50).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	entries := make([]WatchTimeByAnimeEntry, len(results))
	for i, r := range results {
		var anime models.Anime
		s.db.WithContext(ctx).First(&anime, "id = ?", r.AnimeID)
		entries[i] = WatchTimeByAnimeEntry{
			AnimeID:      r.AnimeID,
			Title:        anime.Title,
			AvgWatchTime: r.AvgWatchTime,
			TotalViews:   r.TotalViews,
		}
	}
	return entries, nil
}

func (s *AnalyticsService) GetWatchTimeByEpisode(ctx context.Context) ([]WatchTimeByEpisodeEntry, error) {
	type wtEpisode struct {
		AnimeID      string  `gorm:"anime_id"`
		EpisodeID    string  `gorm:"episode_id"`
		AvgWatchTime float64 `gorm:"avg_watch_time"`
		TotalViews   int64   `gorm:"total_views"`
	}
	var results []wtEpisode
	err := s.db.WithContext(ctx).
		Model(&models.WatchHistory{}).
		Select("anime_id, episode_id, COALESCE(AVG(duration), 0) as avg_watch_time, COUNT(*) as total_views").
		Group("anime_id, episode_id").
		Order("total_views DESC").
		Limit(50).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	entries := make([]WatchTimeByEpisodeEntry, len(results))
	for i, r := range results {
		var ep models.Episode
		s.db.WithContext(ctx).First(&ep, "id = ?", r.EpisodeID)
		entries[i] = WatchTimeByEpisodeEntry{
			AnimeID:      r.AnimeID,
			EpisodeID:    r.EpisodeID,
			EpisodeTitle: ep.Title,
			AvgWatchTime: r.AvgWatchTime,
			TotalViews:   r.TotalViews,
		}
	}
	return entries, nil
}

func (s *AnalyticsService) GetWatchTimeHistogram(ctx context.Context) ([]HistogramBucket, error) {
	type histResult struct {
		Range string `gorm:"range"`
		Count int64  `gorm:"count"`
	}

	query := `
		SELECT 
			CASE
				WHEN duration < 300 THEN '0-5 min'
				WHEN duration < 600 THEN '5-10 min'
				WHEN duration < 1200 THEN '10-20 min'
				WHEN duration < 1800 THEN '20-30 min'
				WHEN duration < 3600 THEN '30-60 min'
				ELSE '60+ min'
			END as range,
			COUNT(*) as count
		FROM watch_history
		GROUP BY range
		ORDER BY MIN(duration)
	`
	var results []histResult
	if err := s.db.WithContext(ctx).Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}

	var total int64
	for _, r := range results {
		total += r.Count
	}

	buckets := make([]HistogramBucket, len(results))
	for i, r := range results {
		pct := 0.0
		if total > 0 {
			pct = float64(r.Count) / float64(total) * 100
		}
		buckets[i] = HistogramBucket{
			Range:      r.Range,
			Count:      r.Count,
			Percentage: pct,
		}
	}
	return buckets, nil
}

func (s *AnalyticsService) GetDevices(ctx context.Context) ([]DeviceEntry, error) {
	type devResult struct {
		Type  string `gorm:"type"`
		Count int64  `gorm:"count"`
	}
	var results []devResult
	err := s.db.WithContext(ctx).
		Model(&models.Device{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Order("count DESC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var total int64
	for _, r := range results {
		total += r.Count
	}

	entries := make([]DeviceEntry, len(results))
	for i, r := range results {
		pct := 0.0
		if total > 0 {
			pct = float64(r.Count) / float64(total) * 100
		}
		entries[i] = DeviceEntry{Type: r.Type, Count: r.Count, Pct: pct}
	}
	return entries, nil
}

func (s *AnalyticsService) GetDevicesBrowsers(ctx context.Context) ([]BrowserEntry, error) {
	type brResult struct {
		Browser string `gorm:"browser"`
		Count   int64  `gorm:"count"`
	}
	var results []brResult
	err := s.db.WithContext(ctx).
		Model(&models.Device{}).
		Select("COALESCE(browser, 'unknown') as browser, COUNT(*) as count").
		Where("browser IS NOT NULL AND browser != ''").
		Group("browser").
		Order("count DESC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var total int64
	for _, r := range results {
		total += r.Count
	}

	entries := make([]BrowserEntry, len(results))
	for i, r := range results {
		pct := 0.0
		if total > 0 {
			pct = float64(r.Count) / float64(total) * 100
		}
		entries[i] = BrowserEntry{Browser: r.Browser, Count: r.Count, Pct: pct}
	}
	return entries, nil
}

func (s *AnalyticsService) GetDevicesOS(ctx context.Context) ([]OSEntry, error) {
	type osResult struct {
		OS    string `gorm:"os"`
		Count int64  `gorm:"count"`
	}
	var results []osResult
	err := s.db.WithContext(ctx).
		Model(&models.Device{}).
		Select("COALESCE(os, 'unknown') as os, COUNT(*) as count").
		Where("os IS NOT NULL AND os != ''").
		Group("os").
		Order("count DESC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var total int64
	for _, r := range results {
		total += r.Count
	}

	entries := make([]OSEntry, len(results))
	for i, r := range results {
		pct := 0.0
		if total > 0 {
			pct = float64(r.Count) / float64(total) * 100
		}
		entries[i] = OSEntry{OS: r.OS, Count: r.Count, Pct: pct}
	}
	return entries, nil
}

func (s *AnalyticsService) GetPopular(ctx context.Context, limit int) ([]PopularEntry, error) {
	type popResult struct {
		ContentID    string  `gorm:"content_id"`
		Title        string  `gorm:"title"`
		Views        int64   `gorm:"views"`
		AvgWatchTime float64 `gorm:"avg_watch_time"`
	}
	var results []popResult
	err := s.db.WithContext(ctx).
		Model(&models.WatchHistory{}).
		Select("anime_id as content_id, '' as title, COUNT(*) as views, COALESCE(AVG(duration), 0) as avg_watch_time").
		Group("anime_id").
		Order("views DESC").
		Limit(limit).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	entries := make([]PopularEntry, len(results))
	for i, r := range results {
		var anime models.Anime
		s.db.WithContext(ctx).First(&anime, "id = ?", r.ContentID)
		entries[i] = PopularEntry{
			ContentID:    r.ContentID,
			Title:        anime.Title,
			Views:        r.Views,
			AvgWatchTime: r.AvgWatchTime,
		}
	}
	return entries, nil
}

func (s *AnalyticsService) GetPopularTrending(ctx context.Context) ([]TrendingEntry, error) {
	now := time.Now().UTC()
	weekAgo := now.AddDate(0, 0, -7)
	twoWeeksAgo := now.AddDate(0, 0, -14)

	type trendResult struct {
		ContentID string  `gorm:"content_id"`
		Views     int64   `gorm:"views"`
	}
	var thisWeek []trendResult
	err := s.db.WithContext(ctx).
		Model(&models.WatchHistory{}).
		Select("anime_id as content_id, COUNT(*) as views").
		Where("watched_at >= ?", weekAgo).
		Group("anime_id").
		Order("views DESC").
		Limit(20).
		Scan(&thisWeek).Error
	if err != nil {
		return nil, err
	}

	entries := make([]TrendingEntry, len(thisWeek))
	for i, r := range thisWeek {
		var lastWeekViews int64
		s.db.WithContext(ctx).
			Model(&models.WatchHistory{}).
			Where("anime_id = ? AND watched_at >= ? AND watched_at < ?", r.ContentID, twoWeeksAgo, weekAgo).
			Count(&lastWeekViews)

		var anime models.Anime
		s.db.WithContext(ctx).First(&anime, "id = ?", r.ContentID)

		growth := 0.0
		if lastWeekViews > 0 {
			growth = float64(r.Views-lastWeekViews) / float64(lastWeekViews) * 100
		}
		entries[i] = TrendingEntry{
			ContentID: r.ContentID,
			Title:     anime.Title,
			Views:     r.Views,
			Growth:    growth,
		}
	}
	return entries, nil
}

func (s *AnalyticsService) GetPopularNew(ctx context.Context, limit int) ([]NewPopularEntry, error) {
	type newPopResult struct {
		ContentID  string `gorm:"content_id"`
		Title      string `gorm:"title"`
		Views      int64  `gorm:"views"`
		UploadedAt string `gorm:"uploaded_at"`
	}
	var results []newPopResult
	err := s.db.WithContext(ctx).
		Model(&models.MediaAsset{}).
		Select("anime_id as content_id, filename as title, 0 as views, created_at::text as uploaded_at").
		Order("created_at DESC").
		Limit(limit).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for i := range results {
		var viewCount int64
		s.db.WithContext(ctx).Model(&models.WatchHistory{}).Where("anime_id = ?", results[i].ContentID).Count(&viewCount)
		results[i].Views = viewCount
	}

	entries := make([]NewPopularEntry, len(results))
	for i, r := range results {
		entries[i] = NewPopularEntry{
			ContentID:  r.ContentID,
			Title:      r.Title,
			Views:      r.Views,
			UploadedAt: r.UploadedAt,
		}
	}
	return entries, nil
}

func (s *AnalyticsService) GetGeography(ctx context.Context) ([]GeographyEntry, error) {
	type geoResult struct {
		Country string `gorm:"country"`
		Views   int64  `gorm:"views"`
	}
	var results []geoResult
	err := s.db.WithContext(ctx).
		Model(&models.Device{}).
		Select("COALESCE(country, 'Unknown') as country, COUNT(*) as views").
		Group("country").
		Order("views DESC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var total int64
	for _, r := range results {
		total += r.Views
	}

	entries := make([]GeographyEntry, len(results))
	for i, r := range results {
		pct := 0.0
		if total > 0 {
			pct = float64(r.Views) / float64(total) * 100
		}
		entries[i] = GeographyEntry{Country: r.Country, Views: r.Views, Pct: pct}
	}
	return entries, nil
}

func (s *AnalyticsService) GetGeographyTopCountries(ctx context.Context, limit int) ([]GeographyEntry, error) {
	all, err := s.GetGeography(ctx)
	if err != nil {
		return nil, err
	}
	if len(all) > limit {
		all = all[:limit]
	}
	return all, nil
}

func (s *AnalyticsService) GetActiveUsers(ctx context.Context) ([]ActiveUsersEntry, error) {
	now := time.Now().UTC()
	monthAgo := now.AddDate(0, 0, -30)

	type activeResult struct {
		Date  string `gorm:"date"`
		Count int64  `gorm:"count"`
	}
	var results []activeResult
	err := s.db.WithContext(ctx).
		Model(&models.WatchHistory{}).
		Select("DATE(watched_at) as date, COUNT(DISTINCT user_id) as count").
		Where("watched_at >= ?", monthAgo).
		Group("DATE(watched_at)").
		Order("date ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	entries := make([]ActiveUsersEntry, len(results))
	for i, r := range results {
		entries[i] = ActiveUsersEntry{Date: r.Date, Count: r.Count}
	}
	return entries, nil
}

func (s *AnalyticsService) GetActiveUsersRetention(ctx context.Context) ([]RetentionEntry, error) {
	now := time.Now().UTC()
	weekAgo := now.AddDate(0, 0, -7)

	type retentionResult struct {
		Period    string  `gorm:"period"`
		UserCount int64   `gorm:"user_count"`
		Retained  int64   `gorm:"retained"`
	}
	var results []retentionResult
	query := `
		SELECT 
			'week' as period,
			COUNT(DISTINCT user_id) as user_count,
			COUNT(DISTINCT CASE WHEN watched_at >= ? THEN user_id END) as retained
		FROM watch_history
		WHERE watched_at >= ?
	`
	err := s.db.WithContext(ctx).Raw(query, weekAgo, now.AddDate(0, 0, -28)).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	entries := make([]RetentionEntry, len(results))
	for i, r := range results {
		rate := 0.0
		if r.UserCount > 0 {
			rate = float64(r.Retained) / float64(r.UserCount) * 100
		}
		entries[i] = RetentionEntry{
			Period:    r.Period,
			Rate:      rate,
			UserCount: r.UserCount,
		}
	}
	return entries, nil
}

func (s *AnalyticsService) GetActiveUsersSessions(ctx context.Context) ([]SessionEntry, error) {
	now := time.Now().UTC()
	monthAgo := now.AddDate(0, 0, -30)

	type sessResult struct {
		Period     string  `gorm:"period"`
		AvgSessions float64 `gorm:"avg_sessions"`
		TotalUsers int64   `gorm:"total_users"`
	}
	var results []sessResult
	query := `
		SELECT 
			'month' as period,
			COALESCE(AVG(session_count), 0) as avg_sessions,
			COUNT(DISTINCT user_id) as total_users
		FROM (
			SELECT user_id, COUNT(*) as session_count
			FROM watch_history
			WHERE watched_at >= ?
			GROUP BY user_id
		) sub
	`
	err := s.db.WithContext(ctx).Raw(query, monthAgo).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	entries := make([]SessionEntry, len(results))
	for i, r := range results {
		entries[i] = SessionEntry{
			Period:      r.Period,
			AvgSessions: r.AvgSessions,
			TotalUsers:  r.TotalUsers,
		}
	}
	return entries, nil
}
