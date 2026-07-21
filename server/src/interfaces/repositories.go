package interfaces

import (
	"context"
	"time"

	"github.com/kami-sama-fr/platform/server/src/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	ListStale(ctx context.Context, before time.Time, limit int) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
}

type UserSettingsRepository interface {
	GetByUserID(ctx context.Context, userID string) (*models.UserSettings, error)
	Upsert(ctx context.Context, settings *models.UserSettings) error
}

type NotificationPreferenceRepository interface {
	GetByUserID(ctx context.Context, userID string) (*models.NotificationPreference, error)
	Upsert(ctx context.Context, preference *models.NotificationPreference) error
}

type LocalCredentialRepository interface {
	Create(ctx context.Context, credential *models.LocalCredential) error
	GetByUserID(ctx context.Context, userID string) (*models.LocalCredential, error)
	Update(ctx context.Context, credential *models.LocalCredential) error
}

type AuthSessionRepository interface {
	Create(ctx context.Context, session *models.AuthSession) error
	GetByID(ctx context.Context, id string) (*models.AuthSession, error)
	ListActiveByUser(ctx context.Context, userID string, now time.Time) ([]models.AuthSession, error)
	ListByUser(ctx context.Context, userID string) ([]models.AuthSession, error)
	Update(ctx context.Context, session *models.AuthSession) error
	Revoke(ctx context.Context, id string, reason string, revokedAt time.Time) error
	RevokeAllByUser(ctx context.Context, userID string, reason string, revokedAt time.Time, exceptSessionID string) error
	RevokeFamily(ctx context.Context, familyID string, reason string, revokedAt time.Time) error
	DeleteExpired(ctx context.Context, before time.Time) error
}

type AuthRefreshTokenRepository interface {
	Create(ctx context.Context, token *models.AuthRefreshToken) error
	GetByHash(ctx context.Context, tokenHash string) (*models.AuthRefreshToken, error)
	GetByID(ctx context.Context, id string) (*models.AuthRefreshToken, error)
	Update(ctx context.Context, token *models.AuthRefreshToken) error
	RevokeFamily(ctx context.Context, familyID string, revokedAt time.Time) error
	DeleteExpired(ctx context.Context, before time.Time) error
}

type EmailVerificationTokenRepository interface {
	Create(ctx context.Context, token *models.EmailVerificationToken) error
	GetByHash(ctx context.Context, tokenHash string) (*models.EmailVerificationToken, error)
	Update(ctx context.Context, token *models.EmailVerificationToken) error
	DeleteExpired(ctx context.Context, before time.Time) error
}

type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *models.PasswordResetToken) error
	GetByHash(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error)
	Update(ctx context.Context, token *models.PasswordResetToken) error
	DeleteExpired(ctx context.Context, before time.Time) error
}

type AuthAuditEventRepository interface {
	Create(ctx context.Context, event *models.AuthAuditEvent) error
}

type WorkspaceRepository interface {
	Create(ctx context.Context, workspace *models.Workspace) error
	ListByUser(ctx context.Context, userID string) ([]models.Workspace, error)
	GetByID(ctx context.Context, id string) (*models.Workspace, error)
	Update(ctx context.Context, workspace *models.Workspace) error
	Archive(ctx context.Context, id string, archivedAt time.Time) error
}

type WorkspaceMemberRepository interface {
	Create(ctx context.Context, member *models.WorkspaceMember) error
	Get(ctx context.Context, workspaceID, userID string) (*models.WorkspaceMember, error)
	ListByWorkspace(ctx context.Context, workspaceID string) ([]models.WorkspaceMember, error)
	Update(ctx context.Context, member *models.WorkspaceMember) error
	Delete(ctx context.Context, workspaceID, userID string) error
}

type AuthAccountRepository interface {
	Create(ctx context.Context, account *models.AuthAccount) error
	GetByProvider(ctx context.Context, provider, providerAccountID string) (*models.AuthAccount, error)
	ListByUserID(ctx context.Context, userID string) ([]models.AuthAccount, error)
	GetByUserIDAndProvider(ctx context.Context, userID, provider string) (*models.AuthAccount, error)
	Update(ctx context.Context, account *models.AuthAccount) error
	Delete(ctx context.Context, id string) error
}

// Anime repository
type AnimeRepository interface {
	Create(ctx context.Context, anime *models.Anime) error
	GetByID(ctx context.Context, id string) (*models.Anime, error)
	GetBySlug(ctx context.Context, slug string) (*models.Anime, error)
	List(ctx context.Context, opts ListAnimeOpts) ([]models.Anime, int64, error)
	Update(ctx context.Context, anime *models.Anime) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, limit int) ([]models.Anime, error)
}

