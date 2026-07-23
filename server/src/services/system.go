package services

import (
	"context"
	"runtime"
	"time"

	redisclient "github.com/skygenesisenterprise/kami-sama/server/internal/redis"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemService struct {
	db        *gorm.DB
	redis     *redisclient.Client
	startTime time.Time
}

func NewSystemService(db *gorm.DB, redis *redisclient.Client) *SystemService {
	return &SystemService{db: db, redis: redis, startTime: time.Now().UTC()}
}

func (s *SystemService) GetServiceHealth(ctx context.Context) gin.H {
	dbHealth := "healthy"
	if err := s.db.WithContext(ctx).Raw("SELECT 1").Error; err != nil {
		dbHealth = "unhealthy"
	}

	redisHealth := "disabled"
	redisLatency := ""
	if s.redis != nil && s.redis.IsAvailable() {
		h := s.redis.Health(ctx)
		redisHealth = string(h.Status)
		redisLatency = h.Latency
	}

	return gin.H{
		"database": gin.H{
			"status":  dbHealth,
			"latency": "",
		},
		"redis": gin.H{
			"status":  redisHealth,
			"latency": redisLatency,
		},
		"jellyfin": gin.H{
			"status": "unknown",
		},
		"cdn": gin.H{
			"status": "unknown",
		},
		"updatedAt": time.Now().UTC().Format(time.RFC3339),
	}
}

func (s *SystemService) GetUptime() gin.H {
	uptime := time.Since(s.startTime)
	return gin.H{
		"startedAt": s.startTime.Format(time.RFC3339),
		"uptime":    uptime.String(),
		"seconds":   int(uptime.Seconds()),
	}
}

func (s *SystemService) GetMetrics() gin.H {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return gin.H{
		"cpu": gin.H{
			"description": "not available without external collector",
		},
		"memory": gin.H{
			"allocMB":       m.Alloc / 1024 / 1024,
			"totalAllocMB":  m.TotalAlloc / 1024 / 1024,
			"sysMB":         m.Sys / 1024 / 1024,
			"numGC":         m.NumGC,
			"description": "approximate — use external monitoring for accuracy",
		},
		"disk": gin.H{
			"description": "not available without external collector",
		},
		"goroutines": runtime.NumGoroutine(),
	}
}

func (s *SystemService) GetLogs(ctx context.Context, level, source, date string, limit, offset int) ([]gin.H, int64, error) {
	return []gin.H{}, 0, nil
}

func (s *SystemService) GetLogByID(ctx context.Context, logID string) (*gin.H, error) {
	return nil, utils.ErrNotFound
}

func (s *SystemService) SearchLogs(ctx context.Context, query string) ([]gin.H, error) {
	return []gin.H{}, nil
}

func (s *SystemService) GetQueueStatus(ctx context.Context) gin.H {
	return gin.H{
		"pending": 0,
		"running": 0,
		"failed":  0,
		"completed": 0,
		"total":   0,
	}
}

func (s *SystemService) ListQueueJobs(ctx context.Context, status, jobType string, limit, offset int) ([]gin.H, int64, error) {
	return []gin.H{}, 0, nil
}

func (s *SystemService) GetCacheStatus(ctx context.Context) gin.H {
	hitRate := 0.0
	usedMemory := ""
	var totalKeys int64

	if s.redis != nil && s.redis.IsAvailable() {
		if dbSize, err := s.redis.Raw.DBSize(ctx).Result(); err == nil {
			totalKeys = dbSize
		}
	}

	return gin.H{
		"hitRate":    hitRate,
		"size":       totalKeys,
		"usedMemory": usedMemory,
	}
}

func (s *SystemService) FlushCache(ctx context.Context) error {
	if s.redis != nil && s.redis.IsAvailable() {
		return s.redis.Raw.FlushAll(ctx).Err()
	}
	return nil
}

func (s *SystemService) FlushCacheByPattern(ctx context.Context, pattern string) error {
	if s.redis != nil && s.redis.IsAvailable() {
		iter := s.redis.Raw.Scan(ctx, 0, pattern, 100).Iterator()
		for iter.Next(ctx) {
			if err := s.redis.Raw.Del(ctx, iter.Val()).Err(); err != nil {
				return err
			}
		}
		return iter.Err()
	}
	return nil
}

func (s *SystemService) ListCacheKeys(ctx context.Context) ([]string, error) {
	return []string{}, nil
}

func (s *SystemService) DeleteCacheKey(ctx context.Context, key string) error {
	if s.redis != nil && s.redis.IsAvailable() {
		return s.redis.Raw.Del(ctx, key).Err()
	}
	return nil
}

func (s *SystemService) GetSearchStatus(ctx context.Context) gin.H {
	return gin.H{
		"status":  "not_configured",
		"healthy": false,
	}
}

func (s *SystemService) TriggerReindex(ctx context.Context) error {
	return nil
}

func (s *SystemService) ListSearchIndexes(ctx context.Context) ([]gin.H, error) {
	return []gin.H{}, nil
}

func (s *SystemService) GetSearchIndexStats(ctx context.Context, indexName string) (*gin.H, error) {
	return nil, utils.ErrNotFound
}

func (s *SystemService) UpdateSearchIndex(ctx context.Context, indexName string) error {
	return nil
}

func (s *SystemService) ListBackgroundJobs(ctx context.Context) ([]gin.H, error) {
	return []gin.H{
		{"id": "session-cleanup", "name": "Session Cleanup", "schedule": "*/30 * * * *", "status": "active", "lastRun": nil, "nextRun": nil},
		{"id": "refresh-token-cleanup", "name": "Refresh Token Cleanup", "schedule": "0 * * * *", "status": "active", "lastRun": nil, "nextRun": nil},
		{"id": "cache-warm", "name": "Cache Warm", "schedule": "0 4 * * *", "status": "paused", "lastRun": nil, "nextRun": nil},
	}, nil
}

func (s *SystemService) RunBackgroundJob(ctx context.Context, jobID string) error {
	return nil
}

func (s *SystemService) PauseBackgroundJob(ctx context.Context, jobID string) error {
	return nil
}

func (s *SystemService) ResumeBackgroundJob(ctx context.Context, jobID string) error {
	return nil
}

func (s *SystemService) GetBackgroundJobHistory(ctx context.Context) ([]gin.H, error) {
	return []gin.H{}, nil
}

func (s *SystemService) GetEncodingProfiles() []gin.H {
	return []gin.H{
		{"name": "1080p", "width": 1920, "height": 1080, "bitrate": "5000k"},
		{"name": "720p", "width": 1280, "height": 720, "bitrate": "2500k"},
		{"name": "480p", "width": 854, "height": 480, "bitrate": "1000k"},
	}
}