type ListAnimeOpts struct {
	Page, Limit      int
	Status           string
	Genres           []string
	Studio           string
	Year             int
	Season           string
	Sort             string
	Query            string
	Featured, Trending *bool
}

// Genre repository
type GenreRepository interface {
	Create(ctx context.Context, genre *models.Genre) error
	GetByID(ctx context.Context, id string) (*models.Genre, error)
	GetBySlug(ctx context.Context, slug string) (*models.Genre, error)
	List(ctx context.Context) ([]models.Genre, error)
	Update(ctx context.Context, genre *models.Genre) error
	Delete(ctx context.Context, id string) error
}

// Studio repository
type StudioRepository interface {
	Create(ctx context.Context, studio *models.Studio) error
	GetByID(ctx context.Context, id string) (*models.Studio, error)
	GetBySlug(ctx context.Context, slug string) (*models.Studio, error)
	List(ctx context.Context) ([]models.Studio, error)
	Update(ctx context.Context, studio *models.Studio) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, limit int) ([]models.Studio, error)
}

// Character repository
type CharacterRepository interface {
	Create(ctx context.Context, character *models.Character) error
	GetByID(ctx context.Context, id string) (*models.Character, error)
	GetBySlug(ctx context.Context, slug string) (*models.Character, error)
	List(ctx context.Context) ([]models.Character, error)
	Update(ctx context.Context, character *models.Character) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, limit int) ([]models.Character, error)
}

// Episode repository
type EpisodeRepository interface {
	Create(ctx context.Context, episode *models.Episode) error
	GetByID(ctx context.Context, id string) (*models.Episode, error)
	ListByAnime(ctx context.Context, animeID string, seasonID *string) ([]models.Episode, error)
	Update(ctx context.Context, episode *models.Episode) error
	Delete(ctx context.Context, id string) error
}

// MediaAsset repository
type MediaAssetRepository interface {
	Create(ctx context.Context, asset *models.MediaAsset) error
	GetByID(ctx context.Context, id string) (*models.MediaAsset, error)
	List(ctx context.Context, mediaType string, limit, offset int) ([]models.MediaAsset, int64, error)
	Update(ctx context.Context, asset *models.MediaAsset) error
	Delete(ctx context.Context, id string) error
}

// EncodingJob repository
type EncodingJobRepository interface {
	Create(ctx context.Context, job *models.EncodingJob) error
	GetByID(ctx context.Context, id string) (*models.EncodingJob, error)
	GetByMediaAssetID(ctx context.Context, mediaAssetID string) (*models.EncodingJob, error)
	List(ctx context.Context, status string, limit, offset int) ([]models.EncodingJob, int64, error)
	Update(ctx context.Context, job *models.EncodingJob) error
}

// Review repository
type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	GetByID(ctx context.Context, id string) (*models.Review, error)
	ListByAnime(ctx context.Context, animeID string, limit, offset int) ([]models.Review, int64, error)
	GetByUserAndAnime(ctx context.Context, userID, animeID string) (*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id string) error
}

// Comment repository
type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id string) (*models.Comment, error)
	ListByAnime(ctx context.Context, animeID string, limit, offset int) ([]models.Comment, int64, error)
	ListByReview(ctx context.Context, reviewID string, limit, offset int) ([]models.Comment, int64, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id string) error
}

// Report repository
type ReportRepository interface {
	Create(ctx context.Context, report *models.Report) error
	GetByID(ctx context.Context, id string) (*models.Report, error)
	List(ctx context.Context, status string, limit, offset int) ([]models.Report, int64, error)
	Update(ctx context.Context, report *models.Report) error
}

// Watchlist repository
type WatchlistRepository interface {
	Create(ctx context.Context, watchlist *models.Watchlist) error
	GetByID(ctx context.Context, id string) (*models.Watchlist, error)
	ListByUser(ctx context.Context, userID string) ([]models.Watchlist, error)
	Update(ctx context.Context, watchlist *models.Watchlist) error
	Delete(ctx context.Context, id string) error
	AddAnime(ctx context.Context, item *models.WatchlistItem) error
	RemoveAnime(ctx context.Context, watchlistID, animeID string) error
	ListAnime(ctx context.Context, watchlistID string) ([]models.WatchlistItem, error)
}

// WatchProgress repository
type WatchProgressRepository interface {
	Upsert(ctx context.Context, progress *models.WatchProgress) error
	GetByUserAndEpisode(ctx context.Context, userID, episodeID string) (*models.WatchProgress, error)
	ListByUser(ctx context.Context, userID string) ([]models.WatchProgress, error)
	GetContinueWatching(ctx context.Context, userID string, limit int) ([]models.WatchProgress, error)
}

// WatchHistory repository
type WatchHistoryRepository interface {
	Create(ctx context.Context, history *models.WatchHistory) error
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]models.WatchHistory, int64, error)
}

// Simulcast repository
type SimulcastRepository interface {
	Create(ctx context.Context, simulcast *models.Simulcast) error
	GetByID(ctx context.Context, id string) (*models.Simulcast, error)
	ListActive(ctx context.Context) ([]models.Simulcast, error)
	ListByWeek(ctx context.Context) (map[string][]models.Simulcast, error)
	Update(ctx context.Context, simulcast *models.Simulcast) error
	Delete(ctx context.Context, id string) error
}

// ReleaseSchedule repository
type ReleaseScheduleRepository interface {
	Create(ctx context.Context, schedule *models.ReleaseSchedule) error
	GetByID(ctx context.Context, id string) (*models.ReleaseSchedule, error)
	ListUpcoming(ctx context.Context, limit int) ([]models.ReleaseSchedule, error)
	Update(ctx context.Context, schedule *models.ReleaseSchedule) error
}

// Notification repository
type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	GetByID(ctx context.Context, id string) (*models.Notification, error)
	ListByUser(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]models.Notification, int64, error)
	UnreadCount(ctx context.Context, userID string) (int64, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
}

// SystemSetting repository
type SystemSettingRepository interface {
	GetByKey(ctx context.Context, key string) (*models.SystemSetting, error)
	List(ctx context.Context, category string) ([]models.SystemSetting, error)
	Upsert(ctx context.Context, setting *models.SystemSetting) error
	Delete(ctx context.Context, key string) error
}

// AuditLog repository
type AuditLogRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	List(ctx context.Context, workspaceID string, limit, offset int) ([]models.AuditLog, int64, error)
}

// Contact repository
type ContactRepository interface {
	Create(ctx context.Context, contact *models.Contact) error
	GetByID(ctx context.Context, id string) (*models.Contact, error)
	ListByOwner(ctx context.Context, ownerID string) ([]models.Contact, error)
	Update(ctx context.Context, contact *models.Contact) error
	Delete(ctx context.Context, id string) error
}

// ContactGroup repository
type ContactGroupRepository interface {
	Create(ctx context.Context, group *models.ContactGroup) error
	GetByID(ctx context.Context, id string) (*models.ContactGroup, error)
	ListByOwner(ctx context.Context, ownerID string) ([]models.ContactGroup, error)
	Update(ctx context.Context, group *models.ContactGroup) error
	Delete(ctx context.Context, id string) error
}

// WorkspaceSSOConfig repository
type WorkspaceSSOConfigRepository interface {
	GetByWorkspaceID(ctx context.Context, workspaceID string) (*models.WorkspaceSSOConfig, error)
	Upsert(ctx context.Context, config *models.WorkspaceSSOConfig) error
}

type RepositorySet interface {
	Users() UserRepository
	UserSettings() UserSettingsRepository
	NotificationPreferences() NotificationPreferenceRepository
	LocalCredentials() LocalCredentialRepository
	AuthSessions() AuthSessionRepository
	AuthRefreshTokens() AuthRefreshTokenRepository
	EmailVerificationTokens() EmailVerificationTokenRepository
	PasswordResetTokens() PasswordResetTokenRepository
	AuthAuditEvents() AuthAuditEventRepository
	AuthAccounts() AuthAccountRepository
	Workspaces() WorkspaceRepository
	WorkspaceMembers() WorkspaceMemberRepository
	Anime() AnimeRepository
	Genres() GenreRepository
	Studios() StudioRepository
	Characters() CharacterRepository
	Episodes() EpisodeRepository
	MediaAssets() MediaAssetRepository
	EncodingJobs() EncodingJobRepository
	Reviews() ReviewRepository
	Comments() CommentRepository
	Reports() ReportRepository
	Watchlists() WatchlistRepository
	WatchProgresses() WatchProgressRepository
	WatchHistories() WatchHistoryRepository
	Simulcasts() SimulcastRepository
	ReleaseSchedules() ReleaseScheduleRepository
	Notifications() NotificationRepository
	SystemSettings() SystemSettingRepository
	AuditLogs() AuditLogRepository
	Contacts() ContactRepository
	ContactGroups() ContactGroupRepository
	WorkspaceSSOConfigs() WorkspaceSSOConfigRepository
}
